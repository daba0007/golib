package scan

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	manuf "github.com/timest/gomanuf"
)

const (
	// START 启动
	START = "start"
	// END 结束
	END = "end"
)

// HostNetworkInfo 本机网络信息
type HostNetworkInfo struct {
	ipNet      *net.IPNet
	localHaddr net.HardwareAddr
	iface      string
	data       map[string]Info
	t          *time.Ticker
	do         chan string
	duration   time.Duration
}

// NewScanner 初始化扫描器
func NewScanner(f string) (*HostNetworkInfo, error) {
	var ifs []net.Interface
	var err error
	if f == "" {
		ifs, err = net.Interfaces()
	} else {
		var it *net.Interface
		it, err = net.InterfaceByName(f)
		if err == nil {
			ifs = append(ifs, *it)
		}
	}
	if err != nil {
		return nil, err
	}
	var info HostNetworkInfo
	for _, it := range ifs {
		addr, _ := it.Addrs()
		for _, a := range addr {
			if ip, ok := a.(*net.IPNet); ok && !ip.IP.IsLoopback() {
				if ip.IP.To4() != nil {
					info.ipNet = ip
					info.localHaddr = it.HardwareAddr
					info.iface = it.Name
					info.data = make(map[string]Info)
					info.do = make(chan string)
					// 默认扫描30s
					info.duration = 30 * time.Second
					goto END
				}
			}
		}
	}
END:
	if info.ipNet == nil || len(info.localHaddr) == 0 {
		return nil, fmt.Errorf("无法获取本地网络信息")
	}
	return &info, nil
}

// ScanSubnet 网段扫描
func (hostinfo *HostNetworkInfo) ScanSubnet(ctx context.Context) error {
	go hostinfo.listenARP(ctx)
	go hostinfo.listenMDNS(ctx)
	go hostinfo.listenNBNS(ctx)
	go hostinfo.sendARP()
	go hostinfo.localHost()
	hostinfo.t = time.NewTicker(hostinfo.duration)
	for {
		select {
		case <-hostinfo.t.C:
			hostinfo.PrintData()
			ctx.Done()
			goto END
		case d := <-hostinfo.do:
			switch d {
			case START:
				hostinfo.t.Stop()
			case END:
				// 接收到新数据，重置2秒的计数器
				hostinfo.t = time.NewTicker(2 * time.Second)
			}
		}
	}
END:
	return nil
}

// SetDuration 设置扫描的时间
func (hostinfo *HostNetworkInfo) SetDuration(t time.Duration) {
	hostinfo.duration = t
}

// listenARP 监听ARP
func (hostinfo *HostNetworkInfo) listenARP(ctx context.Context) {
	// 为网卡打开一个句柄
	handle, err := pcap.OpenLive(hostinfo.iface, 1024, false, 10*time.Second)
	if err != nil {
		log.Fatal("pcap打开失败:", err)
	}
	defer handle.Close()
	// 过滤收集arp包
	handle.SetBPFFilter("arp")
	// 获取收集到的数据包集合
	ps := gopacket.NewPacketSource(handle, handle.LinkType())
	for {
		select {
		case <-ctx.Done():
			return
		case p := <-ps.Packets():
			arp := p.Layer(layers.LayerTypeARP).(*layers.ARP)
			// arp operation表示回应包
			if arp.Operation == 2 {
				mac := net.HardwareAddr(arp.SourceHwAddress)
				m := manuf.Search(mac.String())
				fmt.Println(ParseIP(arp.SourceProtAddress).String(), mac, m)
				hostinfo.pushData(ParseIP(arp.SourceProtAddress).String(), mac, "", m)
				if strings.Contains(m, "Apple") {
					go hostinfo.sendMdns(ParseIP(arp.SourceProtAddress), mac)
				} else {
					go hostinfo.sendNbns(ParseIP(arp.SourceProtAddress), mac)
				}
			}
		}
	}
}

// listenMDNS 监听MDNS
func (hostinfo *HostNetworkInfo) listenMDNS(ctx context.Context) error {
	handle, err := pcap.OpenLive(hostinfo.iface, 1024, false, 10*time.Second)
	if err != nil {
		return fmt.Errorf("pcap打开失败:#%v", err)
	}
	defer handle.Close()
	handle.SetBPFFilter("udp and port 5353")
	ps := gopacket.NewPacketSource(handle, handle.LinkType())
	for {
		select {
		case <-ctx.Done():
			return nil
		case p := <-ps.Packets():
			if len(p.Layers()) == 4 {
				c := p.Layers()[3].LayerContents()
				if c[2] == 0x84 && c[3] == 0x00 && c[6] == 0x00 && c[7] == 0x01 {
					// 从网络层(ipv4)拿IP, 不考虑IPv6
					i := p.Layer(layers.LayerTypeIPv4)
					if i == nil {
						continue
					}
					ipv4 := i.(*layers.IPv4)
					ip := ipv4.SrcIP.String()
					// 把 hostname 存入到数据库
					h := parseMdns(c)
					if len(h) > 0 {
						fmt.Println(ip, h)
						hostinfo.pushData(ip, nil, h, "")
					}
				}
			}
		}
	}
}

// sendMdns 发送Mdns包
func (hostinfo *HostNetworkInfo) sendMdns(ip IP, mhaddr net.HardwareAddr) error {
	srcIP := net.ParseIP(hostinfo.ipNet.IP.String()).To4()
	dstIP := net.ParseIP(ip.String()).To4()
	ether := &layers.Ethernet{
		SrcMAC:       hostinfo.localHaddr,
		DstMAC:       mhaddr,
		EthernetType: layers.EthernetTypeIPv4,
	}

	ip4 := &layers.IPv4{
		Version:  uint8(4),
		IHL:      uint8(5),
		TTL:      uint8(255),
		Protocol: layers.IPProtocolUDP,
		SrcIP:    srcIP,
		DstIP:    dstIP,
	}
	bf := newBuffer()
	mdns(bf, ip.String())
	udpPayload := bf.data
	udp := &layers.UDP{
		SrcPort: layers.UDPPort(60666),
		DstPort: layers.UDPPort(5353),
	}
	udp.SetNetworkLayerForChecksum(ip4)
	udp.Payload = udpPayload // todo
	buffer := gopacket.NewSerializeBuffer()
	opt := gopacket.SerializeOptions{
		FixLengths:       true, // 自动计算长度
		ComputeChecksums: true, // 自动计算checksum
	}
	err := gopacket.SerializeLayers(buffer, opt, ether, ip4, udp, gopacket.Payload(udpPayload))
	if err != nil {
		return fmt.Errorf("Serialize layers出现问题:%#v", err)
	}
	outgoingPacket := buffer.Bytes()

	handle, err := pcap.OpenLive(hostinfo.iface, 1024, false, 10*time.Second)
	if err != nil {
		return fmt.Errorf("pcap打开失败:%#v", err)
	}
	defer handle.Close()
	err = handle.WritePacketData(outgoingPacket)
	if err != nil {
		return fmt.Errorf("发送udp数据包失败")
	}
	return nil
}

// listenNBNS 监听NBNS
func (hostinfo *HostNetworkInfo) listenNBNS(ctx context.Context) error {
	handle, err := pcap.OpenLive(hostinfo.iface, 1024, false, 10*time.Second)
	if err != nil {
		return fmt.Errorf("pcap打开失败:#%v", err)
	}
	defer handle.Close()
	handle.SetBPFFilter("udp and port 137 and dst host " + hostinfo.ipNet.IP.String())
	ps := gopacket.NewPacketSource(handle, handle.LinkType())
	for {
		select {
		case <-ctx.Done():
			return nil
		case p := <-ps.Packets():
			if len(p.Layers()) == 4 {
				c := p.Layers()[3].LayerContents()
				if len(c) > 8 && c[2] == 0x84 && c[3] == 0x00 && c[6] == 0x00 && c[7] == 0x01 {
					// 从网络层(ipv4)拿IP, 不考虑IPv6
					i := p.Layer(layers.LayerTypeIPv4)
					if i == nil {
						continue
					}
					ipv4 := i.(*layers.IPv4)
					ip := ipv4.SrcIP.String()
					// 把 hostname 存入到数据库
					m := parseNBNS(c)
					if len(m) > 0 {
						fmt.Println(ip, m)
						hostinfo.pushData(ip, nil, m, "")
					}
				}
			}
		}
	}
}

// SendNbns 发送Nbns包
func (hostinfo *HostNetworkInfo) sendNbns(ip IP, mhaddr net.HardwareAddr) error {
	srcIP := net.ParseIP(hostinfo.ipNet.IP.String()).To4()
	dstIP := net.ParseIP(ip.String()).To4()
	ether := &layers.Ethernet{
		SrcMAC:       hostinfo.localHaddr,
		DstMAC:       mhaddr,
		EthernetType: layers.EthernetTypeIPv4,
	}

	ip4 := &layers.IPv4{
		Version:  uint8(4),
		IHL:      uint8(5),
		TTL:      uint8(255),
		Protocol: layers.IPProtocolUDP,
		SrcIP:    srcIP,
		DstIP:    dstIP,
	}
	bf := newBuffer()
	nbns(bf)
	udpPayload := bf.data
	udp := &layers.UDP{
		SrcPort: layers.UDPPort(61666),
		DstPort: layers.UDPPort(137),
	}
	udp.SetNetworkLayerForChecksum(ip4)
	udp.Payload = udpPayload
	buffer := gopacket.NewSerializeBuffer()
	opt := gopacket.SerializeOptions{
		FixLengths:       true, // 自动计算长度
		ComputeChecksums: true, // 自动计算checksum
	}
	err := gopacket.SerializeLayers(buffer, opt, ether, ip4, udp, gopacket.Payload(udpPayload))
	if err != nil {
		return fmt.Errorf("Serialize layers出现问题:%#v", err)
	}
	outgoingPacket := buffer.Bytes()

	handle, err := pcap.OpenLive(hostinfo.iface, 1024, false, 10*time.Second)
	if err != nil {
		return fmt.Errorf("pcap打开失败:%#v", err)
	}
	defer handle.Close()
	err = handle.WritePacketData(outgoingPacket)
	if err != nil {
		return fmt.Errorf("发送udp数据包失败")
	}
	return nil
}

// pushData 将抓到的数据集加入到data中，同时重置计时器
func (hostinfo *HostNetworkInfo) pushData(ip string, mac net.HardwareAddr, hostname, manuf string) {
	// 停止计时器
	hostinfo.do <- START
	var mu sync.RWMutex
	mu.RLock()
	defer func() {
		// 重置计时器
		hostinfo.do <- END
		mu.RUnlock()
	}()
	if _, ok := hostinfo.data[ip]; !ok {
		hostinfo.data[ip] = Info{Mac: mac, Hostname: hostname, Manuf: manuf}
		return
	}
	info := hostinfo.data[ip]
	if len(hostname) > 0 && len(info.Hostname) == 0 {
		info.Hostname = hostname
	}
	if len(manuf) > 0 && len(info.Manuf) == 0 {
		info.Manuf = manuf
	}
	if mac != nil {
		info.Mac = mac
	}
	hostinfo.data[ip] = info
}

// PrintData 打印数据
func (hostinfo *HostNetworkInfo) PrintData() {
	var keys IPSlice
	for k := range hostinfo.data {
		keys = append(keys, ParseIPString(k))
	}
	sort.Sort(keys)
	for _, k := range keys {
		d := hostinfo.data[k.String()]
		mac := ""
		if d.Mac != nil {
			mac = d.Mac.String()
		}
		fmt.Printf("%-15s %-17s %-30s %-10s\n", k.String(), mac, d.Hostname, d.Manuf)
	}
}

// localHost 加入本地数据
func (hostinfo *HostNetworkInfo) localHost() {
	host, _ := os.Hostname()
	hostinfo.data[hostinfo.ipNet.IP.String()] = Info{Mac: hostinfo.localHaddr, Hostname: strings.TrimSuffix(host, ".local"), Manuf: manuf.Search(hostinfo.localHaddr.String())}
}

// sendARP 发送arp
func (hostinfo *HostNetworkInfo) sendARP() {
	// ips 是内网IP地址集合
	ips := getHostZone(hostinfo.ipNet)
	for _, ip := range ips {
		go sendArpPackage(ip, hostinfo)
	}
}

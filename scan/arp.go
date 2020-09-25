package scan

import (
	"fmt"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// sendArpPackage 发送ARP包
func sendArpPackage(ip IP, hostinfo *HostNetworkInfo) error {
	srcIP := net.ParseIP(hostinfo.ipNet.IP.String()).To4()
	dstIP := net.ParseIP(ip.String()).To4()
	if srcIP == nil || dstIP == nil {
		return fmt.Errorf("ip 解析出问题")
	}
	// 以太网首部
	// EthernetType 0x0806  ARP 发送ARP协议包，目标MAC地址未知 FF-FF-FF
	ether := &layers.Ethernet{
		SrcMAC:       hostinfo.localHaddr,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}

	a := &layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     uint8(6),
		ProtAddressSize:   uint8(4),
		Operation:         uint16(1), // 0x0001 arp request 0x0002 arp response
		SourceHwAddress:   hostinfo.localHaddr,
		SourceProtAddress: srcIP,
		DstHwAddress:      net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		DstProtAddress:    dstIP,
	}

	buffer := gopacket.NewSerializeBuffer()
	var opt gopacket.SerializeOptions
	gopacket.SerializeLayers(buffer, opt, ether, a)
	outgoingPacket := buffer.Bytes()
	// 为网卡打开一个句柄
	handle, err := pcap.OpenLive(hostinfo.iface, 2048, false, 30*time.Second)
	if err != nil {
		return fmt.Errorf("pcap打开失败:%#v", err)
	}
	defer handle.Close()
	// 使用句柄发送arp包
	err = handle.WritePacketData(outgoingPacket)
	if err != nil {
		return fmt.Errorf("发送arp数据包失败")
	}
	return nil
}

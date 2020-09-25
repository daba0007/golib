# SCAN

SCAN会实现一些扫描功能，通过发送网络包请求，根据返回的数据来解析获取主机信息

## ARP

arp协议位于网络层，用于请求主机的mac地址。因此，我们可以伪造往网段内发送arp数据包，通过每个主机给我们返回的单播包来解析获取各主机的mac地址。
1. 通过网卡的IP和子网掩码计算出内网的IP范围
2. 向内网广播ARP请求
3. 监听并抓取ARP相应包，记录mac地址
4. 根据mac地址计算厂家信息

其实MAC地址都是包含了厂家的信息的，我们可以通过mac地址计算出厂家信息
```
...
00:03:8F	Weinsche	Weinschel Corporation
00:03:90	DigitalV	Digital Video Communications, Inc.
00:03:91	Advanced	Advanced Digital Broadcast, Ltd.
00:03:92	HyundaiT	Hyundai Teletek Co., Ltd.
00:03:93	Apple	Apple, Inc.
00:03:94	ConnectO	Connect One
00:03:95	Californ	California Amplifier
00:03:96	EzCast	EZ Cast Co., Ltd.
00:03:97	Watchfro	Watchfront Limited
...
```
在获取ip和mac地址后，其实也可以通过MDNS和NBNS协议去查询主机名

### 最终实现

```Go
// 定义扫描器
// "networkCard" -> 网卡名称
scanner,err := NewScanner("networkCard")
// 默认对网段内所有ip进行扫描
// 可以定义扫描时间长短，在扫描获取完所有主机后会自动停止 SetDuration(2) -> 设置为2s
scanner.ScanSubnet()
```
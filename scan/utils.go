package scan

import (
	"net"
)

// IP uint32
type IP uint32

// Info 机器信息
type Info struct {
	// IP地址
	Mac net.HardwareAddr
	// 主机名
	Hostname string
	// 厂商信息
	Manuf string
}

// Buffer struct
type Buffer struct {
	data  []byte
	start int
}

// PrependBytes 复制buffer
func (b *Buffer) prependBytes(n int) []byte {
	length := cap(b.data) + n
	newData := make([]byte, length)
	copy(newData, b.data)
	b.start = cap(b.data)
	b.data = newData
	return b.data[b.start:]
}

// NewBuffer 创建buffer
func newBuffer() *Buffer {
	return &Buffer{}
}

// Reverse 反转字符串
func reverse(s string) (result string) {
	for _, v := range s {
		result = string(v) + result
	}
	return
}

// IPSlice slice
type IPSlice []IP

func (ip IPSlice) Len() int { return len(ip) }

func (ip IPSlice) Swap(i, j int) {
	ip[i], ip[j] = ip[j], ip[i]
}

func (ip IPSlice) Less(i, j int) bool {
	return ip[i] < ip[j]
}

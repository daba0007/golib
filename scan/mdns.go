package scan

import (
	"bytes"
	"encoding/binary"
	"strings"
)

// 根据ip生成含mdns请求包，包存储在 buffer里
func mdns(buffer *Buffer, ip string) {
	b := buffer.prependBytes(12)
	binary.BigEndian.PutUint16(b, uint16(0))          // 0x0000 标识
	binary.BigEndian.PutUint16(b[2:], uint16(0x0100)) // 标识
	binary.BigEndian.PutUint16(b[4:], uint16(1))      // 问题数
	binary.BigEndian.PutUint16(b[6:], uint16(0))      // 资源数
	binary.BigEndian.PutUint16(b[8:], uint16(0))      // 授权资源记录数
	binary.BigEndian.PutUint16(b[10:], uint16(0))     // 额外资源记录数
	// 查询问题
	ipList := strings.Split(ip, ".")
	for j := len(ipList) - 1; j >= 0; j-- {
		ip := ipList[j]
		b = buffer.prependBytes(len(ip) + 1)
		b[0] = uint8(len(ip))
		for i := 0; i < len(ip); i++ {
			b[i+1] = uint8(ip[i])
		}
	}
	b = buffer.prependBytes(8)
	b[0] = 7 // 后续总字节
	copy(b[1:], []byte{'i', 'n', '-', 'a', 'd', 'd', 'r'})
	b = buffer.prependBytes(5)
	b[0] = 4 // 后续总字节
	copy(b[1:], []byte{'a', 'r', 'p', 'a'})
	b = buffer.prependBytes(1)
	// terminator
	b[0] = 0
	// type 和 classIn
	b = buffer.prependBytes(4)
	binary.BigEndian.PutUint16(b, uint16(12))
	binary.BigEndian.PutUint16(b[2:], 1)
}

// ParseMdns 从 mdns响应报文中获取主机名
// 参数data  开头是 dns的协议头 0x0000 0x8400 0x0000 0x0001(ans) 0x0000 0x0000
func parseMdns(data []byte) string {
	var buf bytes.Buffer
	i := bytes.Index(data, []byte{0x05, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x00})
	if i < 0 {
		return ""
	}

	for s := i - 1; s > 1; s-- {
		num := i - s
		if s-2 < 0 {
			break
		}
		// 包括 .local_ 7 个字符
		if bto16([]byte{data[s-2], data[s-1]}) == uint16(num+7) {
			return reverse(buf.String())
		}
		buf.WriteByte(data[s])
	}

	return ""
}

func bto16(b []byte) uint16 {
	if len(b) != 2 {
		return uint16(0)
	}
	return uint16(b[0])<<8 + uint16(b[1])
}

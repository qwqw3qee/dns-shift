package model

// IPType 定义 IP 地址类型
type IPType byte

const (
	TypeIPv4 IPType = 4
	TypeIPv6 IPType = 6
)

// PrefixLen 返回 IP 类型的默认子网前缀长度（以位为单位）。
func (x IPType) PrefixLen() int {
	switch x {
	case TypeIPv4:
		return 24
	case TypeIPv6:
		return 112
	default:
		return 24
	}
}

// AddrLen 返回 IP 地址的总长度（以位为单位）。
func (x IPType) AddrLen() int {
	switch x {
	case TypeIPv4:
		return 32
	case TypeIPv6:
		return 128
	default:
		return 32
	}
}

// RecordType 返回 IP 类型对应的 DNS 记录类型字符串。
func (x IPType) RecordType() string {
	switch x {
	case TypeIPv4:
		return "A"
	case TypeIPv6:
		return "AAAA"
	default:
		return "A"
	}
}

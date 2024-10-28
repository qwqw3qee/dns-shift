package util

import (
	"dns-shift/model"
	"errors"
	"fmt"
	probing "github.com/prometheus-community/pro-bing"
	"net"
	"time"
)

// CheckIPv6Support 检测是否支持IPv6访问
func CheckIPv6Support() bool {
	conn, err := net.DialTimeout("tcp6", "2001:4860:4860::8888:53", time.Second)
	if err != nil {
		return false // 不支持IPv6
	}
	defer conn.Close()
	return true // 支持IPv6
}

// PingAddress 用于 ping 一个地址，允许设置超时时间和 ping 的次数
func PingAddress(addr string, count int, timeout time.Duration) (time.Duration, error) {
	pinger, err := probing.NewPinger(addr)
	if err != nil {
		return 0, fmt.Errorf("创建 ping 对象失败: %v", err)
	}

	// 配置 pinger
	pinger.Count = count
	pinger.Timeout = timeout
	pinger.SetPrivileged(true)

	// 执行 ping 操作
	if err = pinger.Run(); err != nil {
		return 0, fmt.Errorf("执行 ping 失败: %v", err)
	}

	// 获取统计数据
	stats := pinger.Statistics()
	if len(stats.Rtts) == 0 {
		return 0, errors.New("ping 超时，未收到任何响应")
	}

	// 输出 ping 结果
	// log.Printf("Ping %s 成功: %+v\n", addr, *stats)
	return stats.AvgRtt, nil
}

// incrementIPv4 增加或减少 IPv4 地址
func incrementIPv4(ip net.IP, inc int) (net.IP, error) {
	ip4 := ip.To4()
	if ip4 == nil {
		return nil, errors.New("invalid IPv4 address")
	}
	ipCopy := make(net.IP, len(ip4))
	copy(ipCopy, ip4)

	lastByte := int(ipCopy[3])
	lastByte += inc
	if lastByte > 254 || lastByte < 1 {
		return nil, errors.New("IPv4 address out of range")
	}
	ipCopy[3] = byte(lastByte)
	return ipCopy, nil
}

// ShiftIPv4 进行 IPv4 递增或递减，并根据 checkFlag 决定是否检测 IP 的可用性
func ShiftIPv4(oriIP net.IP, checkFlag bool) (net.IP, error) {

	// 增加 IP 地址
	for i := 1; i <= 5; i++ {
		newIP, err := incrementIPv4(oriIP, i)
		if err != nil {
			break
		}
		if !checkFlag {
			return newIP, nil // 不需要检查时直接返回
		}
		// 检查 IP 的可用性
		if _, err = PingAddress(newIP.String(), 3, 3*time.Second); err == nil {
			return newIP, nil // 找到可用 IP，返回
		}
	}

	// 减少 IP 地址
	for i := -1; i >= -5; i-- {
		newIP, err := incrementIPv4(oriIP, i)
		if err != nil {
			break
		}
		if !checkFlag {
			return newIP, nil
		}
		// 检查 IP 的可用性
		if _, err = PingAddress(newIP.String(), 3, 3*time.Second); err == nil {
			return newIP, nil
		}
	}

	// 如果找不到可用 IP
	return nil, errors.New("no available IPv4 found")
}

// incrementIPv6 增加或减少 IPv6 地址
func incrementIPv6(ip net.IP, inc int) (net.IP, error) {
	ip6 := ip.To16()
	if ip6 == nil {
		return nil, errors.New("invalid IPv6 address")
	}
	ipCopy := make(net.IP, len(ip6))
	copy(ipCopy, ip6)

	// 递增最后两个字节
	for i := 15; i >= 0; i-- {
		ipCopy[i] += byte(inc)
		if ipCopy[i] != 0 {
			break
		}
	}
	return ipCopy, nil
}

// ShiftIPv6 进行 IPv6 递增或递减，并根据 checkFlag 决定是否检测 IP 的可用性
func ShiftIPv6(oriIP net.IP, checkFlag bool) (net.IP, error) {
	// 增加 IP 地址
	for i := 1; i <= 5; i++ {
		newIP, err := incrementIPv6(oriIP, i)
		if err != nil {
			continue
		}
		if !checkFlag {
			return newIP, nil // 不需要检查时直接返回
		}
		// 检查 IP 的可用性
		if _, err = PingAddress(newIP.String(), 3, 3*time.Second); err == nil {
			return newIP, nil // 找到可用 IP，返回
		}
	}

	// 减少 IP 地址
	for i := -1; i >= -5; i-- {
		newIP, err := incrementIPv6(oriIP, i)
		if err != nil {
			continue
		}
		if !checkFlag {
			return newIP, nil
		}
		// 检查 IP 的可用性
		if _, err = PingAddress(newIP.String(), 3, 3*time.Second); err == nil {
			return newIP, nil
		}
	}

	// 如果找不到可用 IP
	return nil, errors.New("no available IPv6 found")
}

// CompareIPPrefix 检查两个 IP 地址是否在同一网段
func CompareIPPrefix(ip1, ip2 net.IP, prefixLength int) bool {
	if ip1 == nil || ip2 == nil {
		return false // 确保 IP 不为 nil
	}

	// 根据 IP 类型处理
	if ip1.To4() != nil && ip2.To4() != nil { // IPv4
		if prefixLength > 32 {
			prefixLength = 32
		}
		mask := net.CIDRMask(prefixLength, 32) // 创建 IPv4 掩码
		return ip1.Mask(mask).Equal(ip2.Mask(mask))
	} else if ip1.To16() != nil && ip2.To16() != nil { // IPv6
		mask := net.CIDRMask(prefixLength, 128) // 创建 IPv6 掩码
		return ip1.Mask(mask).Equal(ip2.Mask(mask))
	}

	return false // 不同类型的 IP
}
func ShiftIP(oriIP net.IP, checkFlag bool, ipType model.IPType) (net.IP, error) {
	switch ipType {
	case model.TypeIPv4:
		return ShiftIPv4(oriIP, checkFlag)
	case model.TypeIPv6:
		return ShiftIPv6(oriIP, checkFlag)
	}
	return ShiftIPv4(oriIP, checkFlag)
}

// CompareIPLists 检查两个有序 IP 列表是否相同
func CompareIPLists(list1, list2 []net.IP, prefixLength int) bool {
	if len(list1) != len(list2) {
		return false
	}
	i, j := 0, 0 // 初始化两个指针

	for i < len(list1) && j < len(list2) {
		if CompareIPPrefix(list1[i], list2[j], prefixLength) {
			i++ // 找到匹配，移动 list1 的指针
			j++ // 找到匹配，移动 list2 的指针
		} else if list1[i].String() < list2[j].String() {
			i++ // list1 的 IP 小于 list2 的 IP，移动 list1 的指针
		} else {
			j++ // list2 的 IP 小于 list1 的 IP，移动 list2 的指针
		}
	}

	return i == len(list1) && j == len(list2) // 如果 list1 的指针到达末尾，说明所有 IP 都找到匹配
}

// HasCommonIP 检查两个有序 IP 列表是否相同
func HasCommonIP(list1, list2 []net.IP, prefixLength int) bool {
	i, j := 0, 0 // 初始化两个指针

	for i < len(list1) && j < len(list2) {
		if CompareIPPrefix(list1[i], list2[j], prefixLength) {
			return true
		} else if list1[i].String() < list2[j].String() {
			i++ // list1 的 IP 小于 list2 的 IP，移动 list1 的指针
		} else {
			j++ // list2 的 IP 小于 list1 的 IP，移动 list2 的指针
		}
	}

	return false
}

// ConvertIPListToStrings 将 []net.IP 转换为 []string
func ConvertIPListToStrings(ipList []net.IP) []string {
	strList := make([]string, len(ipList))
	for i, ip := range ipList {
		strList[i] = ip.String()
	}
	return strList
}

func GenShiftIPList(ipList []net.IP, checkFlag bool, ipType model.IPType) ([]net.IP, error) {
	if checkFlag && ipType == model.TypeIPv6 && !CheckIPv6Support() {
		checkFlag = false
	}
	shiftIPlist := make([]net.IP, 0, len(ipList))
	for _, ip := range ipList {
		shiftIP, err := ShiftIP(ip, checkFlag, ipType)
		if err != nil {
			continue // 如果出错则跳过该IP
		}
		shiftIPlist = append(shiftIPlist, shiftIP)
	}
	if len(shiftIPlist) == 0 {
		return nil, errors.New("无法生成可用IP")
	}
	return shiftIPlist, nil
}

package util

import (
	"bytes"
	"dns-shift/model"
	"errors"
	"fmt"
	"github.com/miekg/dns"
	"net"
	"sort"
)

// DnsLookupA 查询域名的 A 记录并返回对应的 IP 地址列表
func DnsLookupA(domain, dnsAddress string) ([]net.IP, error) {
	// 构造 DNS 客户端
	client := dns.Client{
		SingleInflight: false, // 禁用缓存，确保每次查询都是新的请求
	}
	// 构造 DNS 查询报文，查询 A 记录
	msg := dns.Msg{}
	msg.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	msg.RecursionDesired = true
	// 发起 DNS 查询
	response, _, err := client.Exchange(&msg, dnsAddress)
	if err != nil {
		return nil, fmt.Errorf("dns query failed: %v", err)
	}
	// 检查查询是否成功
	if response.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("dns query failed: invalid answer for domain %s", domain)
	}
	// 提取 A 记录中的 IP 地址
	var ipList []net.IP
	for _, ans := range response.Answer {
		if record, ok := ans.(*dns.A); ok { // 只获取A记录
			ipList = append(ipList, record.A)
		}
	}
	// 如果没有找到 A 记录，返回错误
	if len(ipList) == 0 {
		return nil, errors.New("no A records found for domain: " + domain)
	}
	// 排序
	sort.Slice(ipList, func(i, j int) bool {
		return bytes.Compare(ipList[i], ipList[j]) < 0
	})
	return ipList, nil
}

// DnsLookupAAAA 查询域名的 AAAA 记录并返回对应的 IPv6 地址列表
func DnsLookupAAAA(domain, dnsAddress string) ([]net.IP, error) {
	// 构造 DNS 客户端
	client := dns.Client{
		SingleInflight: false, // 禁用缓存，确保每次查询都是新的请求
	}
	// 构造 DNS 查询报文，查询 AAAA 记录
	msg := dns.Msg{}
	msg.SetQuestion(dns.Fqdn(domain), dns.TypeAAAA)
	msg.RecursionDesired = true
	// 发起 DNS 查询
	response, _, err := client.Exchange(&msg, dnsAddress)
	if err != nil {
		return nil, fmt.Errorf("dns query failed: %v", err)
	}
	// 检查查询是否成功
	if response.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("dns query failed: invalid answer for domain %s", domain)
	}
	// 提取 AAAA 记录中的 IP 地址
	var ipList []net.IP
	for _, ans := range response.Answer {
		if record, ok := ans.(*dns.AAAA); ok { // 只获取AAAA记录
			ipList = append(ipList, record.AAAA)
		}
	}
	// 如果没有找到 AAAA 记录，返回错误
	if len(ipList) == 0 {
		return nil, errors.New("no AAAA records found for domain: " + domain)
	}
	// 排序
	sort.Slice(ipList, func(i, j int) bool {
		return bytes.Compare(ipList[i], ipList[j]) < 0
	})
	return ipList, nil
}
func DnsLookupIP(domain, dnsAddress string, ipType model.IPType) ([]net.IP, error) {
	switch ipType {
	case model.TypeIPv4:
		return DnsLookupA(domain, dnsAddress)
	case model.TypeIPv6:
		return DnsLookupAAAA(domain, dnsAddress)
	}
	return DnsLookupA(domain, dnsAddress)
}

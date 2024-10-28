package util

import (
	"fmt"
	"net"
	"testing"
)

func TestCompareIPPrefix(t *testing.T) {
	ip1 := net.ParseIP("104.17.110.184")
	ip2 := net.ParseIP("104.17.110.185")
	println(CompareIPPrefix(ip1, ip2, 24))
	println(CompareIPPrefix(ip1, ip2, 32))
	println(CompareIPPrefix(ip1, ip2, 128))
}
func TestCompareIPLists(t *testing.T) {
	list1, _ := DnsLookupA("dash.cloudflare.com", "8.8.8.8:53")
	list2, _ := DnsLookupA("cf-shift.12123123.xyz", "8.8.8.8:53")
	for _, ip := range list1 {
		fmt.Println(ip)
	}
	for _, ip := range list2 {
		fmt.Println(ip)
	}
	println(CompareIPLists(list1, list2, 24))
	println(CompareIPLists(list1, list2, 32))
}

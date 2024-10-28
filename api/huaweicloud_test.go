package api

import (
	"dns-shift/model"
	"dns-shift/util"
	"net"
	"strconv"
	"testing"
)

func TestHuaweiCloudAPI(t *testing.T) {
	configMap := model.ConfigMap{
		"ak": "xxxxx",
		"sk": "xxxxx",
	}
	domain := "cf-shift.examole.domain"
	huaweicloud := HuaweiCloudApi{params: configMap}
	huaweicloud.InitClient()
	ipList := make([]net.IP, 2)
	for i := range ipList {
		ipList[i] = net.ParseIP("192.168.3." + strconv.Itoa(i+2))
	}
	util.PrintErr("设置A记录失败", huaweicloud.SetRecord(ipList, domain, model.TypeIPv4))
	ipList6 := make([]net.IP, 2)
	for i := range ipList6 {
		ipList6[i] = net.ParseIP("aaee::" + strconv.Itoa(i+2))
	}
	util.PrintErr("设置AAAA记录失败", huaweicloud.SetRecord(ipList6, domain, model.TypeIPv6))
}

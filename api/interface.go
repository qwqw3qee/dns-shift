package api

import (
	"dns-shift/model"
	"errors"
	"net"
)

type DnsApi interface {
	SetRecord(ipList []net.IP, targetDomain string, ipType model.IPType) error
}

//注册-调用模式

// DnsApiRegistry 保存各 DnsType 对应的构造函数
var dnsApiRegistry = make(map[model.DnsType]func(model.ConfigMap) DnsApi)

// RegisterApi 将 DnsType 和 API 实现关联
func RegisterApi(dnsType model.DnsType, constructor func(model.ConfigMap) DnsApi) {
	dnsApiRegistry[dnsType] = constructor
}

// GetApi 根据 DnsType 获取对应的 API 实现
func GetApi(dnsServer model.DnsServerStruct) (DnsApi, error) {
	constructor, exists := dnsApiRegistry[dnsServer.Provider]
	if !exists {
		return nil, errors.New("unsupported DNS service type")
	}
	return constructor(dnsServer.ApiParam), nil
}

// 初始化注册各个 API 实现
func init() {
	RegisterApi(model.NoAPI, func(configMap model.ConfigMap) DnsApi {
		return &NoApi{}
	})
	RegisterApi(model.HuaweiCloud, NewHuaweiCloudApi)
}

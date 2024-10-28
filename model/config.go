package model

// ConfigStruct 用于存储配置文件的解析结果
type ConfigStruct struct {
	SourceDomain       string          `mapstructure:"sourceDomain"`
	TargetDomain       string          `mapstructure:"targetDomain"`
	QueryServer        string          `mapstructure:"queryServer"`
	EnablePingCheck    bool            `mapstructure:"enablePingCheck"`
	EnableIPv6         bool            `mapstructure:"enableIPv6"`
	EnableSuccessCheck bool            `mapstructure:"enableSuccessCheck"`
	SuccessCheckSecond int             `mapstructure:"successCheckSecond"`
	DnsServer          DnsServerStruct `mapstructure:"dnsServer"`
}
type DnsServerStruct struct {
	Provider DnsType   `mapstructure:"provider"` // DNS 服务商
	ApiParam ConfigMap `mapstructure:"apiParam"` // API 参数
}

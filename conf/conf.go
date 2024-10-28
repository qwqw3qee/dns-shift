package conf

import (
	"dns-shift/model"
	"dns-shift/util"
	"github.com/spf13/viper"
	"sync"
)

var (
	config *viper.Viper
	cfg    *model.ConfigStruct // 缓存解析后的配置
	once   sync.Once           // 确保配置初始化只运行一次
)

// setDefaultValue 设置默认值
func setDefaultValue() {
	config.SetDefault("sourceDomain", "dash.cloudflare.com")
	config.SetDefault("targetDomain", "yourdomain.example")
	config.SetDefault("queryServer", "8.8.8.8:53")
	config.SetDefault("enablePingCheck", true)
	config.SetDefault("enableSuccessCheck", true)
	config.SetDefault("successCheckSecond", 10)
	config.SetDefault("enableIPv6", true)
	config.SetDefault("dnsServer.provider", model.NoAPI)
}

// configInit 初始化配置，确保只运行一次
func configInit() {
	once.Do(func() {
		config = viper.New() // 初始化 viper 对象
		config.AddConfigPath(util.GetCurDir())
		config.SetConfigName("dns-shift")
		config.SetConfigType("yaml")
		setDefaultValue()
		// 读取配置文件
		if err := config.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				util.FatalErr("找不到配置文件", err)
			} else {
				util.FatalErr("读取配置文件出错", err)
			}
		}
		// 解析配置文件并缓存到结构体
		var tempConfig model.ConfigStruct
		if err := config.Unmarshal(&tempConfig); err != nil {
			util.FatalErr("解析配置出错", err)
		} else {
			cfg = &tempConfig // 缓存解析后的配置
		}
	})
}

// GetConfig 返回配置文件的解析结果
func GetConfig() *model.ConfigStruct {
	// 确保配置初始化已经执行
	configInit()
	return cfg
}

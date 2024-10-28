package main

import (
	"dns-shift/api"
	"dns-shift/conf"
	"dns-shift/model"
	"dns-shift/util"
	"log"
	"time"
)

func fullProcess(ipType model.IPType, config *model.ConfigStruct, api api.DnsApi) {
	sourceIPs, err := util.DnsLookupIP(config.SourceDomain, config.QueryServer, ipType)
	util.FatalErr("DNS查询失败", err)
	log.Printf("%s %s: %v\n", config.SourceDomain, ipType.RecordType(), sourceIPs)
	targetIPs, err := util.DnsLookupIP(config.TargetDomain, config.QueryServer, ipType)
	log.Printf("%s %s: %v\n", config.TargetDomain, ipType.RecordType(), targetIPs)
	if util.CompareIPLists(sourceIPs, targetIPs, ipType.PrefixLen()) &&
		!util.HasCommonIP(sourceIPs, targetIPs, ipType.AddrLen()) {
		log.Printf("%s 记录正常，无需修改\n", ipType.RecordType())
		return
	}
	newTargetIPs, err := util.GenShiftIPList(sourceIPs, config.EnablePingCheck, ipType)
	log.Printf("生成IP: %v\n", newTargetIPs)
	err = api.SetRecord(newTargetIPs, config.TargetDomain, ipType)
	if util.PrintErr("使用API设置记录出错", err) {
		return
	}

	if config.EnableSuccessCheck {
		waitSecond := time.Duration(config.SuccessCheckSecond)
		log.Printf("检查DNS设置情况，请等待%d秒...\n", waitSecond)
		time.Sleep(waitSecond * time.Second)
		updatedTargetIPs, err := util.DnsLookupIP(config.TargetDomain, config.QueryServer, ipType)
		if util.PrintErr("查询新DNS记录失败", err) ||
			util.CompareIPLists(newTargetIPs, updatedTargetIPs, ipType.AddrLen()) {
			log.Printf("%s 记录检查成功，DNS解析已生效\n", ipType.RecordType())
		} else {
			log.Printf("%s 检查失败，DNS解析未生效，可能是缓存未过期，请等待一段时间自行检查\n", ipType.RecordType())
		}
	}

}
func logPrintConfig(cfg *model.ConfigStruct) {
	log.Println("┌─配置文件信息")
	log.Println("├─原始域名:", cfg.SourceDomain)
	log.Println("├─目标域名:", cfg.TargetDomain)
	log.Println("├─API类型:", cfg.DnsServer.Provider)
	log.Println("├─启用IPv6:", cfg.EnableIPv6)
	log.Println("├─Ping检测:", cfg.EnablePingCheck)
	log.Println("├─完成检查:", cfg.EnableSuccessCheck)
	log.Println("└─DNS查询地址:", cfg.QueryServer)
}
func main() {
	log.Printf("=====DNS-Shift START=====")
	defer log.Printf("======DNS-Shift END======")
	config := conf.GetConfig()
	logPrintConfig(config)
	dnsApi, err := api.GetApi(config.DnsServer)
	util.FatalErr("获取API接口失败", err)
	fullProcess(model.TypeIPv4, config, dnsApi)
	if config.EnableIPv6 {
		fullProcess(model.TypeIPv6, config, dnsApi)
	}

}

package api

import (
	"dns-shift/model"
	"dns-shift/util"
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	dns "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dns/v2"
	model2 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dns/v2/model"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dns/v2/region"
	dns2 "github.com/miekg/dns"
	"net"
	"strings"
)

// HuaweiCloudApi 华为云Api的实现
type HuaweiCloudApi struct {
	params       model.ConfigMap
	client       *dns.DnsClient
	targetDomain string
	zoneId       string
}

func NewHuaweiCloudApi(params model.ConfigMap) DnsApi {
	api := &HuaweiCloudApi{params: params, targetDomain: "", zoneId: ""}
	util.FatalErr("初始化HuaweiCloudAPI失败", api.InitClient())
	return api
}
func (h *HuaweiCloudApi) InitClient() error {
	ak := h.params.GetString("ak")
	sk := h.params.GetString("sk")
	if ak == "" || sk == "" {
		return fmt.Errorf("get ak/sk error")
	}
	auth, err := basic.NewCredentialsBuilder().
		WithAk(ak).
		WithSk(sk).
		SafeBuild()
	if err != nil {
		return err
	}
	regionVal, _ := region.SafeValueOf("cn-north-1")
	c, err := dns.DnsClientBuilder().
		WithRegion(regionVal).
		WithCredential(auth).
		SafeBuild()
	if err != nil {
		return err
	}
	h.client = dns.NewDnsClient(c)
	return nil
}

// findMatchingDomain 遍历域名列表，找到与输入域名匹配的最合适的父域。
func (h *HuaweiCloudApi) findZoneIdByDomain(domain string, zoneList *[]model2.PublicZoneResp) (string, error) {
	// 将 zoneList 转换为 map 结构，以便快速查找
	zoneMap := make(map[string]string)
	for _, d := range *zoneList {
		if d.Name != nil && d.Id != nil {
			zoneMap[*d.Name] = *d.Id
		}
	}
	// 提前将 domain 分割为各级子域名，避免在循环中重复分割
	parts := strings.Split(domain, ".")
	for i := 0; i < len(parts); i++ {
		subDomain := strings.Join(parts[i:], ".") // 从第 i 个部分开始重组域名
		// 检查是否在 zoneMap 中存在匹配的域名
		if zoneId, exists := zoneMap[subDomain]; exists {
			return zoneId, nil // 找到匹配的域名
		}
	}
	// 如果没有找到匹配的域名，则返回错误
	return "", fmt.Errorf("no matching domain found for %s", domain)
}

func (h *HuaweiCloudApi) listPublicZones() (*[]model2.PublicZoneResp, error) {
	request := &model2.ListPublicZonesRequest{}
	response, err := h.client.ListPublicZones(request)
	if err != nil {
		return nil, err
	}
	if response.HttpStatusCode < 200 || response.HttpStatusCode >= 300 {
		return nil, fmt.Errorf("wrong response code:%d", response.HttpStatusCode)
	}
	if *response.Metadata.TotalCount == 0 {
		return nil, fmt.Errorf("total count is 0")
	}
	return response.Zones, nil
}

func (h *HuaweiCloudApi) checkSetZoneId(targetDomain string) error {
	targetDomain = dns2.Fqdn(targetDomain)
	if h.targetDomain == targetDomain {
		return nil
	}
	zones, err := h.listPublicZones()
	if err != nil {
		return err
	}
	zoneId, err := h.findZoneIdByDomain(targetDomain, zones)
	if err != nil {
		return err
	}
	h.targetDomain = targetDomain
	h.zoneId = zoneId
	return nil
}

// getExistRecordId 返回已有记录ID，若不存在则返回空字符串
func (h *HuaweiCloudApi) getExistRecordId(ipType model.IPType) (string, error) {
	recordType := ipType.RecordType()
	request := &model2.ListRecordSetsByZoneRequest{
		ZoneId: h.zoneId,
		Name:   &h.targetDomain,
		Type:   &recordType,
	}
	response, err := h.client.ListRecordSetsByZone(request)
	if err != nil {
		return "", err
	}
	if response.HttpStatusCode < 200 || response.HttpStatusCode >= 300 {
		return "", fmt.Errorf("wrong response code:%d", response.HttpStatusCode)
	}
	if *response.Metadata.TotalCount == 0 {
		return "", nil
	}
	return *(*response.Recordsets)[0].Id, nil
}

func (h *HuaweiCloudApi) createRecord(ipList []string, ipType model.IPType) error {
	ttl := int32(1)
	request := &model2.CreateRecordSetRequest{ZoneId: h.zoneId}
	request.Body = &model2.CreateRecordSetRequestBody{
		Name:    h.targetDomain,
		Type:    ipType.RecordType(),
		Records: ipList,
		Ttl:     &ttl,
	}
	response, err := h.client.CreateRecordSet(request)
	if err != nil {
		return err
	}
	if response.HttpStatusCode < 200 || response.HttpStatusCode >= 300 {
		return fmt.Errorf("wrong response code:%d", response.HttpStatusCode)
	}
	return nil
}
func (h *HuaweiCloudApi) updateRecord(recordId string, ipList []string, ipType model.IPType) error {
	request := &model2.UpdateRecordSetRequest{
		ZoneId:      h.zoneId,
		RecordsetId: recordId,
	}
	ttlUpdateRecordSetReq := int32(1)
	typeUpdateRecordSetReq := ipType.RecordType()
	request.Body = &model2.UpdateRecordSetReq{
		Name:    &h.targetDomain,
		Ttl:     &ttlUpdateRecordSetReq,
		Type:    &typeUpdateRecordSetReq,
		Records: &ipList,
	}
	response, err := h.client.UpdateRecordSet(request)
	if err != nil {
		return err
	}
	if response.HttpStatusCode < 200 || response.HttpStatusCode >= 300 {
		return fmt.Errorf("wrong response code:%d", response.HttpStatusCode)
	}
	return nil
}
func (h *HuaweiCloudApi) SetRecord(ipList []net.IP, targetDomain string, ipType model.IPType) error {
	if err := h.checkSetZoneId(targetDomain); err != nil {
		return fmt.Errorf("checkZoneId error, %v", err)
	}
	recordId, err := h.getExistRecordId(ipType)
	if err != nil {
		return fmt.Errorf("getExistRecordId error, %v", err)
	}
	ipStrList := util.ConvertIPListToStrings(ipList)
	if recordId != "" {
		err = h.updateRecord(recordId, ipStrList, ipType)
		if err != nil {
			return fmt.Errorf("update record error, %v", err)
		}
	} else {
		err = h.createRecord(ipStrList, ipType)
		if err != nil {
			return fmt.Errorf("create record error, %v", err)
		}
	}
	return nil
}

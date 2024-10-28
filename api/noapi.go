package api

import (
	"dns-shift/model"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

// NoApi 无API，写入本地文件
type NoApi struct {
}

func (n *NoApi) SetRecord(ipList []net.IP, targetDomain string, ipType model.IPType) error {
	// 构建 JSON 数据结构
	record := struct {
		Domain    string   `json:"domain"`
		Type      string   `json:"type"`
		Records   []string `json:"records"`
		Timestamp string   `json:"timestamp"`
	}{
		Domain:    targetDomain,
		Type:      ipType.RecordType(),
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// 将 IP 列表转换为字符串形式
	for _, ip := range ipList {
		record.Records = append(record.Records, ip.String())
	}

	// 转换为 JSON
	jsonData, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	// 将 JSON 数据写入文件
	filename := fmt.Sprintf("%s_record_%s.json", targetDomain, ipType.RecordType())
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed to write JSON to file: %v", err)
	}

	// 在日志中输出保存的 JSON 文件信息
	log.Printf("Saved DNS record for %s (%s) to %s:\n%s\n", targetDomain, ipType.RecordType(), filename, string(jsonData))

	return nil
}

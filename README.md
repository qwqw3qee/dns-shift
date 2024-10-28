# dns-shift
## 描述

根据源域名的 A 或 AAAA 记录，将 IP 进行偏移，通过域名服务商的API将其应用到目标域名上。

## 快速上手

### 直接运行

1. 从[Releases](https://github.com/qwqw3qee/dns-shift/releases)页面获取最新可执行程序。
2. 编辑`dns-shift.yaml`，修改为你的域名与服务商接口信息，参考[配置文件参数说明](./doc/配置文件参数说明.md) 。
2. 设置可执行权限`chmod +x dns-shift`
3. `./dns-shift`运行。

### 定时执行（仅Linux）

**设置执行计划**

```bash
chmod +x setup-cron.sh
./setup-cron.sh
```

根据提示输入时间间隔及执行时间

**删除执行计划**

```bash
./setup-cron.sh rm
```

## 功能概述

- **多 DNS 服务商支持**：可根据不同 DNS 服务商接口灵活扩展，目前仅支持华为云。
- **IPv4 与 IPv6 支持**：根据配置自动选择 IPv4 或 IPv6 地址类型，支持灵活的地址前缀与长度设置。
- **域名解析记录同步**：自动检测并同步源域名和目标域名的解析记录，支持 Ping 检测功能。
- **DNS 查询与更新**：可配置 DNS 查询服务器地址，支持快速的域名记录查询。

## TODO

- [ ] 多 DNS 服务商支持

## 开发

克隆项目仓库：

```bash
git clone https://github.com/qwqw3qee/dns-shift
cd dns-shift
```

安装依赖并构建项目：

```bash
go mod tidy
go build -o dns-shift
```

### 添加更多API

参考： [添加其它API开发说明](./doc/添加其它API开发说明.md) 。

## 用到的开源项目

- [spf13/viper: Go configuration with fangs](https://github.com/spf13/viper)
- [miekg/dns: DNS library in Go](https://github.com/miekg/dns)
- [prometheus-community/pro-bing: A library for creating continuous probers](https://github.com/prometheus-community/pro-bing)
- [huaweicloud/huaweicloud-sdk-go-v3](https://github.com/huaweicloud/huaweicloud-sdk-go-v3)

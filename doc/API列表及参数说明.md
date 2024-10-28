# API列表及参数说明



## 无服务商

```yaml
dnsServer:
  provider: "noapi"
```

该选项为默认配置，在未配置服务商API时，通过该选项，可以将程序生成的IP信息保存至本地，方便调试或者其它程序使用。

## 华为云

```yaml
dnsServer:                                   # DNS 服务提供商相关配置
  provider: "huaweicloud"                    # DNS 服务提供商的名称，目前支持的有 "huaweicloud"
  apiParam:                                  # API 参数，根据不同的服务提供商可能会有所不同
    ak: "your-ak"                            # API 密钥（Access Key），用于身份验证
    sk: "your-sk"                            # API 密钥（Secret Key），用于身份验证
```

- 要使用华为云 Go SDK ，您需要拥有云账号以及该账号对应的 Access Key（AK）和 Secret Access Key（SK）。请在华为云控制台“我的凭证-访问密钥”页面上创建和查看您的 AK&SK 。更多信息请查看 [访问密钥](https://support.huaweicloud.com/usermanual-ca/zh-cn_topic_0046606340.html) 。

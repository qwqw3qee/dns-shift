# 添加其它API开发说明

1. 假设你的API的配置文件如下。

```yaml
dnsServer: 
  provider: "custom"
  apiParam:
    customKey1: "customValue1"
    customKey2: "customValue2"
```

2. 在[`model/dnstype.go`](../model/dnstype.go)中，添加你的API标识。

```go
const (
    // 新的API标识
    Custom       DnsType = "custom"
)
```

3. 在`api/`目录下，新建一个文件，名称任意，如`custom.go`，编写API的实现代码，并实现[`api/interface.go`](../api/interface.go)中的接口。

```go
package api

import (
    "dns-shift/model"
    "net"
)

// CustomApi 自定义Api的实现
type CustomApi struct {
    // `dnsServer.apiParam`中的参数以map[string]interface{}的形式传入params变量中，以便在使用API时进行调用
    params model.ConfigMap
    // 可以将全局用到的变量均定义于此
}

// SetRecord 设置DNS解析记录主流程
func (c *CustomApi) SetRecord(ipList []net.IP, targetDomain string, ipType model.IPType) error {
    //TODO implement me
    customKey1 := c.params.GetString("customKey1")
    customKey2 := c.params.GetString("customKey2")
    // 设置DNS记录主逻辑
    return nil
}

// NewCustomApi 返回一个实现了DnsApi接口的对象，供RegisterApi()使用
func NewCustomApi(params model.ConfigMap) DnsApi {
    api := &CustomApi{params: params}
    return api
}

```

4. 在[`api/interface.go`](../api/interface.go)中，找到`init()`函数，添加如下内容，注册你的API。

```go
// 初始化注册各个 API 实现
func init() {
    // ...
    // 添加CustomApi注册
    RegisterApi(model.Custom, NewCustomApi)
}
```

完成上述操作后，便可通过配置文件，调用新的API了。具体细节可参考[`api/huaweicloud.go`](../api/huaweicloud.go)。


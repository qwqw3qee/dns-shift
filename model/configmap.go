package model

// ConfigMap 是对 map[string]interface{} 的封装
type ConfigMap map[string]interface{}

// Set 设置一个键值对
func (c ConfigMap) Set(key string, value interface{}) {
	c[key] = value
}

// Get 获取键对应的值，如果不存在则返回nil
func (c ConfigMap) Get(key string) interface{} {
	return c[key]
}

// GetString 直接获取字符串值
func (c ConfigMap) GetString(key string) string {
	if value, ok := c[key].(string); ok {
		return value
	}
	return ""
}

// GetInt 获取整数值
func (c ConfigMap) GetInt(key string) int {
	if value, ok := c[key].(int); ok {
		return value
	}
	return 0
}

// 其他类似方法可以继续扩展

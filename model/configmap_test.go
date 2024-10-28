package model

import (
	"fmt"
	"testing"
)

func TestConfigMap(t *testing.T) {

	config := ConfigMap{
		"host": "localhost",
		"port": 8080,
	}

	host := config.GetString("host")
	port := config.GetInt("port")
	config.Set("aaa", "bbb")
	fmt.Println("Host:", host)
	fmt.Println("Port:", port)
	fmt.Printf("%v\n", config)

}

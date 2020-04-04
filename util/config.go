package util

import (
	"fmt"
)

type Address struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

var urls []string

func GetConsulUrls() []string {
	addr := make([]Address, 1)
	//"host": "192.168.0.198",
	//	"port": 8500
	addr[0].Host = "192.168.0.198"
	addr[0].Port = 8500  //"192.168.0.198"
	//if err := config.Get("consul").Scan(&addr); err != nil {
	//	log.Panic(err)
	//}

	for _, addr := range addr {
		urls = append(urls, fmt.Sprintf("%s:%d", addr.Host, addr.Port))
	}

	return urls
}
package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
)

func Register(address string, port int, name string, tags []string, id string) {
	var (
		err    error
		client *api.Client
	)
	config := &api.Config{
		Address: fmt.Sprintf("%s:%d", address, port),
	}
	if client, err = api.NewClient(config); err != nil {
		log.Panicln(err)
	}
	// 生成注册对象
	registration := &api.AgentServiceRegistration{
		ID:      id,
		Name:    name,
		Tags:    tags,
		Port:    8080,
		Address: "172.16.130.1",
		Check: &api.AgentServiceCheck{
			HTTP:                           "http://172.16.130.1:8080/health",
			Timeout:                        "5s",
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "10s",
		},
	}
	if err = client.Agent().ServiceRegister(registration); err != nil {
		log.Fatalln(err)
	}
	/*
		data, err := client.Agent().Services()
		for k, v := range data {
			fmt.Println(k, v)
		}

	*/
	data, err := client.Agent().ServicesWithFilter(`Service == "shop-web"`)
	for k, v := range data {
		fmt.Println(k, v)
	}
}

func main() {
	Register("127.0.0.1", 8500, "shop-web", []string{"web", "local", "web-shop"}, "shop-web")
}

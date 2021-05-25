package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type ProxyConfig struct {
	Namespace string  `yaml:"namespace,omitempty"`
	Service map[string][]struct{
		Port       string `yaml:"port"`
		TargetPort string `yaml:"targetPort"`
		namespace string `yaml:"namespace",omitempty`
	} `yaml:"service,omitempty"`
}

func main() {

	_, err := exec.Command("which", "kubectl").CombinedOutput()
	if err != nil {
		log.Fatalf("Please install kubectl tool")
	}
	proxyConfig := ProxyConfig{}
	namespace:="default"
	config, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err  #%v ", err)
	}

	err = yaml.Unmarshal(config, &proxyConfig)
	if len(proxyConfig.Namespace) != 0 {
		namespace = proxyConfig.Namespace
	}
	pids := []*os.Process{}
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	for serviceName, service := range proxyConfig.Service {
		for _, elm := range service {
			fmt.Printf(namespace)
			cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf(" kubectl port-forward --namespace %s svc/%s %s:%s",
				namespace,serviceName, service[0].Port, service[0].TargetPort))
			err := cmd.Run()
			if err != nil {
				log.Printf("Cannot proxied %s", serviceName)
			} else {
				log.Printf("%s proxied %s-->%s", serviceName, elm.Port, elm.TargetPort)
				pids = append(pids, cmd.Process)
			}

		}

	}

	for {


	}
}

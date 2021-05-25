package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os/exec"
)

type ProxyConfig struct {
	Namespace string  `yaml:"namespace,omitempty"`
	Service map[string][]struct{
		Port       string `yaml:"port"`
		TargetPort string `yaml:"targetPort"`
		Namespace string `yaml:"namespace,omitempty"`
	} `yaml:"service,omitempty"`
}



type Response struct {
	ServiceName string
	Message string
}
func main() {
	_, err := exec.LookPath("kubectl")
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
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	outputChan := make(chan Response)


	fmt.Println("Global namespace: ",namespace)
	ns := namespace
	for serviceName, service := range proxyConfig.Service {
		for _, elm := range service {
			if len(elm.Namespace) != 0 {
				ns = elm.Namespace
			}
			log.Printf("%s service on namespace %s proxied %s --> %s",serviceName,ns,elm.Port,elm.TargetPort)
			go func(serviceName,port,targetPort,ns string){

				cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf(" kubectl port-forward --namespace %s svc/%s %s:%s",
					ns,serviceName, port, targetPort))

				output,err := cmd.CombinedOutput()
				if err != nil {
					outputChan <- Response{ServiceName: serviceName,Message: string(output)}
				}
			}(serviceName,elm.Port,elm.TargetPort,ns)

		}
	}



	for msg := range outputChan{
		log.Println(msg)

	}
}

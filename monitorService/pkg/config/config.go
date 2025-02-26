package config

import (
	"fmt"
	"log"

	"gopkg.in/ini.v1"
)

var serviceConfig = make(map[string]map[string]string)

// 抽象的函数，用于将指定配置节的内容读取到serviceConfig中
func readConfigSection(cfg *ini.File, sectionName string) {
	section, err := cfg.GetSection(sectionName)
	if err != nil {
		log.Fatal("读取配置节错误: ", err)
	}

	serviceConfig[sectionName] = make(map[string]string)
	for _, value := range section.Keys() {
		serviceConfig[sectionName][value.Name()] = value.Value()
	}
}

func ConfigInit() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatal("ConfigInit Load error: ", err)
	}

	for _, sectionName := range cfg.SectionStrings() {
		readConfigSection(cfg, sectionName)
	}
}

func GetHostAndPort(serverName string) (string, error) {
	if serviceConfig, ok := serviceConfig[serverName]; ok {
		return serviceConfig["Host"] + ":" + serviceConfig["Port"], nil
	} else {
		return "", fmt.Errorf("相关服务未注册")
	}
}

func GetPassword(serverName string) (string, error) {
	if serviceConfig, ok := serviceConfig[serverName]; ok {
		return serviceConfig["Password"], nil
	} else {
		return "", fmt.Errorf("相关服务 Password 未注册")
	}
}

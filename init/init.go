package ssrinit

import (
	config2 "SsrMicroClient/config"
	"fmt"
	"io"
	"net/http"
	"os"
)

// PathExists 判断目录是否存在返回布尔类型
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// Init  <-- init
func Init(configPath string) {
	//判断目录是否存在 不存在则创建
	if !PathExists(configPath) {
		err := os.MkdirAll(configPath, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}

	if !PathExists(configPath + "/shadowsocksr") {
		GetSsrPython(configPath)
	}

	if !PathExists(configPath + "/node.json") {
		if config2.InitJSON(configPath) != nil {
			return
		}
	}

	if !PathExists(configPath + "/SsrMicroConfig.json") {
		if config2.SettingInitJSON(configPath) != nil {
			return
		}
	}

	if !PathExists(configPath + "/cidrBypass.conf") {
		res, err := http.Get("https://raw.githubusercontent.com/Asutorufa/SsrMicroClient/ACL/ssrMicroClientBypass.conf")
		if err != nil {
			panic(err)
		}
		f, err := os.Create(configPath + "/cidrBypass.conf")
		if err != nil {
			panic(err)
		}
		_, _ = io.Copy(f, res.Body)
	}

	if !PathExists(configPath + "/domainBypass.conf") {
		res, err := http.Get("https://raw.githubusercontent.com/Asutorufa/SsrMicroClient/ACL/ssrMicroClientDomainBypass.conf")
		if err != nil {
			panic(err)
		}
		f, err := os.Create(configPath + "/domainBypass.conf")
		if err != nil {
			panic(err)
		}
		_, _ = io.Copy(f, res.Body)
	}

	if !PathExists(configPath + "/domainProxy.conf") {
		//res, err := http.Get("https://raw.githubusercontent.com/Asutorufa/SsrMicroClient/ACL/ssrMicroClientBypass.conf")
		//if err != nil {
		//	panic(err)
		//}
		_, err := os.Create(configPath + "/domainProxy.conf")
		if err != nil {
			panic(err)
		}
		//io.Copy(f, res.Body)
	}

	if !PathExists(configPath + "/discordDomain.conf") {
		_, err := os.Create(configPath + "/discordDomain.conf")
		if err != nil {
			panic(err)
		}
	}

	if !PathExists(configPath + "/SsrMicroClient.png") {
		res, err := http.Get("https://raw.githubusercontent.com/Asutorufa/SsrMicroClient/master/SsrMicroClient.png")
		if err != nil {
			panic(err)
		}
		f, err := os.Create(configPath + "/SsrMicroClient.png")
		if err != nil {
			panic(err)
		}
		_, _ = io.Copy(f, res.Body)
	}
}

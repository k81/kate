package config

import (
	"context"
	"io/ioutil"

	goyaml "gopkg.in/yaml.v2"
)

var (
	configFilePath string
	app            string
	mctx           = context.Background()
)

/*************************************************
Description: 	从yaml文件加载配置
Input:
	filePath		yaml配置文件的路径
	config		目标配置变量地址
Output:
Return:			成功，nil；失败，error
Others:
*************************************************/
func LoadConfig(filePath string, config interface{}) (err error) {
	var content []byte
	if content, err = ioutil.ReadFile(filePath); err != nil {
		return
	}

	if err = goyaml.Unmarshal(content, config); err != nil {
		return
	}
	return
}

/*************************************************
Description: 保存配置变量到yaml文件
Input:
	filePath	 yaml配置文件的路径
	config	 配置变量的地址
Output:
Return:		 成功，nil；错误，error
Others:
*************************************************/
func SaveConfig(filePath string, config interface{}) (err error) {
	var content []byte
	if content, err = goyaml.Marshal(config); err != nil {
		return
	}

	err = ioutil.WriteFile(filePath, content, 0644)
	return
}

/*************************************************
Description: 注册并加载全局配置
Input:
	configList	全局配置Entry列表
Output:
Return:
Others:
*************************************************/
func Init(appName, filePath string) {
	app = appName
	configFilePath = filePath

	Local.init()
	Global.init()
}

func Register(entry *Entry) {
	Global.Add(entry)
}

/*************************************************
Description: 获取key关联的配置数据（经过解析的）
Input:
	name	 配置项的Key
Output:
Return:		 Key关联的配置数据
Others:
*************************************************/
func Get(name string) interface{} {
	return Global.GetEntryData(name)
}

/*************************************************
Description: 将解析后的配置数据绑定在key下面
Input:
	name	 配置项的Key
	v		 解析后的配置数据
Output:
Return:
Others:
*************************************************/
func Set(name string, v interface{}) {
	Global.SetEntryData(name, v)
}

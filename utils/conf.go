package utils

import (
	"github.com/go-ini/ini"
)

/* 获取配置conf文件的参数信息 */
func GetConf(name string) *ini.Section {
	cfg, err := ini.Load("conf/app.ini")
	if err != nil {
		panic(err)
	}
	return cfg.Section(name)
}

package viper_etcd

import "github.com/spf13/viper"

var gEtcdRemoteConfig *EtcdRemoteConfig

func init() {
	gEtcdRemoteConfig = &EtcdRemoteConfig{}
	viper.RemoteConfig = gEtcdRemoteConfig
}

func SetUsername(username string) {
	gEtcdRemoteConfig.username = username
}

func SetPassword(password string) {
	gEtcdRemoteConfig.password = password
}

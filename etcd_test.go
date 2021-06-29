package viper_etcd

import (
	"testing"

	"github.com/spf13/viper"
)

func TestEtcd(t *testing.T) {
	endpoint := "192.168.3.132:2379"
	path := "/vipet-etcd.yaml"

	v := viper.New()
	v.AddRemoteProvider("etcd", endpoint, path)
	v.SetConfigType("yaml")
	v.ReadRemoteConfig()
	t.Log(v.Get("test"))
}

func TestEtcdWatch(t *testing.T) {
	endpoint := "192.168.3.132:2379"
	path := "/vipet-etcd.yaml"

	v := viper.New()
	v.AddRemoteProvider("etcd", endpoint, path)
	v.SetConfigType("yaml")
	v.ReadRemoteConfig()
	t.Log(v.Get("test"))

	t.Run(
		"watch",
		func(t *testing.T) {
			v.WatchRemoteConfig()
			t.Log(v.Get("test"))
		},
	)
}

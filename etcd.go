package viper_etcd

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/spf13/viper"
	etcd_client "go.etcd.io/etcd/client/v3"
)

type EtcdRemoteConfig struct {
	username string
	password string
}

func (c *EtcdRemoteConfig) Get(rp viper.RemoteProvider) (io.Reader, error) {
	client, err := c.newClient(rp)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := client.Get(ctx, rp.Path())
	cancel()
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(resp.Kvs[0].Value), nil
}

func (c *EtcdRemoteConfig) Watch(rp viper.RemoteProvider) (io.Reader, error) {
	client, err := c.newClient(rp)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ch := client.Watch(context.Background(), rp.Path())
	resp := <-ch
	return bytes.NewReader(resp.Events[0].Kv.Value), nil
}

func (c *EtcdRemoteConfig) WatchChannel(rp viper.RemoteProvider) (<-chan *viper.RemoteResponse, chan bool) {
	rr := make(chan *viper.RemoteResponse)
	stop := make(chan bool)

	go func() {
		for {
			client, err := c.newClient(rp)
			if err != nil {
				time.Sleep(time.Duration(10 * time.Second))
				continue
			}
			defer client.Close()

			ch := client.Watch(context.Background(), rp.Path())
			select {
			case <-stop:
				return
			case res := <-ch:
				for _, event := range res.Events {
					rr <- &viper.RemoteResponse{
						Value: event.Kv.Value,
					}
				}
			}
		}
	}()

	return rr, stop
}

func (c *EtcdRemoteConfig) newClient(rp viper.RemoteProvider) (*etcd_client.Client, error) {
	client, err := etcd_client.New(
		etcd_client.Config{
			Endpoints: []string{rp.Endpoint()},
			Username:  c.username,
			Password:  c.password,
		},
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}

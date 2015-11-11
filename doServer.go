package main

import (
	"fmt"
	"log"
	"path"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

func doServer(kapi client.KeysAPI) {
	var key = configbase + *servername

	var settings map[string]string
	settings = make(map[string]string)

	resp, err := kapi.Get(context.TODO(), key, &client.GetOptions{Recursive: true})
	if err != nil {
		log.Fatal(err)
	}

	for _, node := range resp.Node.Nodes {
		_, setting := path.Split(node.Key)
		settings[setting] = node.Value
	}

	fmt.Println(settings)

	watcher := kapi.Watcher(key, &client.WatcherOptions{Recursive: true})

	for true {
		resp, err := watcher.Next(context.TODO())

		if err != nil {
			if _, ok := err.(*client.ClusterError); ok {
				continue
			}
			log.Fatal(err)
		}

		switch resp.Action {
		case "set":
			_, setting := path.Split(resp.Node.Key)
			settings[setting] = resp.Node.Value
		case "delete", "expire":
			_, setting := path.Split(resp.Node.Key)
			delete(settings, setting)
		}

		fmt.Println(settings)
	}
}

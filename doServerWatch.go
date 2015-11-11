package main

import (
	"fmt"
	"log"
	"path"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

func doServerWatch(kapi client.KeysAPI) {

	watcher := kapi.Watcher(runningbase, &client.WatcherOptions{Recursive: true})

	for true {
		resp, err := watcher.Next(context.TODO())

		if err != nil {
			if _, ok := err.(*client.ClusterError); ok {
				continue
			}
			log.Fatal(err)
		}

		fmt.Println(resp.Node.Key + " " + resp.Node.Value)

		_, server := path.Split(resp.Node.Key)
		switch resp.Action {
		case "create":
			fmt.Println(server + " has started heart beat")
		case "compareAndSwap":
			fmt.Println(server + " heart beat")
		case "compareAndDelete":
			fmt.Println(server + " has shut down correctly")
		case "expire":
			fmt.Println("*** " + server + " has missed heartbeat")
		default:
			fmt.Println("Didn't handle " + resp.Action)
		}
	}

}

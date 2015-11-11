package main

import (
	"fmt"
	"log"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

func doDumpQueue(kapi client.KeysAPI) {
	var key = queuebase + *dumpqueuename

	resp, err := kapi.Get(context.TODO(), key, &client.GetOptions{Sort: true})

	if err != nil {
		log.Fatal(err)
	}
	for _, v := range resp.Node.Nodes {
		fmt.Println(v.Key + " set to " + v.Value)
	}
}

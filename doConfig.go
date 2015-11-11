package main

import (
	"fmt"
	"log"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

func doConfig(kapi client.KeysAPI) {
	var key = configbase + *configserver + "/" + *configvar

	resp, err := kapi.Set(context.TODO(), key, *configval, nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Action + " " + resp.Node.Key + " to " + resp.Node.Value)
}

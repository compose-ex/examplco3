package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

func doFillQueue(kapi client.KeysAPI) {
	var key = queuebase + *queuename
	list := rand.Perm(10)
	for _, v := range list {
		value := "Value" + strconv.Itoa(v)

		resp, err := kapi.CreateInOrder(context.TODO(), key, value, nil)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(resp.Action + " " + resp.Node.Key + " to " + resp.Node.Value)
	}
}

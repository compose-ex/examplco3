package main

import (
	"fmt"
	"log"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/twinj/uuid"
	"golang.org/x/net/context"
)

func doServerBeat(kapi client.KeysAPI) {
	var key = runningbase + *serverbeatname

	myuuid := uuid.NewV4()
	uuidstring := myuuid.String()

	fmt.Println("Badum")
	_, err := kapi.Set(context.TODO(), key, uuidstring, &client.SetOptions{PrevExist: client.PrevNoExist, TTL: time.Second * 60})
	if err != nil {
		log.Fatal(err)
	}

	running := true
	counter := *serverbeatcount

	for running {
		time.Sleep(time.Second * time.Duration(*serverbeattime))
		fmt.Println("Badum")
		_, err := kapi.Set(context.TODO(), key, uuidstring, &client.SetOptions{PrevExist: client.PrevExist, TTL: time.Second * 60, PrevValue: uuidstring})
		if err != nil {
			log.Fatal(err)
		}
		if *serverbeatcount != 0 {
			counter = counter - 1
			if counter == 0 {
				running = false
			}
		}
	}

	_, err = kapi.Delete(context.TODO(), key, &client.DeleteOptions{PrevValue: uuidstring})
	if err != nil {
		log.Fatal(err)
	}
}

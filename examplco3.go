package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/coreos/etcd/client"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app      = kingpin.New("examplco2", "An extended etcd demonstration")
	peerlist = app.Flag("peers", "etcd peers").Default("http://127.0.0.1:4001,http://127.0.0.1:2379").OverrideDefaultFromEnvar("EX_PEERS").String()
	username = app.Flag("user", "etcd User").OverrideDefaultFromEnvar("EX_USER").String()
	password = app.Flag("pass", "etcd Password").OverrideDefaultFromEnvar("EX_PASS").String()
	cafile   = app.Flag("cacert", "CA Certificate").OverrideDefaultFromEnvar("EX_CERT").String()

	config       = app.Command("config", "Change config data")
	configserver = config.Arg("server", "Server name").Required().String()
	configvar    = config.Arg("var", "Config variable").Required().String()
	configval    = config.Arg("val", "Config value").Required().String()

	server     = app.Command("server", "Go into server mode and listen for changes")
	servername = server.Arg("server", "Server name").Required().String()

	serverbeat      = app.Command("serverbeat", "Go into server mode and heartbeat in etcd")
	serverbeatname  = serverbeat.Arg("server", "Server name").Required().String()
	serverbeattime  = serverbeat.Flag("rate", "Time between beats").Default("60").Int()
	serverbeatcount = serverbeat.Flag("count", "Number of beats - 0 forever").Default("0").Int()

	serverwatch = app.Command("serverwatch", "Watch for expiring servers")

	fillqueue = app.Command("fillqueue", "Fill an ordered named queue with values")
	queuename = fillqueue.Arg("queue", "Queue name").Default("jobqueue").String()

	dumpqueue     = app.Command("dumpqueue", "Print an sorted named queue")
	dumpqueuename = dumpqueue.Arg("queue", "Queue name").Default("jobqueue").String()
)

var configbase = "/config/"
var runningbase = "/running/"
var queuebase = "/queues/"

func main() {
	kingpin.Version("0.0.3")
	command := kingpin.MustParse(app.Parse(os.Args[1:]))

	peers := strings.Split(*peerlist, ",")

	// Read the certificate into a file
	caCert, err := ioutil.ReadFile(*cafile)
	if err != nil {
		log.Fatal(err)
	}

	// Create a certificate pool
	caCertPool := x509.NewCertPool()
	// and add the freshly read certificate to the new pool
	caCertPool.AppendCertsFromPEM(caCert)

	// Create a TLS configuration structure
	// with the certificate pool as it's list of certificate authorities
	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}

	// Then create a HTTP transport with that configuration
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	// When we create the etcd client configuration, use that transport
	cfg := client.Config{
		Endpoints:               peers,
		Transport:               transport,
		HeaderTimeoutPerRequest: time.Minute,
		Username:                *username,
		Password:                *password,
	}

	// And create your client as normal.
	etcdclient, err := client.New(cfg)

	if err != nil {
		log.Fatal(err)
	}

	kapi := client.NewKeysAPI(etcdclient)

	switch command {
	case config.FullCommand():
		doConfig(kapi)
	case server.FullCommand():
		doServer(kapi)
	case serverbeat.FullCommand():
		doServerBeat(kapi)
	case serverwatch.FullCommand():
		doServerWatch(kapi)
	case fillqueue.FullCommand():
		doFillQueue(kapi)
	case dumpqueue.FullCommand():
		doDumpQueue(kapi)
	}
}

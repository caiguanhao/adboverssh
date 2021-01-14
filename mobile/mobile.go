package adboverssh

import (
	"log"
	"time"

	"github.com/caiguanhao/adboverssh"
)

var (
	adbAddress, listen, sshAddr, sshUser, privKey string
)

func Run() {
	Start(adbAddress, listen, sshAddr, sshUser, privKey)
}

func Start(adbAddress, listen, sshAddr, sshUser, privKey string) {
	c := &adboverssh.Client{
		ADBAddress:              adbAddress,
		SSHListenAddress:        listen,
		SSHServerAddress:        sshAddr,
		SSHServerUser:           sshUser,
		SSHServerUserPrivateKey: []byte(privKey),
		SSHConnectTimeout:       5 * time.Second,
		OnConnected: func(addr string) {
			log.Println("connected to", addr)
		},
		OnListening: func(addr string) {
			log.Println("listening", addr)
		},
		OnNewConnection: func(a, b string) {
			log.Println("connected", a, "<->", b)
		},
		OnError: func(err error) {
			log.Println(err)
		},
	}
	for {
		c.Connect()
		time.Sleep(2 * time.Second)
	}
}

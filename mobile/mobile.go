package adboverssh

import (
	"log"
	"time"

	"github.com/caiguanhao/adboverssh"
)

var (
	adbAddress, sshListen, sshAddress, sshUser, sshPrivKey string
)

func Run() {
	Start(adbAddress, sshListen, sshAddress, sshUser, sshPrivKey)
}

func Start(adbAddress, sshListen, sshAddress, sshUser, sshPrivKey string) {
	c := &adboverssh.Client{
		ADBAddress:              adbAddress,
		SSHListenAddress:        sshListen,
		SSHServerAddress:        sshAddress,
		SSHServerUser:           sshUser,
		SSHServerUserPrivateKey: []byte(sshPrivKey),
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

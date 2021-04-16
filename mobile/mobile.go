package adboverssh

import (
	"log"
	"time"

	"github.com/caiguanhao/adboverssh"
)

var (
	adbAddress, sshListen, sshAddress, sshUser, sshPrivKey string

	currentAddress string

	currentClient *adboverssh.Client
)

func GetCurrentAddress() string {
	return currentAddress
}

func StartDefaultClient() {
	StartClient(adbAddress, sshListen, sshAddress, sshUser, sshPrivKey)
}

func StartClient(adbAddress, sshListen, sshAddress, sshUser, sshPrivKey string) {
	currentClient = newClient(adbAddress, sshListen, sshAddress, sshUser, sshPrivKey)
	currentClient.Connect()
}

func Stop() {
	if currentClient == nil {
		return
	}
	currentClient.Stop()
	currentAddress = ""
	log.Println("stopped")
}

func newClient(adbAddress, sshListen, sshAddress, sshUser, sshPrivKey string) *adboverssh.Client {
	return &adboverssh.Client{
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
			currentAddress = addr
			log.Println("listening", addr)
		},
		OnNewConnection: func(a, b string) {
			log.Println("connected", a, "<->", b)
		},
		OnError: func(err error) {
			log.Println(err)
		},
	}
}

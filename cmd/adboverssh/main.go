package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/caiguanhao/adboverssh"
)

func main() {
	adbAddress := flag.String("adb", "127.0.0.1:5555", "adb address")
	listen := flag.String("l", "127.0.0.1:0", "listen address on ssh server")
	keyFile := flag.String("i", "", "ssh private key")
	timeout := flag.Int("t", 5, "ssh connection timeout in seconds")
	flag.Parse()

	sshStr := strings.SplitN(flag.Arg(0), "@", 2)
	if len(sshStr) != 2 {
		log.Fatalln("you must provide user@ssh-server-address")
	}
	sshUser, sshAddr := sshStr[0], sshStr[1]

	if *keyFile == "" {
		log.Fatalln("you must provide private key")
	}
	privKey, err := ioutil.ReadFile(*keyFile)
	if err != nil {
		log.Fatalln(err)
	}

	c := &adboverssh.Client{
		ADBAddress:              *adbAddress,
		SSHListenAddress:        *listen,
		SSHServerAddress:        sshAddr,
		SSHServerUser:           sshUser,
		SSHServerUserPrivateKey: privKey,
		SSHConnectTimeout:       time.Duration(*timeout) * time.Second,
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

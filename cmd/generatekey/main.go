package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	adbAddress := os.Getenv("ADB_ADDRESS")
	if adbAddress == "" {
		reader := bufio.NewReader(os.Stdin)
		defaultAdbAddress := "127.0.0.1:5555"
		fmt.Print("Enter adb address (" + defaultAdbAddress + "): ")
		adbAddress, _ = reader.ReadString('\n')
		adbAddress = strings.TrimSpace(adbAddress)
		if adbAddress == "" {
			adbAddress = defaultAdbAddress
		}
	}

	listen := os.Getenv("LISTEN")
	if listen == "" {
		reader := bufio.NewReader(os.Stdin)
		defaultListen := "127.0.0.1:0"
		fmt.Print("Enter listen address (" + defaultListen + "): ")
		listen, _ = reader.ReadString('\n')
		listen = strings.TrimSpace(listen)
		if listen == "" {
			listen = defaultListen
		}
	}

	sshAddress := os.Getenv("SSH_ADDRESS")
	for !strings.Contains(sshAddress, ":") {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter ssh address (host/ip:port): ")
		sshAddress, _ = reader.ReadString('\n')
		sshAddress = strings.TrimSpace(sshAddress)
	}

	sshUser := os.Getenv("SSH_USER")
	if sshUser == "" {
		reader := bufio.NewReader(os.Stdin)
		defaultSshUser := "root"
		fmt.Print("Enter ssh user (" + defaultSshUser + "): ")
		sshUser, _ = reader.ReadString('\n')
		sshUser = strings.TrimSpace(sshUser)
		if sshUser == "" {
			sshUser = defaultSshUser
		}
	}

	sshPrivKey := os.Getenv("SSH_PRIVATE_KEY")
	sshPrivKeyContent := readFile(sshPrivKey)
	for sshPrivKeyContent == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter ssh private key file location: ")
		sshPrivKey, _ = reader.ReadString('\n')
		sshPrivKey = strings.TrimSpace(sshPrivKey)
		sshPrivKeyContent = readFile(sshPrivKey)
	}

	content := `package adboverssh

import (
	"strings"
)

func init() {
	adbAddress = ` + toJson(adbAddress) + `

	sshListen = ` + toJson(listen) + `

	sshAddress = ` + toJson(sshAddress) + `

	sshUser = ` + toJson(sshUser) + `

	sshPrivKey = strings.Join([]string{`
	for i, c := range sshPrivKeyContent {
		if i%8 == 0 {
			content += "\n\t\t"
		} else if i > 0 {
			content += " "
		}
		content += fmt.Sprintf(`"\x%02X",`, c)
	}
	content += `
	}, "")
}
`
	writeFile("mobile/key.go", content)
}

func toJson(i string) string {
	b, _ := json.Marshal(i)
	return string(b)
}

func readFile(file string) string {
	c, _ := ioutil.ReadFile(file)
	return string(c)
}

func writeFile(file, content string) {
	err := ioutil.WriteFile(file, []byte(content), 0644)
	if err != nil {
		panic(err)
	}
	log.Println("written", file)
}

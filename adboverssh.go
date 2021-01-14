package adboverssh

import (
	"errors"
	"io"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

type (
	Client struct {
		ADBAddress              string
		SSHListenAddress        string
		SSHServerAddress        string
		SSHServerUser           string
		SSHServerUserPrivateKey []byte
		SSHConnectTimeout       time.Duration

		OnConnected     func(string)
		OnListening     func(string)
		OnNewConnection func(string, string)
		OnError         func(error)
	}

	sshDialResult struct {
		client *ssh.Client
		err    error
	}
)

func (c *Client) Connect() error {
	conn, err := c.dial()
	if err != nil {
		if c.OnError != nil {
			c.OnError(err)
		}
		return err
	}
	defer conn.Close()

	if c.OnConnected != nil {
		c.OnConnected(c.SSHServerAddress)
	}

	listener, err := conn.Listen("tcp", c.SSHListenAddress)
	if err != nil {
		if c.OnError != nil {
			c.OnError(err)
		}
		return err
	}
	defer listener.Close()

	if c.OnListening != nil {
		c.OnListening(listener.Addr().String())
	}

	for {
		if err := c.accept(listener); err != nil {
			if c.OnError != nil {
				c.OnError(err)
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func (c *Client) accept(listener net.Listener) error {
	client, err := listener.Accept()
	if err != nil {
		return err
	}
	defer client.Close()

	adbInAndroid, err := net.Dial("tcp", c.ADBAddress)
	if err != nil {
		return err
	}
	defer adbInAndroid.Close()

	if c.OnNewConnection != nil {
		c.OnNewConnection(client.RemoteAddr().String(), adbInAndroid.RemoteAddr().String())
	}

	go io.Copy(client, adbInAndroid)
	io.Copy(adbInAndroid, client)
	return nil
}

func (c *Client) dial() (*ssh.Client, error) {
	ch := make(chan sshDialResult)

	go func() {
		signer, err := ssh.ParsePrivateKey(c.SSHServerUserPrivateKey)
		if err != nil {
			ch <- sshDialResult{nil, err}
			return
		}
		client, err := ssh.Dial("tcp", c.SSHServerAddress, &ssh.ClientConfig{
			User: c.SSHServerUser,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		})
		ch <- sshDialResult{client, err}
	}()

	if c.SSHConnectTimeout < 1*time.Second {
		result := <-ch
		return result.client, result.err
	}

	select {
	case result := <-ch:
		return result.client, result.err
	case <-time.After(c.SSHConnectTimeout):
		return nil, errors.New("timed out")
	}
}

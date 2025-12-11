package main

import (
	"io"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/ssh"
)

type TunnelConfig struct {
	SSHUser       string
	SSHPassword   string
	SSHServerAddr string
	RemoteAddr    string
	LocalAddr     string
}

func pipe(write, render net.Conn) {
	defer write.Close()
	defer render.Close()
	_, err := io.Copy(write, render)
	if err != nil {
		log.Printf("failed to copy: %s", err)
	}
}

func reverseTunnel(conn *ssh.Client, remoteAddr, localAddr string) error {
	listener, err := net.Listen("tcp", remoteAddr)
	if err != nil {
		return err
	}

	log.Printf("SSH server listening on %s -> forwarding to local %s", remoteAddr, localAddr)

	for {
		remoteConn, err := listener.Accept()
		if err != nil {
			return err
		}

		localConn, err := conn.Dial("tcp", localAddr)
		if err != nil {
			log.Fatalf("failed to dial remote: %q", err)
			remoteConn.Close()
			continue
		}
		go pipe(remoteConn, localConn)
		go pipe(localConn, remoteConn)
	}
}

func loadConfig() *TunnelConfig {
	godotenv.Load()
	sshUser := os.Getenv("SSH_USER")
	sshPassword := os.Getenv("SSH_PASSWORD")
	sshServerAddr := os.Getenv("SSH_SERVER_ADDR")
	remoteAddr := os.Getenv("REMOTE_ADDR")
	localAddr := os.Getenv("LOCAL_ADDR")
	if sshUser == "" {
		sshUser = "root"
	}
	if sshPassword == "" {
		sshPassword = "123456"
	}
	if sshServerAddr == "" {
		sshServerAddr = "example.com:22"
	}
	if remoteAddr == "" {
		remoteAddr = "0.0.0.0:9001"
	}
	if localAddr == "" {
		localAddr = "127.0.0.1:9091"
	}
	return &TunnelConfig{
		SSHUser:       sshUser,
		SSHPassword:   sshPassword,
		SSHServerAddr: sshServerAddr,
		RemoteAddr:    remoteAddr,
		LocalAddr:     localAddr,
	}
}

func main() {
	tunnelConfig := loadConfig()
	config := &ssh.ClientConfig{
		User: tunnelConfig.SSHUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(tunnelConfig.SSHPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	serverAddr := tunnelConfig.SSHServerAddr
	conn, err := ssh.Dial("tcp", serverAddr, config)
	if err != nil {
		log.Fatalf("failed to dial: %s", err)
	}
	defer conn.Close()

	remoteExpose := tunnelConfig.RemoteAddr
	localService := tunnelConfig.LocalAddr
	if err := reverseTunnel(conn, remoteExpose, localService); err != nil {
		log.Fatalf("tunnel failed: %s", err)
	}
}

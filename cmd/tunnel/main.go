package main

import (
	"log"

	"golang.org/x/crypto/ssh"

	"github.com/atrox39/ReverseReach/internal/config"
	"github.com/atrox39/ReverseReach/internal/tunnel"
	"github.com/atrox39/ReverseReach/internal/web"
)

func main() {
	cfg := config.Load()

	sshConfig := &ssh.ClientConfig{
		User: cfg.SSHUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(cfg.SSHPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", cfg.SSHServerAddr, sshConfig)
	if err != nil {
		log.Fatalf("failed to dial ssh: %s", err)
	}
	defer conn.Close()

	go web.Start(cfg.WebPort)

	err = tunnel.StartReverse(conn, cfg.RemoteAddr, cfg.LocalAddr)
	if err != nil {
		log.Fatalf("tunnel failed: %s", err)
	}
}

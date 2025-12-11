package tunnel

import (
	"io"
	"log"
	"net"

	"golang.org/x/crypto/ssh"

	"github.com/atrox39/ReverseReach/internal/metrics"
)

func pipe(write, read net.Conn) {
	defer write.Close()
	defer read.Close()

	buf := make([]byte, 32*1024)

	for {
		n, err := read.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("pipe error: %s", err)
			}
			return
		}

		metrics.AddBytesReceived(uint64(n))

		_, err = write.Write(buf[:n])
		if err != nil {
			log.Printf("pipe write error: %s", err)
			return
		}

		metrics.AddBytesSent(uint64(n))
	}
}

func StartReverse(conn *ssh.Client, remoteAddr, localAddr string) error {
	listener, err := conn.Listen("tcp", remoteAddr)
	if err != nil {
		return err
	}

	log.Printf("SSH Tunnel listening on %s -> forwarding to %s", remoteAddr, localAddr)

	for {
		remoteConn, err := listener.Accept()
		if err != nil {
			log.Printf("accept error: %s", err)
			continue
		}

		metrics.AddConnection()

		localConn, err := net.Dial("tcp", localAddr)
		if err != nil {
			log.Printf("local dial error: %s", err)
			remoteConn.Close()
			continue
		}

		go pipe(remoteConn, localConn)
		go pipe(localConn, remoteConn)
	}
}

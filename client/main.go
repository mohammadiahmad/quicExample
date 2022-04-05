package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/lucas-clemente/quic-go"
	"io"
	"time"
)

const addr = "localhost:4242"

const message = "foobar"

// We start a server echoing data on the first stream the client opens,
// then connect with a client, send the message, and wait for its receipt.
func main() {

	err := clientMain()
	if err != nil {
		panic(err)
	}
	select {}
}

func clientMain() error {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}
	conn, err := quic.DialAddr(addr, tlsConf, nil)
	if err != nil {
		return err
	}

	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		return err
	}
	for {

		fmt.Printf("Client: Sending '%s'\n", message)
		_, err = stream.Write([]byte(message))
		if err != nil {
			return err
		}

		buf := make([]byte, len(message))
		_, err = io.ReadFull(stream, buf)
		if err != nil {
			return err
		}
		fmt.Printf("Client: Got '%s'\n", buf)
		time.Sleep(1 * time.Second)
	}
	return nil
}

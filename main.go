package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

func main() {
	proxy := os.Args[1]
	connectto := os.Args[2]

	tr := &http3.Transport{
		TLSClientConfig: &tls.Config{
			NextProtos: []string{http3.NextProtoH3},
		},
	}

	conn := tr.NewClientConn(must(quic.DialAddr(context.Background(), proxy, &tls.Config{
		NextProtos: []string{http3.NextProtoH3},
	}, &quic.Config{})))

	str := must(conn.OpenRequestStream(conn.Context()))
	must(0, str.SendRequestHeader(&http.Request{
		Method: "CONNECT",
		URL:    &url.URL{Opaque: connectto},
		Host:   connectto,
		Header: make(http.Header),
	}))

	rsp := must(str.ReadResponse())
	if rsp.StatusCode != 200 {
		must(0, fmt.Errorf("bad status code! %d", rsp.StatusCode))
	}

	//TODO: errors and stuff
	go io.Copy(os.Stdout, str)
	io.Copy(str, os.Stdin)
}

func must[T any](a T, err error) T {
	if err != nil {
		panic(err)
	}
	return a
}

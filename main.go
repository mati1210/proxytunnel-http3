package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

func multStrFlag(p *string, name []string, value, usage string) {
	for _, v := range name {
		flag.StringVar(p, v, value, usage)
	}
}

func main() {
	var proxy, to string
	multStrFlag(&proxy, []string{"proxy", "p"}, "", "Proxy to connect to (host:port)")
	multStrFlag(&to, []string{"dest", "d"}, "", "Destination to connect to (host:port)")
	flag.Parse()

	conn := (&http3.Transport{}).NewClientConn(must(quic.DialAddr(
		context.Background(),
		proxy,
		&tls.Config{
			// TODO: ech?
			NextProtos: []string{http3.NextProtoH3},
		},
		&quic.Config{
			// MaxIdleTimeout is by default 30s, so half that
			KeepAlivePeriod: 15 * time.Second,
		},
	)))

	str := must(conn.OpenRequestStream(conn.Context()))
	must(0, str.SendRequestHeader(&http.Request{
		Method: "CONNECT",
		URL:    &url.URL{Opaque: to},
		Host:   to,
	}))

	rsp := must(str.ReadResponse())
	if rsp.StatusCode < 200 && rsp.StatusCode >= 300 {
		panic(fmt.Errorf("bad status code! %d", rsp.StatusCode))
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

package proxy

import (
	"cargo-m/internal/until"
	"fmt"

	"github.com/txthinking/socks5"
)

type SocksProxy struct {
	ProxyServer *socks5.Server
}

func NewSocksProxy() *SocksProxy {
	return &SocksProxy{}
}

func (proxy *SocksProxy) Run(port int, authUser string, authPwd string) {
	addr := fmt.Sprintf(":%d", port)
	server, err := socks5.NewClassicServer(addr, "", authUser, authPwd, 0, 0)
	if err != nil {
		panic(err)
	}
	proxy.ProxyServer = server
	go func() {
		runnerError := server.ListenAndServe(nil)
		if runnerError != nil {
			panic(err)
		}
	}()
	until.Log.Info("socks5 proxy server started: " + addr)
}

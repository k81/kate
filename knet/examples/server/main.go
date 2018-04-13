package main

import (
	"context"

	"net/http"
	_ "net/http/pprof"

	"github.com/k81/kate/knet"
	"github.com/k81/kate/knet/examples/echo"
	"github.com/k81/kate/log"
)

var (
	mctx = log.SetContext(context.Background(), "moudle", "echosrv")
)

type srvEchoHandler struct {
	knet.IoHandlerAdapter
}

func (h *srvEchoHandler) OnConnected(session *knet.IoSession) error {
	log.Info(session.Context(), "session connected", "remote_addr", session.RemoteAddr())
	//m := echo.NewEchoMessage("welcome to the echo server, enjoy it!")
	//session.Send(m)
	return nil
}

func (h *srvEchoHandler) OnDisconnected(session *knet.IoSession) {
	log.Info(session.Context(), "session disconnected", "remote_addr", session.RemoteAddr(), "stats", session.String())
}

func (h *srvEchoHandler) OnError(session *knet.IoSession, err error) {
	log.Debug(session.Context(), "sesson error", "error", err)
}

func (h *srvEchoHandler) OnMessage(session *knet.IoSession, m knet.Message) error {
	echoMsg := m.(*echo.EchoMessage)
	log.Debug(session.Context(), "RECV", "msg", echoMsg.Content)
	session.Send(mctx, m)
	return nil
}

func main() {
	srvConf := knet.NewTCPServerConfig()
	//srvConf.MaxConnection = 2

	log.SetLevelByName("INFO")

	srv := knet.NewTCPServer(mctx, srvConf)
	srv.SetProtocol(&echo.EchoProtocol{})
	srv.SetIoHandler(&srvEchoHandler{})

	go func() {
		http.ListenAndServe("127.0.0.1:8890", nil)
	}()

	addr := "127.0.0.1:8888"
	log.Info(mctx, "server started", "addr", addr)
	srv.ListenAndServe(addr)
}

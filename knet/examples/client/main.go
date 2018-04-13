package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/k81/kate/knet"
	"github.com/k81/kate/knet/examples/echo"

	"net/http"
	_ "net/http/pprof"
)

type clientEchoHandler struct {
	knet.IoHandlerAdapter
}

func (h *clientEchoHandler) OnConnected(session *knet.IoSession) error {
	fmt.Println("----server connected----")
	return nil
}

func (h *clientEchoHandler) OnDisconnected(session *knet.IoSession) {
	fmt.Println("----server connection lost----")
}

func (h *clientEchoHandler) OnMessage(session *knet.IoSession, msg knet.Message) error {
	echoMsg := msg.(*echo.EchoMessage)
	fmt.Printf("\n----server message: %v----\n\n", echoMsg.Content)
	return nil
}

func main() {
	var (
		clientConf = knet.NewTCPClientConfig()
		client     knet.Client
		console    *bufio.Scanner
		err        error
	)
	clientConf.AutoReconnect = true
	client = knet.NewTCPClient(context.Background(), clientConf)

	client.SetProtocol(&echo.EchoProtocol{})
	client.SetIoHandler(&clientEchoHandler{})

	if err = client.Dial("127.0.0.1:8888"); err != nil {
		fmt.Printf("dial to server failed: error=%v\n", err)
		return
	}
	defer client.Close()

	time.Sleep(time.Second)

	go func() {
		http.ListenAndServe("127.0.0.1:8889", nil)
	}()

	console = bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("send>")
		if !console.Scan() {
			return
		}

		var (
			msg   = echo.NewEchoMessage(console.Text())
			reply knet.Response
		)

		if reply, err = client.Call(context.Background(), msg); err != nil {
			fmt.Printf("send error: error=%v\n", err)
			return
		}

		msg = reply.(*echo.EchoMessage)
		fmt.Printf("recv>%s\n", msg.Content)
	}
}

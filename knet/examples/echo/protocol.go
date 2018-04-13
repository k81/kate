package echo

import (
	"bufio"
	"bytes"
	"io"

	"github.com/k81/kate/knet"
)

type EchoMessage struct {
	Content string
}

func NewEchoMessage(data string) *EchoMessage {
	return &EchoMessage{
		Content: data,
	}
}

func (m *EchoMessage) Id() uint64 {
	return 0
}

const (
	KeyDecoder      = "decoder"
	KeyDecodeBuffer = "decode_buf"
)

type EchoProtocol struct{}

func (p *EchoProtocol) Decode(session *knet.IoSession, reader io.Reader) (m knet.Message, err error) {
	var (
		r        *bufio.Reader
		b        *bytes.Buffer
		ok       bool
		line     []byte
		isPrefix bool
	)

	if r, ok = session.GetAttr(KeyDecoder).(*bufio.Reader); !ok {
		r = bufio.NewReader(reader)
		session.SetAttr(KeyDecoder, r)
	}

	if b, ok = session.GetAttr(KeyDecodeBuffer).(*bytes.Buffer); !ok {
		b = &bytes.Buffer{}
		session.SetAttr(KeyDecodeBuffer, b)
	}

	if line, isPrefix, err = r.ReadLine(); err != nil {
		return
	}

	b.Write(line)

	if isPrefix {
		return
	}

	m = NewEchoMessage(b.String())
	b.Reset()
	return
}

func (p *EchoProtocol) Encode(session *knet.IoSession, m knet.Message) (data []byte, err error) {
	echoMsg := m.(*EchoMessage)
	return []byte(echoMsg.Content + "\n"), nil
}

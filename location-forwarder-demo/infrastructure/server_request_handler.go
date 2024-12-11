package infrastructure

import (
	"encoding/binary"
	"log"
	"net"
	"strconv"
)

type ServerRequestHandler struct {
	Host       string
	Port       int
	Connection net.Conn
	Ln         net.Listener
}

var err error

func NewServerRequestHandler(h string, p int) *ServerRequestHandler {
	r := new(ServerRequestHandler)

	r.Host = h
	r.Port = p
	r.Connection = nil

	r.Ln, err = net.Listen("tcp", h+":"+strconv.Itoa(p))
	if err != nil {
		log.Fatalf("ServerRequestHandler 0:: %s", err)
	}

	return r
}

func (srh *ServerRequestHandler) Receive() []byte {

	srh.Connection, err = srh.Ln.Accept()
	if err != nil {
		log.Fatalf("ServerRequestHandler:: Error receiving message's size %s", err)
	}

	// receive message's size
	size := make([]byte, 4)
	_, err = srh.Connection.Read(size)
	if err != nil {
		if _, ok := err.(*net.OpError); ok {
			srh.Connection.Close()
			return nil
		} else {
			log.Fatalf("ServerRequestHandler:: Error receiving message's size: %s", err)
		}
	}
	sizeInt := binary.LittleEndian.Uint32(size)

	// receive message
	msg := make([]byte, sizeInt)
	_, err = srh.Connection.Read(msg)
	if err != nil {
		if _, ok := err.(*net.OpError); ok {
			srh.Connection.Close()
			return nil
		} else {
			log.Fatalf("ServerRequestHandler:: Error receiving message: %s", err)
		}
	}
	return msg
}

func (srh *ServerRequestHandler) Send(msgToClient []byte) {

	// send message's size
	size := make([]byte, 4)
	l := uint32(len(msgToClient))
	binary.LittleEndian.PutUint32(size, l)
	_, err = srh.Connection.Write(size)
	if err != nil {
		if _, ok := err.(*net.OpError); ok {
			srh.Connection.Close()
			return
		} else {
			log.Fatalf("ServerRequestHandler 4:: %s", err)
		}
	}

	// end message
	_, err = srh.Connection.Write(msgToClient)
	if err != nil {
		if _, ok := err.(*net.OpError); ok {
			srh.Connection.Close()
			srh.Ln.Close()
			return
		} else {
			log.Fatalf("ServerRequestHandler 5:: %s", err)
		}
	}
}

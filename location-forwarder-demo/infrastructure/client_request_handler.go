package infrastructure

import (
	"encoding/binary"
	"log"
	"net"
	"strconv"
	"time"
)

type ClientRequestHandler struct {
	Host       string
	Port       int
	Connection net.Conn
}

func NewClientRequestHandler(h string, p int) *ClientRequestHandler {
	r := new(ClientRequestHandler)

	r.Host = h
	r.Port = p
	r.Connection = nil

	return r
}

func (crh *ClientRequestHandler) Handle(msgToServer []byte) []byte {
	var msgFromServer []byte

	for {
		crh.Connection, err = net.Dial("tcp", crh.Host+":"+strconv.Itoa(crh.Port))
		if err != nil {
			log.Println("Trying to connect to server...")
			time.Sleep(time.Second * 5)
			continue
		}

		// end message's size
		sizeMsgToServer := make([]byte, 4)
		l := uint32(len(msgToServer))
		binary.LittleEndian.PutUint32(sizeMsgToServer, l)
		_, err = crh.Connection.Write(sizeMsgToServer)
		if err != nil {
			log.Fatalf("ClientRequestHandler 1:: %s", err)
		}

		// send message
		_, err = crh.Connection.Write(msgToServer)
		if err != nil {
			log.Fatalf("ClientRequestHandler 2:: %s", err)
		}

		// receive message's size
		sizeMsgFromServer := make([]byte, 4)
		_, err = crh.Connection.Read(sizeMsgFromServer)
		if err != nil {
			log.Fatalf("ClientRequestHandler 3:: %s", err)
		}
		sizeFromServerInt := binary.LittleEndian.Uint32(sizeMsgFromServer)

		// receive reply
		msgFromServer = make([]byte, sizeFromServerInt)
		_, err = crh.Connection.Read(msgFromServer)
		if err != nil {
			log.Fatalf("ClientRequestHandler 4:: %s", err)
		}

		return msgFromServer
	}
	return msgFromServer
}

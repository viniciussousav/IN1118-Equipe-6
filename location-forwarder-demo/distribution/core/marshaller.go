package core

import (
	"encoding/gob"
	"encoding/json"
	"log"
)

type Marshaller struct{}

func (Marshaller) Marshall(msg Packet) []byte {
	r, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("Marshaller:: Marshall:: %s", err)
	}

	return r
}

func (Marshaller) Unmarshall(msg []byte) Packet {
	r := Packet{}
	err := json.Unmarshal(msg, &r)
	if err != nil {
		log.Fatalf("Marshaller:: Unmarshall:: %s", err)
	}
	return r
}

func (Marshaller) MarshallerFactory() Marshaller {
	gob.Register(Packet{})
	return Marshaller{}
}

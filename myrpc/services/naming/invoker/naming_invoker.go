package naminginvoker

import (
	"fmt"
	"log"
	"test/myrpc/distribution/core"
	"test/myrpc/infrastructure"
	"test/myrpc/services/naming/service"
	"test/shared"
)

type Invoker struct {
	Ior shared.IOR
}

func New(h string, p int) Invoker {
	i := shared.IOR{Host: h, Port: p}
	r := Invoker{Ior: i}
	return r
}

func (i Invoker) Invoke() {
	s := infrastructure.NewServerRequestHandler(i.Ior.Host, i.Ior.Port)
	m := core.Marshaller{}
	packet := core.Packet{}
	var rep interface{}

	// Create an instance of Calculadora
	n := naming_service.NamingService{}

	for {
		// Invoke ServerRequestHandler
		b := s.Receive()

		// Unmarshall miop packet
		packet = m.Unmarshall(b)

		// Extract request from publisher
		r := core.ExtractRequest(packet)

		// Demultiplex request
		switch r.Op {
		case "Find":
			_p1 := r.Params[0].(string)
			rep = n.Find(_p1)
		case "Bind":
			_p1 := r.Params[0].(string)
			_p22 := r.Params[1].(map[string]interface{})
			_ior := shared.IOR{Host: _p22["Host"].(string), Port: int(_p22["Port"].(float64)), Id: int(_p22["Id"].(float64)), TypeName: _p22["TypeName"].(string)}
			_p2 := _ior
			rep = n.Bind(_p1, _p2)
		case "List":
			rep = n.List()
		default:
			log.Fatal("Invoker:: Operation '" + r.Op + "' is unknown:: ")
		}

		// Prepare reply
		var params []interface{}
		params = append(params, rep)

		// Create reply packet
		replyPacket := core.CreateReplyPacket(params)

		// Marshall packet
		b = m.Marshall(replyPacket)

		// Send marshalled packet
		s.Send(b)

		fmt.Println("Services: ", n.List())
	}
}

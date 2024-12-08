package calculadorainvoker

import (
	"log"
	"test/myrpc/app/businesses/calculadora"
	"test/myrpc/distribution/core"
	"test/myrpc/distribution/interceptors"
	"test/myrpc/infrastructure"
	"test/shared"
)

type Invoker struct {
	ior               shared.IOR
	locationForwarder *interceptors.LocationForwarder
}

func NewInvoker(host string, port int, locationForwarder *interceptors.LocationForwarder) Invoker {
	return Invoker{
		ior:               shared.IOR{Host: host, Port: port},
		locationForwarder: locationForwarder}
}

func (i Invoker) Invoke() {
	s := infrastructure.NewServerRequestHandler(i.ior.Host, i.ior.Port)
	m := core.Marshaller{}
	packet := core.Packet{}

	var rep int

	c := calculadora.Calculadora{}

	for {
		log.Print("Received!")

		// Invoke ServerRequestHandler
		b := s.Receive()

		/*
			isAvailable := rand.Intn(2) == 1

			if !isAvailable {
				log.Print("Object not available locally. Forwarding request...")

				// Redireciona a requisição para o servidor remoto
				response, err := i.locationForwarder.ForwardRequest("Calculadora", b)
				if err != nil {
					log.Printf("Failed to forward request: %v", err)
					continue
				}

				// Processa a resposta recebida do servidor remoto
				log.Printf("Response from remote server: %s", string(response))
				continue
			}
		*/

		// Unmarshall miop packet
		packet = m.Unmarshall(b)

		// Extract request from publisher
		r := core.ExtractRequest(packet)

		_p1 := int(r.Params[0].(float64))
		_p2 := int(r.Params[1].(float64))

		// Demultiplex request & invoke Location Forwarder
		switch r.Op {
		case "Som":
			rep = c.Som(_p1, _p2)
		case "Dif":
			rep = c.Dif(_p1, _p2)
		case "Mul":
			rep = c.Mul(_p1, _p2)
		case "Div":
			rep = c.Div(_p1, _p2)
		default:
			log.Fatal("Invoker:: Operation '" + r.Op + "' is unknown:: ")
		}

		// Prepare reply
		var params []interface{}
		params = append(params, rep)

		// Create miop reply packet
		miop := core.CreateReplyMIOP(params)

		// Marshall miop packet
		b = m.Marshall(miop)

		// Send marshalled packet
		s.Send(b)
	}
}

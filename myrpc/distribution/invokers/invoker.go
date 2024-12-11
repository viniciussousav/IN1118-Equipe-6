package invokers

import (
	"log"
	"test/myrpc/app/businesses"
	"test/myrpc/distribution/core"
	"test/myrpc/distribution/interceptors"
	"test/myrpc/infrastructure"
	"test/shared"
)

type Invoker struct {
	ior               shared.IOR
	locationForwarder *interceptors.LocationForwarder
	localObjects      map[string]interface{}
	redirect          bool
}

func NewInvoker(host string, port int, locationForwarder *interceptors.LocationForwarder, redirect bool) Invoker {

	localObject := map[string]interface{}{
		"Calculadora1": businesses.Calculadora{},
	}

	return Invoker{
		ior:               shared.IOR{Host: host, Port: port},
		locationForwarder: locationForwarder,
		localObjects:      localObject,
		redirect:          redirect}
}

func (i Invoker) Invoke() {
	s := infrastructure.NewServerRequestHandler(i.ior.Host, i.ior.Port)
	m := core.Marshaller{}
	packet := core.Packet{}

	var rep int

	for {
		// Invoke ServerRequestHandler
		b := s.Receive()

		// Unmarshall packet
		packet = m.Unmarshall(b)

		// Extract request from publisher
		r := core.ExtractRequest(packet)
		log.Print("Request received:", r)

		c, exists := i.localObjects[r.ObjKey]

		if !exists && i.redirect {

			red, err := i.locationForwarder.ForwardRequest(r.ObjKey, b)

			if err != nil {
				// Prepare reply
				var params []interface{}
				params = append(params, "object not found")

				// Create reply packet
				replyPacket := core.CreateReplyPacket(params, 404)

				// Marshall packet
				b = m.Marshall(replyPacket)

				// Send marshalled packet
				s.Send(b)
				log.Print("Response sent:", replyPacket.Bd.RepBody)
				continue
			}

			// Send redirected packet
			s.Send(red)

			continue
		}

		if !exists && !i.redirect {
			newLocation, err := i.locationForwarder.GetLocation(r.ObjKey)
			if err != nil {
				// Prepare reply
				var params []interface{}
				params = append(params, "object not found")

				// Create reply packet
				replyPacket := core.CreateReplyPacket(params, 404)

				// Marshall packet
				b = m.Marshall(replyPacket)

				// Send marshalled packet
				s.Send(b)
				log.Print("Response sent:", replyPacket.Bd.RepBody)
				continue
			}

			var params []interface{}
			params = append(params, newLocation.Port)

			// Create reply packet
			replyPacket := core.CreateReplyPacket(params, 301)

			// Marshall packet
			b = m.Marshall(replyPacket)

			// Send marshalled packet
			s.Send(b)
			continue
		}

		if r.ObjKey == "Calculadora1" {

			calc := c.(businesses.Calculadora)

			_p1 := int(r.Params[0].(float64))
			_p2 := int(r.Params[1].(float64))

			// Demultiplex request & invoke Location Forwarder
			switch r.Op {
			case "Som":
				rep = calc.Som(_p1, _p2)
			case "Dif":
				rep = calc.Dif(_p1, _p2)
			default:
				log.Fatal("Invoker:: Operation '" + r.Op + "' is unknown:: ")
			}
		}

		// Prepare reply
		var params []interface{}
		params = append(params, rep)

		// Create reply packet
		replyPacket := core.CreateReplyPacket(params, 200)

		// Marshall packet
		b = m.Marshall(replyPacket)

		// Send marshalled packet
		s.Send(b)
		log.Print("Response sent:", replyPacket.Bd.RepBody)

	}
}

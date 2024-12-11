package invokers

import (
	"log"
	"test/location-forwarder-demo/app/businesses"
	"test/location-forwarder-demo/distribution/core"
	"test/location-forwarder-demo/distribution/interceptors"
	"test/location-forwarder-demo/infrastructure"
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

		b := s.Receive()
		packet = m.Unmarshall(b)
		r := core.ExtractRequest(packet)

		c, exists := i.localObjects[r.ObjKey]

		var params []interface{}

		if !exists && i.redirect {

			red, err := i.locationForwarder.ForwardRequest(r.ObjKey, b)

			if err != nil {
				params = append(params, "object not found")

				replyPacket := core.CreateReplyPacket(params, 404)
				b = m.Marshall(replyPacket)
				s.Send(b)

				log.Print("Response sent:", replyPacket.Bd.RepBody)
				continue
			}

			s.Send(red)
			continue
		}

		if !exists && !i.redirect {
			newLocation, err := i.locationForwarder.GetLocation(r.ObjKey)
			if err != nil {
				params = append(params, "object not found")

				replyPacket := core.CreateReplyPacket(params, 404)
				b = m.Marshall(replyPacket)
				s.Send(b)

				log.Print("Response sent:", replyPacket.Bd.RepBody)
				continue
			}

			params = append(params, newLocation.Port)

			replyPacket := core.CreateReplyPacket(params, 301)
			b = m.Marshall(replyPacket)
			s.Send(b)

			continue
		}

		calc := c.(businesses.Calculadora)

		_p1 := int(r.Params[0].(float64))
		_p2 := int(r.Params[1].(float64))

		switch r.Op {
		case "Som":
			rep = calc.Som(_p1, _p2)
		default:
			log.Println("Invoker:: Operation '" + r.Op + "' is unknown:: ")
			continue
		}

		params = append(params, rep)

		replyPacket := core.CreateReplyPacket(params, 200)
		b = m.Marshall(replyPacket)
		s.Send(b)

		log.Print("Response sent:", replyPacket.Bd.RepBody)
	}
}

func (i *Invoker) AddLocalObject(objKey string, objImpl interface{}) {
	i.localObjects[objKey] = objImpl
}

func (i *Invoker) RemoveLocalObject(objKey string) {
	delete(i.localObjects, objKey)
}

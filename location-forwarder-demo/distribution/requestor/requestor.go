package requestor

import (
	"test/location-forwarder-demo/distribution/core"
	"test/location-forwarder-demo/infrastructure"
	"test/shared"
)

type Requestor struct {
}

func (Requestor) Invoke(i shared.Invocation) shared.Termination {
	requestPacket := core.CreateRequestPackage(i.Request.ObjKey, i.Request.Op, i.Request.Params)

	m := core.Marshaller{}
	b := m.Marshall(requestPacket)

	c := infrastructure.NewClientRequestHandler(i.Ior.Host, i.Ior.Port)
	r := c.Handle(b)

	replyPacket := m.Unmarshall(r)
	rt := core.ExtractReply(replyPacket)

	t := shared.Termination{
		Status: replyPacket.Bd.RepHeader.Status,
		Rep:    rt,
	}

	return t
}

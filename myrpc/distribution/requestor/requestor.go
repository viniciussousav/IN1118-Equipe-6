package requestor

import (
	"test/myrpc/distribution/core"
	"test/myrpc/infrastructure"
	"test/shared"
)

type Requestor struct {
}

func NewRequestor() Requestor {
	return Requestor{}
}

func (Requestor) Invoke(i shared.Invocation) shared.Termination {
	// 1. Create packet
	requestPacket := core.CreateRequestPackage(i.Request.Op, i.Request.Params)

	// 2. Serialise MIOP packet
	m := core.Marshaller{}
	b := m.Marshall(requestPacket)

	// 3. Create & invoke ClientRequestHandler
	c := infrastructure.NewClientRequestHandler(i.Ior.Host, i.Ior.Port)
	r := c.Handle(b)

	// 4. Extract reply from subscriber
	replyPacket := m.Unmarshall(r)
	rt := core.ExtractReply(replyPacket)

	t := shared.Termination{Rep: rt}

	return t
}

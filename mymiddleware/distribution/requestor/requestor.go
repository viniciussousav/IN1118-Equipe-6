package requestor

import (
	"test/mymiddleware/distribution/marshaller"
	"test/mymiddleware/distribution/miop"
	"test/mymiddleware/infrastructure/crh"
	"test/shared"
)

type Requestor struct {
	Inv shared.Invocation
}

func (Requestor) Invoke(i shared.Invocation) shared.Termination {
	// 1. Create MIOP packet
	miopReqPacket := miop.CreateRequestMIOP(i.Request.Op, i.Request.Params)

	// 2. Serialise MIOP packet
	m := marshaller.Marshaller{}
	b := m.Marshall(miopReqPacket)

	// 3. Create & invoke CRH
	c := crh.NewCRH(i.Ior.Host, i.Ior.Port)
	r := c.SendReceive(b)

	// 4. Extract reply from server
	miopRepPacket := m.Unmarshall(r)
	rt := miop.ExtractReply(miopRepPacket)

	t := shared.Termination{Rep: rt}

	return t
}
package core

import (
	"test/shared"
)

type Packet struct {
	Hdr Header
	Bd  Body
}

type Header struct {
	Magic       string
	Version     string
	ByteOrder   bool
	MessageType int
	Size        int
}

type Body struct {
	ReqHeader RequestHeader
	ReqBody   RequestBody
	RepHeader ReplyHeader
	RepBody   ReplyBody
}

type RequestHeader struct {
	Context          string
	RequestId        int
	ResponseExpected bool
	ObjectKey        string
	Operation        string
}

type RequestBody struct {
	Body []interface{}
}

type ReplyHeader struct {
	Context   string
	RequestId int
	Status    int
}

type ReplyBody struct {
	OperationResult []interface{}
}

func CreateRequestPackage(objKey string, op string, p []interface{}) Packet {
	r := Packet{}

	header := Header{}
	body := Body{}
	reqHeader := RequestHeader{Operation: op, ObjectKey: objKey}
	reqBody := RequestBody{Body: p}
	body = Body{ReqHeader: reqHeader, ReqBody: reqBody}

	r.Hdr = header
	r.Bd = body

	return r
}

func CreateReplyPacket(params []interface{}, status int) Packet {
	r := Packet{}

	header := Header{}
	body := Body{}
	repHeader := ReplyHeader{"", 1313, status}
	repBody := ReplyBody{OperationResult: params}
	body = Body{RepHeader: repHeader, RepBody: repBody}

	r.Hdr = header
	r.Bd = body

	return r
}

func ExtractRequest(m Packet) shared.Request {
	i := shared.Request{}

	i.Op = m.Bd.ReqHeader.Operation
	i.Params = m.Bd.ReqBody.Body
	i.ObjKey = m.Bd.ReqHeader.ObjectKey

	return i
}

func ExtractReply(m Packet) shared.Reply {
	var r shared.Reply

	r.Result = m.Bd.RepBody.OperationResult

	return r
}

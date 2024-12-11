package shared

// Other configurations

const SampleSize = 1000
const LocalHost = "localhost"

type Invocation struct {
	Ior     IOR
	Request Request
}

type Termination struct {
	Status int
	Rep    Reply
}

type IOR struct {
	Host     string
	Port     int
	Id       int
	TypeName string
}

type Request struct {
	ObjKey string
	Op     string
	Params []interface{}
}

type Reply struct {
	Result []interface{}
}

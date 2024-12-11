package shared

const MaxConnectionAttempts = 100

// Other configurations
const StatisticSample = 30
const SampleSize = 1000
const CalculatorPort = 4040
const GrpcPort = 5050
const NAMING_PORT = 1414
const MIOP_REQUEST = 1
const MIOP_REPLY = 2
const MAX_NUMBER_CLIENTS = 1

const NamingPort = 1313
const CalculadoraPort = 1314
const LocalHost = "localhost"
const CallBackPort = 1317

type Message struct {
	Payload string
}

type Args struct {
	A, B int
}

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

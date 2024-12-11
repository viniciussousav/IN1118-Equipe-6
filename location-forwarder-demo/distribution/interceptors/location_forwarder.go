package interceptors

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"test/shared"
)

type LocationForwarder struct {
	RemoteLocations map[string]shared.IOR
}

// NewLocationForwarder Cria uma nova instância do Location Forwarder
func NewLocationForwarder() LocationForwarder {
	return LocationForwarder{
		RemoteLocations: map[string]shared.IOR{},
	}
}

// GetLocation retorna a localização do objeto remoto
func (lf *LocationForwarder) GetLocation(objectName string) (shared.IOR, error) {
	location, exists := lf.RemoteLocations[objectName]
	if !exists {
		return shared.IOR{}, fmt.Errorf("object %s not found in Location Forwarder", objectName)
	}
	log.Printf("Location of %s found: %s:%d", objectName, location.Host, location.Port)
	return location, nil
}

// ForwardRequest redireciona a requisição para o servidor remoto
func (lf *LocationForwarder) ForwardRequest(objectName string, request []byte) ([]byte, error) {
	location, err := lf.GetLocation(objectName)
	if err != nil {
		return nil, err
	}

	address := fmt.Sprintf("%s:%d", location.Host, location.Port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to remote server at %s: %v", address, err)
	}
	defer conn.Close()

	// send message's size
	sizeMsgToServer := make([]byte, 4)
	l := uint32(len(request))
	binary.LittleEndian.PutUint32(sizeMsgToServer, l)
	_, err = conn.Write(sizeMsgToServer)
	if err != nil {
		log.Fatalf("ClientRequestHandler 1:: %s", err)
	}

	// send message
	_, err = conn.Write(request)
	if err != nil {
		log.Fatalf("ClientRequestHandler 2:: %s", err)
	}

	// receive message's size
	sizeMsgFromServer := make([]byte, 4)
	_, err = conn.Read(sizeMsgFromServer)
	if err != nil {
		log.Fatalf("ClientRequestHandler 3:: %s", err)
	}
	sizeFromServerInt := binary.LittleEndian.Uint32(sizeMsgFromServer)

	// receive reply
	msgFromServer := make([]byte, sizeFromServerInt)
	_, err = conn.Read(msgFromServer)
	if err != nil {
		log.Fatalf("ClientRequestHandler 4:: %s", err)
	}

	return msgFromServer, nil
}

func (lf *LocationForwarder) AddLocation(objKey string, ior shared.IOR) {
	lf.RemoveLocation(objKey)
	lf.RemoteLocations[objKey] = ior
}

func (lf *LocationForwarder) RemoveLocation(objKey string) {
	delete(lf.RemoteLocations, objKey)
}

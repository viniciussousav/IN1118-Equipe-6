package interceptors

import (
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
		RemoteLocations: map[string]shared.IOR{
			"Calculator": {Host: "localhost", Port: 8082},
		},
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

	// Envia a requisição ao servidor remoto
	_, err = conn.Write(request)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}

	// Recebe a resposta do servidor remoto
	response := make([]byte, 1024)
	n, err := conn.Read(response)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	return response[:n], nil
}

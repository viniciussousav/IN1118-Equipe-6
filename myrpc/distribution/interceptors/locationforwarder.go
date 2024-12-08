package interceptors

import "test/shared"

type LocationForwarder struct {
	redirections map[string]shared.IOR
}

func NewLocationForwarder() *LocationForwarder {
	redirections := make(map[string]shared.IOR)
	return &LocationForwarder{redirections: redirections}
}

func (lf *LocationForwarder) SetLocation(operation string, location shared.IOR) {
	lf.redirections[operation] = location
}

func (lf *LocationForwarder) GetLocation(operation string) shared.IOR {
	return lf.redirections[operation]
}

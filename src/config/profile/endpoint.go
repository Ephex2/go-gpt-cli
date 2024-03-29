package profile
// Users can view create and modify profiles for different endpoints.
// This package aims to provide a central point for client code to refer to when trying to retrieve information about an endpoint, such as its default profile.

import (
	"errors"
)

// Endpoints are meant to represent different API endpoints that will be wrapped by client code.
// They are related to profiles in the sense that multiple profiles can be associated with an Endpoint
// The DefaultProfile() function should return the profile to use for commands related to a given endpoint
// This allows callers to know the endpoint and immediately have a profile ready to use to call functions related to that endpoint.
type Endpoint interface {
    Name() string
    DefaultProfile() Profile
    ProfileFromJsonBuf([]byte) (Profile, error)
}


// Collects endpoints added at runtime and allows them to be exposed to client code through endpoint.EndpointRegistry
type endpointRegistry struct {
    Endpoints []Endpoint
}

// The endpoint.EndpointRegistry is to be used by different endpoints that want to be callable by client code.
// The adding of profiles to the EndpointRegistry should be done in init() functions within the endpoints' respective packages.
// They should include an implementation of an Endpoint type and then call endpoint.EndpointRegistry.Add(theirImplementation).
var EndpointRegistry endpointRegistry


// Add Endpoint to list of Endpoints that will be recognized at runtime
func (r *endpointRegistry) Add(e Endpoint) {
    r.Endpoints = append(r.Endpoints, e)
}

func (r endpointRegistry) List() (s []string, err error) {
    if len(r.Endpoints) == 0 {
        err = errors.New("No endpoints to list") 
        return
    }

    for _, endpoint := range r.Endpoints {
       s = append(s, endpoint.Name()) 
    }

    return
}

func (r endpointRegistry) Get(name string) (e Endpoint, err error) {
    for _, endpoint := range r.Endpoints {
       if endpoint.Name() == name {
            e = endpoint
            break
       }
    }

    if e == nil {
        err = errors.New("No endpoint matches name: " + name)
    }

    return
}


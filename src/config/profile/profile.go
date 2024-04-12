package profile

type Repository interface {
	// Disk operations
	Create(endpoint Endpoint, profileName string) error
	Read(name string, endpointName string) ([]byte, error)
	Update(Profile) error
	Delete(endpointName string, profileName string) error
	GetAll(endpointName string) ([]string, error)
}

type Profile interface {
	Name() string
	Endpoint() Endpoint
    OverrideUrl() string
	ProfileRepository() Repository
	SetName(string) Profile
}

var RuntimeRepository Repository

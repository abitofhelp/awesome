package client_go

import awesomev1 "github.com/abitofhelp/awesome/gen/go/awesome/v1"

type (
	// Persistence is an interface for consumers of service's resources.
	// Usually, it simply wraps the grpc generated interface...
	Persistence interface {
		awesomev1.AwesomeServiceClient
	}
)

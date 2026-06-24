package apidef

import "github.com/gin-gonic/gin"

type Verb string

const (
	GET    Verb = "GET"
	POST   Verb = "POST"
	PUT    Verb = "PUT"
	PATCH  Verb = "PATCH"
	DELETE Verb = "DELETE"
)

type EndpointSpec struct {
	Verb        Verb
	Path        string
	RequireAuth bool
}

func (s EndpointSpec) RegisterOn(r gin.IRoutes, handler gin.HandlerFunc) {
	switch s.Verb {
	case GET:
		r.GET(s.Path, handler)
	case POST:
		r.POST(s.Path, handler)
	case PUT:
		r.PUT(s.Path, handler)
	case PATCH:
		r.PATCH(s.Path, handler)
	case DELETE:
		r.DELETE(s.Path, handler)
	}
}

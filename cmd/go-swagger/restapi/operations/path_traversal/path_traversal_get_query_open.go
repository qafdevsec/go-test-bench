// Code generated by go-swagger; DO NOT EDIT.

package path_traversal

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// PathTraversalGetQueryOpenHandlerFunc turns a function with the right signature into a path traversal get query open handler
type PathTraversalGetQueryOpenHandlerFunc func(PathTraversalGetQueryOpenParams) middleware.Responder

// Handle executing the request and returning a response
func (fn PathTraversalGetQueryOpenHandlerFunc) Handle(params PathTraversalGetQueryOpenParams) middleware.Responder {
	return fn(params)
}

// PathTraversalGetQueryOpenHandler interface for that can handle valid path traversal get query open params
type PathTraversalGetQueryOpenHandler interface {
	Handle(PathTraversalGetQueryOpenParams) middleware.Responder
}

// NewPathTraversalGetQueryOpen creates a new http.Handler for the path traversal get query open operation
func NewPathTraversalGetQueryOpen(ctx *middleware.Context, handler PathTraversalGetQueryOpenHandler) *PathTraversalGetQueryOpen {
	return &PathTraversalGetQueryOpen{Context: ctx, Handler: handler}
}

/* PathTraversalGetQueryOpen swagger:route GET /pathTraversal/os.Open/query/{safety} path_traversal pathTraversalGetQueryOpen

demonstrates Path Traversal via query, with vulnerable function os.Open

*/
type PathTraversalGetQueryOpen struct {
	Context *middleware.Context
	Handler PathTraversalGetQueryOpenHandler
}

func (o *PathTraversalGetQueryOpen) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPathTraversalGetQueryOpenParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

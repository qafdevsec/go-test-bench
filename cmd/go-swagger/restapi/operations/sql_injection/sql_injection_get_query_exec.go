// Code generated by go-swagger; DO NOT EDIT.

package sql_injection

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// SQLInjectionGetQueryExecHandlerFunc turns a function with the right signature into a SQL injection get query exec handler
type SQLInjectionGetQueryExecHandlerFunc func(SQLInjectionGetQueryExecParams) middleware.Responder

// Handle executing the request and returning a response
func (fn SQLInjectionGetQueryExecHandlerFunc) Handle(params SQLInjectionGetQueryExecParams) middleware.Responder {
	return fn(params)
}

// SQLInjectionGetQueryExecHandler interface for that can handle valid SQL injection get query exec params
type SQLInjectionGetQueryExecHandler interface {
	Handle(SQLInjectionGetQueryExecParams) middleware.Responder
}

// NewSQLInjectionGetQueryExec creates a new http.Handler for the SQL injection get query exec operation
func NewSQLInjectionGetQueryExec(ctx *middleware.Context, handler SQLInjectionGetQueryExecHandler) *SQLInjectionGetQueryExec {
	return &SQLInjectionGetQueryExec{Context: ctx, Handler: handler}
}

/* SQLInjectionGetQueryExec swagger:route GET /sqlInjection/sqlite3.exec/query/{safety} sql_injection sqlInjectionGetQueryExec

demonstrates SQL Injection via query, with vulnerable function sqlite3.exec

*/
type SQLInjectionGetQueryExec struct {
	Context *middleware.Context
	Handler SQLInjectionGetQueryExecHandler
}

func (o *SQLInjectionGetQueryExec) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewSQLInjectionGetQueryExecParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

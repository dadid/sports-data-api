package app

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

const (
	ctxKeyParams        CtxKey = "params"
	ctxKeyRequestID     CtxKey = "requestid"
	httpRequestIDHeader string = "X-Request-Id"
)

var (
	authSema = make(chan struct{}, 10) // counting semaphore using a buffered channel
)

// CtxKey represents the key for a value stored in a context
type CtxKey string

func (c CtxKey) String() string {
	return string(c)
}

func getParams(ctx context.Context) httprouter.Params {
	p := ctx.Value(ctxKeyParams)
	if params, ok := p.(httprouter.Params); ok {
		return params
	}
	return nil
}

func getRequestID(ctx context.Context) string {
	id := ctx.Value(ctxKeyRequestID)
	if requestID, ok := id.(string); ok {
		return requestID
	}
	return ""
}

func attachRequestID(ctx context.Context) context.Context {
	requestID := uuid.New()                                            // create new v4 UUID
	return context.WithValue(ctx, ctxKeyRequestID, requestID.String()) // return context with requestID attached
}

// addRequestID is a middleware that will attach a UUID requestID to the *http.Request and to an http header in the response
func (s *Server) addRequestID(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := attachRequestID(r.Context())                    // attach requestID to current request context
		next(w, r.WithContext(ctx))                            // call HandlerFunc with modified context
		w.Header().Set(httpRequestIDHeader, getRequestID(ctx)) // set requestID header
	})
}

// LimitNumClients is a middleware to ensure no more than maxClients requests are passed concurrently to the given handler f
func (s *Server) limitNumClients(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authSema <- struct{}{}
		defer func() { <-authSema }()
		next(w, r)
	})
}

// authenticate is a middleware that wraps a handlefunc and checks/validates Bearer Token cookie or header
func (s *Server) authenticate(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var bearerToken string
		bearerToken, err := checkBearerTokenCookie(r) // check for bearer token cookie
		switch err {
		case nil: // if err is nil do nothing
		case http.ErrNoCookie: // if no cookie is present check for the token header
			bearerToken, err = checkBearerTokenHeader(r)
			if err != nil {
				json.NewEncoder(w).Encode(Exception{Status: http.StatusUnauthorized, Message: err.Error()})
				return
			}
		default:
			json.NewEncoder(w).Encode(Exception{Status: http.StatusUnauthorized, Message: errors.Wrap(err, "no authorization token found").Error()})
			return
		}

		ok := authenticateToken(bearerToken) // if token is found then attempt to authenticate
		if ok {
			next(w, r)
			return
		}
		json.NewEncoder(w).Encode(Exception{Status: http.StatusUnauthorized, Message: "invalid authorization token"})
		return
	})
}

// httpRouterHandleWrapper is a generic wrapper that takes an http.HandlerFunc and returns and httprouter.Handle;
// Params are transferred to the request context
func (s *Server) httpRouterHandleWrapper(next http.HandlerFunc) httprouter.Handle {

	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := context.WithValue(r.Context(), ctxKeyParams, p)
		next(w, r.WithContext(ctx))
	})
}

// Middleware is a func type that receives and returns an http.HandlerFunc
type middleware func(http.HandlerFunc) http.HandlerFunc

// chainMiddleware allows you to chain multiple middlewares together for repeated usage
func (s *Server) chainMiddleware(h http.HandlerFunc, m ...middleware) http.HandlerFunc {
	if len(m) < 1 {
		return h
	}
	wrapped := h
	// loop in reverse to preserve middleware order
	for i := len(m) - 1; i >= 0; i-- {
		wrapped = m[i](wrapped)
	}
	return wrapped
}

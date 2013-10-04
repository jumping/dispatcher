// Package dispatcher provides a route dispatcher/HTTP request multiplexer
// for serving HTTP requests.
package dispatcher

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

// Regular expressions used for splitting paths and generating
// Route matchers for paths.
var (
	replaceCaptureParams = regexp.MustCompile(`\/\(`)
	replaceSlashes       = regexp.MustCompile(`([\/.])`)
	replaceWildcards     = regexp.MustCompile(`\*`)
	splitRoutePathParams = regexp.MustCompile(`(\/)?(\.)?:(\w+)(?:(\(.*?\)))?(\?)?`)
)

// Constants representing supported HTTP methods.
const (
	GET     = "GET"
	PUT     = "PUT"
	POST    = "POST"
	DELETE  = "DELETE"
	OPTIONS = "OPTIONS"
	HEAD    = "HEAD"
	TRACE   = "TRACE"
	CONNECT = "CONNECT"
	PATCH   = "PATCH"
)

// httpMethods is an array of strings containing the supported
// HTTP methods.
var (
	httpMethods = []string{GET, PUT, POST, DELETE, OPTIONS, HEAD, TRACE, CONNECT, PATCH}
)

// The Dispatcher type is an adapter to shorten creation
// of a map of strings to a map of pointers to Routes to
// HandlerFunc types from the http package.
type Dispatcher map[string]map[*Route]http.Handler

// The Middleware type is an adapter to allow the use of
// ordinary functions as middleware handlers.
type MiddlewareHandler func(res http.ResponseWriter, req *http.Request) bool

// ServeHTTP calls m(res, req)
func (m MiddlewareHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) bool {
	return m(res, req)
}

type Middleware interface {
	ServeHTTP(res http.ResponseWriter, req *http.Request) bool
}

type Router struct {
	mutex *sync.Mutex

	// Dispatcher map used for looking up the Router's Routes.
	dispatcher Dispatcher
	// Middleware each request served by the router should pass through.
	middleware []Middleware
	// handler used when Middleware and Routes fail to service the request.
	notFoundHandler http.Handler
	// strict flag to use when creating new Routes.
	strict bool
}

type Route struct {
	path    string         // path is the original path the Route was created for.
	keys    []string       // keys represents the names of the Route's parameters.
	matcher *regexp.Regexp // matcher is the regular expression used for matching the Route.
}

// fragmentedPathParameter is a struct that represents the strings
// generated by splitting a path with the `splitRoutePathParams`
// Regexp.
type fragmentedPathParameter struct {
	definition string
	slash      string
	format     string
	name       string
	capture    string
	optional   string
}

// RestrictRouteMatching sets a flag on the router causing
// routes ending with an unexpected trailing slash `/` to
// fail to match.
func (r *Router) RestrictRouteMatching() *Router {
	r.strict = true
	return r
}

// UnrestrictRouteMatching unsets a flag on the router
// allowing routes ending with an unexpected trailing slash `/` to
// match. By default, unrestricted routing is enabled.
func (r *Router) UnrestrictRouteMatching() *Router {
	r.strict = false
	return r
}

// Get registers a route to match the given path argument for
// HTTP GET requests. When a route is encounted that matches
// the path, the handler function argument is used to serve the
// requests.
func (r *Router) Get(path string, handler http.Handler) *Router {
	return r.AddHandler(GET, path, handler)
}

// Put registers a route to match the given path argument for
// HTTP PUT requests. When a route is encounted that matches
// the path, the handler function argument is used to serve the
// requests.
func (r *Router) Put(path string, handler http.Handler) *Router {
	return r.AddHandler(PUT, path, handler)
}

// Post registers a route to match the given path argument for
// HTTP POST requests. When a route is encounted that matches
// the path, the handler function argument is used to serve the
// requests.
func (r *Router) Post(path string, handler http.Handler) *Router {
	return r.AddHandler(POST, path, handler)
}

// Delete registers a route to match the given path argument for
// HTTP DELETE requests. When a route is encounted that matches
// the path, the handler function argument is used to serve the
// requests.
func (r *Router) Delete(path string, handler http.Handler) *Router {
	return r.AddHandler(DELETE, path, handler)
}

// Options registers a route to match the given path argument for
// HTTP OPTIONS requests. When a route is encounted that matches
// the path, the handler function argument is used to serve the
// requests.
func (r *Router) Options(path string, handler http.Handler) *Router {
	return r.AddHandler(OPTIONS, path, handler)
}

// Head registers a route to match the given path argument for
// HTTP HEAD requests. When a route is encounted that matches
// the path, the handler function argument is used to serve the
// requests.
func (r *Router) Head(path string, handler http.Handler) *Router {
	return r.AddHandler(HEAD, path, handler)
}

// Trace registers a route to match the given path argument for
// HTTP TRACE requests. When a route is encounted that matches
// the path, the handler function argument is used to serve the
// requests.
func (r *Router) Trace(path string, handler http.Handler) *Router {
	return r.AddHandler(TRACE, path, handler)
}

// Connect registers a route to match the given path argument for
// HTTP CONNECT requests. When a route is encounted that matches
// the path, the handler function argument is used to serve the
// requests.
func (r *Router) Connect(path string, handler http.Handler) *Router {
	return r.AddHandler(CONNECT, path, handler)
}

// Patch registers a route to match the given path argument for
// HTTP PATCH requests. When a route is encounted that matches
// the path, the handler function argument is used to serve the
// requests.
func (r *Router) Patch(path string, handler http.Handler) *Router {
	return r.AddHandler(PATCH, path, handler)
}

// Match registers a route to match the given path argument for
// any supported HTTP method. When a route is encounted that
// matches the path, the handler function argument is used to serve
// the requests.
func (r *Router) Match(path string, handler http.Handler) *Router {
	for _, method := range httpMethods {
		r.AddHandler(method, path, handler)
	}

	return r
}

// AddHandler creates a new Route matching path `path` under the
// Router's dispatchers `method`, setting its handler to `handler`.
// If the Router's dispatcher map does not previously have a key
// for `method`, the AddHandler assumes the `method` is unsupported
// and the Route created nor its handler will be added to the
// dispatcher.
func (r *Router) AddHandler(method, path string, handler http.Handler) *Router {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if routes, ok := r.dispatcher[method]; ok {
		route := NewRoute(path, r.strict)
		routes[route] = handler
	}

	return r
}

// RegisterMiddleware registers routing handlers that will be called
// with each HTTP request served.
func (r *Router) RegisterMiddleware(middleware Middleware) *Router {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.middleware = append(r.middleware, middleware)
	return r
}

// NotFound sets the routers handler that will be called when
// middleware does not handle the request's response and the
// path fails to match a known route.
func (r *Router) NotFound(handler http.Handler) *Router {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.notFoundHandler = handler
	return r
}

// findMatchingRouteAndHandler looks into the Router's dispatcher
// object in an attempt to find a matching route and handler function.
// If a pair are found, they are returned, else both will be nil.
func (r *Router) findMatchingRouteAndHandler(req *http.Request) (*Route, http.Handler) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	method := strings.ToUpper(req.Method)

	if routes, ok := r.dispatcher[method]; ok {
		for route, handler := range routes {
			if route.matcher.MatchString(req.URL.Path) {
				return route, handler
			}
		}
	}

	// Found no route or handler
	return nil, nil
}

// ServeHTTP handles all incoming HTTP requests. The request is first
// passed to each of the registered middleware functions. If the middleware
// function returns a boolean value of `true`, ServeHTTP returns early,
// assuming that the response has been served by it. If a middleware
// function fails to serve the request by returning `false`, ServeHTTP
// attempts to search for a Route that matches the requests URL. If a
// route is found, the request and response writer are handed over to
// the matched handler. If no middleware or route is found to handle
// the request, the Router's not found handler is used.
func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	for _, middleware := range r.middleware {
		if middleware.ServeHTTP(res, req) {
			// Midleware returned true meaning it handled the response, return
			// early.
			return
		}
	}

	route, handler := r.findMatchingRouteAndHandler(req)

	if nil == route || nil == handler {
		// No appropriate route and handler combination was found, allow
		// the notFoundHandler to serve the HTTP Request.
		r.notFoundHandler.ServeHTTP(res, req)
		return
	}

	// Middleware did not serve the request, pass it to the
	// handler.
	handler.ServeHTTP(res, req)
}

// NewDispatcher creates a new Dispatcher map, creating
// submaps for all supported HTTP methods.
func NewDispatcher() (dispatcher Dispatcher) {
	dispatcher = make(Dispatcher)

	for _, method := range httpMethods {
		dispatcher[method] = make(map[*Route]http.Handler)
	}

	return
}

// NewRouter creates a new Router object, returning a pointer
// to it. The Router's dispatcher is set with by calling the
// NewDispatcher method, and its not found handler is set to
// the http packages NotFoundHandler by default.
func NewRouter() (r *Router) {
	r = new(Router)
	r.dispatcher = NewDispatcher()
	r.notFoundHandler = http.NotFoundHandler()
	r.mutex = &sync.Mutex{}
	return
}

// NewRoute creates a new Route object, returning a pointer
// to it. A regular expression is generated to match the
// path provided. If strict is passed as true, routes ending
// with unexpected trailing slashes will fail to match
// the Route's regular expression.
func NewRoute(path string, strict bool) (route *Route) {
	route = new(Route)
	route.path = path

	compiled := replaceCaptureParams.ReplaceAllString(path, `(?:/`)
	parameters := splitRoutePathParams.FindAllStringSubmatch(path, -1)

	if !strict {
		compiled = fmt.Sprintf("%v/?", compiled)
	}

	for _, parameter := range parameters {
		fragmented := generateFragmentedPathParameter(parameter)

		var formatted string

		if 0 == len(fragmented.optional) {
			formatted = fmt.Sprintf("%v", fragmented.slash)
		}

		formatted = fmt.Sprintf("%v(?:", formatted)

		if 0 < len(fragmented.optional) {
			formatted = fmt.Sprintf("%v%v", formatted, fragmented.slash)
		}

		if 0 < len(fragmented.format) {
			formatted = fmt.Sprintf("%v%v", formatted, fragmented.format)
		}

		if 0 < len(fragmented.capture) {
			formatted = fmt.Sprintf("%v%v", formatted, fragmented.capture)
		} else if 0 < len(fragmented.format) {
			formatted = fmt.Sprintf("%v([^/.]+?)", formatted)
		} else {
			formatted = fmt.Sprintf("%v([^/]+?)", formatted)
		}

		formatted = fmt.Sprintf("%v)", formatted)

		if 0 < len(fragmented.optional) {
			formatted = fmt.Sprintf("%v%v", formatted, fragmented.optional)
		}

		compiled = strings.Replace(compiled, fragmented.definition, formatted, -1)
		route.keys = append(route.keys, fragmented.name)
	}

	compiled = replaceSlashes.ReplaceAllString(compiled, "\\$1")
	compiled = replaceWildcards.ReplaceAllString(compiled, "(.*)")
	route.matcher = regexp.MustCompile(fmt.Sprintf(`^%v$`, compiled))

	return
}

// generateFragmentedPathParameter returns a fragmentedPathParameter based
// on the parameter array provided.
func generateFragmentedPathParameter(parameter []string) (fragment fragmentedPathParameter) {
	fragment.definition = parameter[0]
	fragment.slash = parameter[1]
	fragment.format = parameter[2]
	fragment.name = parameter[3]
	fragment.capture = parameter[4]
	fragment.optional = parameter[5]
	return
}

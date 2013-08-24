// Package dispatcher provides a route dispatch for
// serving HTTP requests.
package dispatcher

import (
  "fmt"
  "net/http"
  "regexp"
  "strings"
)

var (
  replaceCaptureParams = regexp.MustCompile(`\/\(`)
  replaceSlashes       = regexp.MustCompile(`([\/.])`)
  replaceWildcards     = regexp.MustCompile(`\*`)
  splitRoutePathParams = regexp.MustCompile(`(\/)?(\.)?:(\w+)(?:(\(.*?\)))?(\?)?`)
)

const (
  GET    = "GET"
  PUT    = "PUT"
  POST   = "POST"
  DELETE = "DELETE"
)

var (
  httpMethods = []string{GET, PUT, POST, DELETE}
)

type Dispatcher map[string]map[*Route]http.HandlerFunc

type Router struct {
  dispatcher      Dispatcher
  notFoundHandler http.Handler
  strict          bool
}

type Route struct {
  path    string
  method  string
  keys    []string
  matcher *regexp.Regexp
}

type fragmentedPathParameter struct {
  definition string
  slash      string
  format     string
  name       string
  capture    string
  optional   string
}

func (r *Router) Get(path string, handler http.HandlerFunc) *Router {
  route := NewRoute(path, r.strict)
  r.dispatcher[GET][route] = handler
  return r
}

func (r *Router) Put(path string, handler http.HandlerFunc) *Router {
  route := NewRoute(path, r.strict)
  r.dispatcher[PUT][route] = handler
  return r
}

func (r *Router) Post(path string, handler http.HandlerFunc) *Router {
  route := NewRoute(path, r.strict)
  r.dispatcher[POST][route] = handler
  return r
}

func (r *Router) Delete(path string, handler http.HandlerFunc) *Router {
  route := NewRoute(path, r.strict)
  r.dispatcher[DELETE][route] = handler
  return r
}

func (r *Router) NotFound(handler http.HandlerFunc) *Router {
  r.notFoundHandler = handler
  return r
}

func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  method := strings.ToUpper(req.Method)

  if routes, ok := r.dispatcher[method]; ok {
    for route, handler := range routes {
      if route.matcher.MatchString(req.URL.Path) {
        handler.ServeHTTP(res, req)
        return
      }
    }
  }

  r.notFoundHandler.ServeHTTP(res, req)
}

func (r *Router) RestrictRouteMatching() *Router {
  r.strict = true
  return r
}

func (r *Router) UnrestrictRouteMatching() *Router {
  r.strict = false
  return r
}

func NewDispatcher() (dispatcher Dispatcher) {
  dispatcher = make(Dispatcher)

  for _, method := range httpMethods {
    dispatcher[method] = make(map[*Route]http.HandlerFunc)
  }

  return
}

func NewRouter() (r *Router) {
  r = new(Router)
  r.dispatcher = NewDispatcher()
  r.notFoundHandler = http.NotFoundHandler()
  return
}

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

func generateFragmentedPathParameter(segments []string) (fragment fragmentedPathParameter) {
  fragment.definition = segments[0]
  fragment.slash = segments[1]
  fragment.format = segments[2]
  fragment.name = segments[3]
  fragment.capture = segments[4]
  fragment.optional = segments[5]
  return
}

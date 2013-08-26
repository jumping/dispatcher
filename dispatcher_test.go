package dispatcher

import (
  "net/http"
  "testing"
)

// TestUnrestrictedRoute ensures routes with trailing slash are
// matched when route is unrestricted.
func TestUnrestrictedRoute(t *testing.T) {
  route := NewRoute("/test/:required", false)
  path := "/test/one/"

  if !route.matcher.MatchString(path) {
    t.Error("Expected unrestricted route to match path with trailing slash.")
  }
}

// TestUnrestrictedRoute ensures routes with trailing slash fails
// to match when route is restricted.
func TestRestrictedRoute(t *testing.T) {
  route := NewRoute("/test/:required", true)
  path := "/test/one/"

  if route.matcher.MatchString(path) {
    t.Error("Expected restricted route to fail to match path with trailing slash.")
  }
}

// TestOptionalRouteParameters ensures routes are matched when optional
// route parameters are ommited.
func TestOptionalRouteParameters(t *testing.T) {
  route := NewRoute("/test/:required/:optional?", false)
  path := "/test/one/two"

  if !route.matcher.MatchString(path) {
    t.Error("Expected route to match path with required and optional parameters supplied.")
  } else if path = "/test/one"; !route.matcher.MatchString(path) {
    t.Error("Expected route to match path with optional parameter missing.")
  }
}

// TestRoutedGetRequest ensures only the function registered via
// the Routers Get method respond to HTTP GET requests.
func TestRoutedGetRequest(t *testing.T) {
  counter := 0
  path := "/path/:to/:use"

  NewRouter().
    Get(path, generateCountableHandler(&counter)).
    Put(path, generateCountableHandler(&counter)).
    Post(path, generateCountableHandler(&counter)).
    Delete(path, generateCountableHandler(&counter)).
    Options(path, generateCountableHandler(&counter)).
    Head(path, generateCountableHandler(&counter)).
    Trace(path, generateCountableHandler(&counter)).
    Connect(path, generateCountableHandler(&counter)).
    Patch(path, generateCountableHandler(&counter)).
    Options(path, generateCountableHandler(&counter)).
    Head(path, generateCountableHandler(&counter)).
    Trace(path, generateCountableHandler(&counter)).
    Connect(path, generateCountableHandler(&counter)).
    Patch(path, generateCountableHandler(&counter)).
    ServeHTTP(nil, generateHttpRequest(GET, "/path/to/use"))

  if 1 != counter {
    t.Errorf("Expected counter to be set to 1, was set to %d.", counter)
  }
}

// TestRoutedPutRequest ensures only the function registered via
// the Routers Put method respond to HTTP PUT requests.
func TestRoutedPutRequest(t *testing.T) {
  counter := 0
  path := "/path/:to/:use"

  NewRouter().
    Get(path, generateCountableHandler(&counter)).
    Put(path, generateCountableHandler(&counter)).
    Post(path, generateCountableHandler(&counter)).
    Delete(path, generateCountableHandler(&counter)).
    Options(path, generateCountableHandler(&counter)).
    Head(path, generateCountableHandler(&counter)).
    Trace(path, generateCountableHandler(&counter)).
    Connect(path, generateCountableHandler(&counter)).
    Patch(path, generateCountableHandler(&counter)).
    ServeHTTP(nil, generateHttpRequest(PUT, "/path/to/use"))

  if 1 != counter {
    t.Errorf("Expected counter to be set to 1, was set to %d.", counter)
  }
}

// TestRoutedPostRequest ensures only the function registered via
// the Routers Post method respond to HTTP POST requests.
func TestRoutedPostRequest(t *testing.T) {
  counter := 0
  path := "/path/:to/:use"

  NewRouter().
    Get(path, generateCountableHandler(&counter)).
    Put(path, generateCountableHandler(&counter)).
    Post(path, generateCountableHandler(&counter)).
    Delete(path, generateCountableHandler(&counter)).
    Options(path, generateCountableHandler(&counter)).
    Head(path, generateCountableHandler(&counter)).
    Trace(path, generateCountableHandler(&counter)).
    Connect(path, generateCountableHandler(&counter)).
    Patch(path, generateCountableHandler(&counter)).
    ServeHTTP(nil, generateHttpRequest(POST, "/path/to/use"))

  if 1 != counter {
    t.Errorf("Expected counter to be set to 1, was set to %d.", counter)
  }
}

// TestRoutedDeleteRequest ensures only the function registered via
// the Routers Delete method respond to HTTP DELETE requests.
func TestRoutedDeleteRequest(t *testing.T) {
  counter := 0
  path := "/path/:to/:use"

  NewRouter().
    Get(path, generateCountableHandler(&counter)).
    Put(path, generateCountableHandler(&counter)).
    Post(path, generateCountableHandler(&counter)).
    Delete(path, generateCountableHandler(&counter)).
    Options(path, generateCountableHandler(&counter)).
    Head(path, generateCountableHandler(&counter)).
    Trace(path, generateCountableHandler(&counter)).
    Connect(path, generateCountableHandler(&counter)).
    Patch(path, generateCountableHandler(&counter)).
    ServeHTTP(nil, generateHttpRequest(DELETE, "/path/to/use"))

  if 1 != counter {
    t.Errorf("Expected counter to be set to 1, was set to %d.", counter)
  }
}

// TestRoutedOptionsRequest ensures only the function registered via
// the Routers Options method respond to HTTP OPTIONS requests.
func TestRoutedOptionsRequest(t *testing.T) {
  counter := 0
  path := "/path/:to/:use"

  NewRouter().
    Get(path, generateCountableHandler(&counter)).
    Put(path, generateCountableHandler(&counter)).
    Post(path, generateCountableHandler(&counter)).
    Delete(path, generateCountableHandler(&counter)).
    Options(path, generateCountableHandler(&counter)).
    Head(path, generateCountableHandler(&counter)).
    Trace(path, generateCountableHandler(&counter)).
    Connect(path, generateCountableHandler(&counter)).
    Patch(path, generateCountableHandler(&counter)).
    ServeHTTP(nil, generateHttpRequest(OPTIONS, "/path/to/use"))

  if 1 != counter {
    t.Errorf("Expected counter to be set to 1, was set to %d.", counter)
  }
}

// TestRoutedHeadRequest ensures only the function registered via
// the Routers Head method respond to HTTP HEAD requests.
func TestRoutedHeadRequest(t *testing.T) {
  counter := 0
  path := "/path/:to/:use"

  NewRouter().
    Get(path, generateCountableHandler(&counter)).
    Put(path, generateCountableHandler(&counter)).
    Post(path, generateCountableHandler(&counter)).
    Delete(path, generateCountableHandler(&counter)).
    Options(path, generateCountableHandler(&counter)).
    Head(path, generateCountableHandler(&counter)).
    Trace(path, generateCountableHandler(&counter)).
    Connect(path, generateCountableHandler(&counter)).
    Patch(path, generateCountableHandler(&counter)).
    ServeHTTP(nil, generateHttpRequest(HEAD, "/path/to/use"))

  if 1 != counter {
    t.Errorf("Expected counter to be set to 1, was set to %d.", counter)
  }
}

// TestRoutedTraceRequest ensures only the function registered via
// the Routers Trace method respond to HTTP TRACE requests.
func TestRoutedTraceRequest(t *testing.T) {
  counter := 0
  path := "/path/:to/:use"

  NewRouter().
    Get(path, generateCountableHandler(&counter)).
    Put(path, generateCountableHandler(&counter)).
    Post(path, generateCountableHandler(&counter)).
    Delete(path, generateCountableHandler(&counter)).
    Options(path, generateCountableHandler(&counter)).
    Head(path, generateCountableHandler(&counter)).
    Trace(path, generateCountableHandler(&counter)).
    Connect(path, generateCountableHandler(&counter)).
    Patch(path, generateCountableHandler(&counter)).
    ServeHTTP(nil, generateHttpRequest(TRACE, "/path/to/use"))

  if 1 != counter {
    t.Errorf("Expected counter to be set to 1, was set to %d.", counter)
  }
}

// TestRoutedConnectRequest ensures only the function registered via
// the Routers Connect method respond to HTTP CONNECT requests.
func TestRoutedConnectRequest(t *testing.T) {
  counter := 0
  path := "/path/:to/:use"

  NewRouter().
    Get(path, generateCountableHandler(&counter)).
    Put(path, generateCountableHandler(&counter)).
    Post(path, generateCountableHandler(&counter)).
    Delete(path, generateCountableHandler(&counter)).
    Options(path, generateCountableHandler(&counter)).
    Head(path, generateCountableHandler(&counter)).
    Trace(path, generateCountableHandler(&counter)).
    Connect(path, generateCountableHandler(&counter)).
    Patch(path, generateCountableHandler(&counter)).
    ServeHTTP(nil, generateHttpRequest(CONNECT, "/path/to/use"))

  if 1 != counter {
    t.Errorf("Expected counter to be set to 1, was set to %d.", counter)
  }
}

// TestRoutedPatchRequest ensures only the function registered via
// the Routers Patch method respond to HTTP PATCH requests.
func TestRoutedPatchRequest(t *testing.T) {
  counter := 0
  path := "/path/:to/:use"

  NewRouter().
    Get(path, generateCountableHandler(&counter)).
    Put(path, generateCountableHandler(&counter)).
    Post(path, generateCountableHandler(&counter)).
    Delete(path, generateCountableHandler(&counter)).
    Options(path, generateCountableHandler(&counter)).
    Head(path, generateCountableHandler(&counter)).
    Trace(path, generateCountableHandler(&counter)).
    Connect(path, generateCountableHandler(&counter)).
    Patch(path, generateCountableHandler(&counter)).
    ServeHTTP(nil, generateHttpRequest(PATCH, "/path/to/use"))

  if 1 != counter {
    t.Errorf("Expected counter to be set to 1, was set to %d.", counter)
  }
}

// TestRoutedMatchRequest ensures functions registered via
// the Routers Match method respond to any supported HTTP methods.
func TestRoutedMatchRequest(t *testing.T) {
  counter := 0
  path := "/path/:to/:use"

  router := NewRouter().
    Match(path, generateCountableHandler(&counter))

  for _, method := range httpMethods {
    router.ServeHTTP(nil, generateHttpRequest(method, "/path/to/use"))
  }

  if counter != len(httpMethods) {
    t.Error("Expected Router's Match method to match requests of any method")
  }
}

// TestUnhandledMiddleware ensures Middleware is executed on
// a routed request and eventually handled by a registered route
// handler.
func TestUnhandledMiddleware(t *testing.T) {
  counter := 0
  path := "/path/:to/:use"

  NewRouter().
    Get(path, generateCountableHandler(&counter)).
    RegisterMiddleware(generateCountableMiddleware(&counter, false)).
    ServeHTTP(nil, generateHttpRequest(GET, "/path/to/use"))

  if 2 != counter {
    t.Error("Expected the Middleware handler and Route handler to be called.")
  }
}

// TestHandledMiddleware ensures Middleware is executed on
// a routed request and all Middleware/Routes following it
// are not.
func TestHandledMiddleware(t *testing.T) {
  counter := 0
  path := "/path/:to/:use"

  NewRouter().
    Get(path, generateCountableHandler(&counter)).
    RegisterMiddleware(generateCountableMiddleware(&counter, true)).
    ServeHTTP(nil, generateHttpRequest(GET, "/path/to/use"))

  if 1 != counter {
    t.Error("Expected the Middleware handler to handle the request.")
  }
}

// TestNotFoundRoute ensures that any Route not matched is
// handled by the handler registered via the Router's
// NotFound method.
func TestNotFoundRoute(t *testing.T) {
  counter := 0
  path := "/path/:to/:use"

  NewRouter().
    Get(path, generateCountableHandler(&counter)).
    Put(path, generateCountableHandler(&counter)).
    Post(path, generateCountableHandler(&counter)).
    Delete(path, generateCountableHandler(&counter)).
    Options(path, generateCountableHandler(&counter)).
    Head(path, generateCountableHandler(&counter)).
    Trace(path, generateCountableHandler(&counter)).
    Connect(path, generateCountableHandler(&counter)).
    Patch(path, generateCountableHandler(&counter)).
    NotFound(generateCountableHandler(&counter)).
    ServeHTTP(nil, generateHttpRequest(GET, "/test/unkown/route"))

  if 1 != counter {
    t.Errorf("Expected counter to be set to 1, was set to %d.", counter)
  }
}

// generateHttpRequest is a helper function to generate a Request
// from the http package to use with testing.
func generateHttpRequest(method, path string) (req *http.Request) {
  req, err := http.NewRequest(method, path, nil)

  if nil != err {
    panic(err)
  }

  return
}

// generateCountableHandler is a helper function to generate a
// HandlerFunc from the http package to use with testing.
func generateCountableHandler(n *int) http.HandlerFunc {
  return func(res http.ResponseWriter, req *http.Request) {
    *n += 1
  }
}

// generateCountableMiddleware is a helper function to generate a
// Middleware from the dispatcher package to use with testing.
func generateCountableMiddleware(n *int, handle bool) Middleware {
  return func(res http.ResponseWriter, req *http.Request) bool {
    *n += 1
    return handle
  }
}

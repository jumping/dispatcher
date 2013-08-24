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
  ReplaceCaptureParams = regexp.MustCompile(`\/\(`)
  ReplaceSlashes       = regexp.MustCompile(`([\/.])`)
  ReplaceWildcard      = regexp.MustCompile(`\*`)
  SplitRoutePath       = regexp.MustCompile(`(\/)?(\.)?:(\w+)(?:(\(.*?\)))?(\?)?`)
)

type HttpMethod string

const (
  GET    HttpMethod = "GET"
  PUT    HttpMethod = "PUT"
  POST   HttpMethod = "POST"
  DELETE HttpMethod = "DELETE"
)

type Dispatcher map[HttpMethod]map[*Route]http.HandlerFunc

type Router struct {
  dispatcher Dispatcher
  strict     bool
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

func NewRoute(path string, strict bool) (route *Route) {
  route = new(Route)
  route.path = path

  compiled := ReplaceCaptureParams.ReplaceAllString(path, `(?:/`)
  parameters := SplitRoutePath.FindAllStringSubmatch(path, -1)

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

  compiled = ReplaceSlashes.ReplaceAllString(compiled, "\\$1")
  compiled = ReplaceWildcard.ReplaceAllString(compiled, "(.*)")
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

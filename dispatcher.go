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
  ReplaceOptionalParams = regexp.MustCompile(`\/\(`)
  ReplaceRequiredParams = regexp.MustCompile(`(\/)?(\.)?:(\w+)(?:(\(.*?\)))?(\?)?`)
  ReplaceSlashes        = regexp.MustCompile(`([\/.])`)
  ReplaceWildcard       = regexp.MustCompile(`\*`)
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
  Dispatcher    Dispatcher
  StrictRouting bool
}

type Route struct {
  Path   string
  Method string
  Keys   []string
  Regexp *regexp.Regexp
}

type FragmentedPathParameter struct {
  Definition string
  Slash      string
  Format     string
  Name       string
  Capture    string
  Optional   string
}

func NewRoute(path string, strict bool) (route *Route) {
  route = new(Route)
  route.Path = path

  compiled := ReplaceOptionalParams.ReplaceAllString(path, `(?:/`)
  parameters := ReplaceRequiredParams.FindAllStringSubmatch(path, -1)

  if !strict {
    compiled = fmt.Sprintf("%v/?", compiled)
  }

  for _, paramenter := range parameters {
    fragmented := GenerateFragmentedPathParameter(paramenter)

    var formatted string

    if 0 == len(fragmented.Optional) {
      formatted = fmt.Sprintf("%v", fragmented.Slash)
    }

    formatted = fmt.Sprintf("%v(?:", formatted)

    if 0 < len(fragmented.Optional) {
      formatted = fmt.Sprintf("%v%v", formatted, fragmented.Slash)
    }

    if 0 < len(fragmented.Format) {
      formatted = fmt.Sprintf("%v%v", formatted, fragmented.Format)
    }

    if 0 < len(fragmented.Capture) {
      formatted = fmt.Sprintf("%v%v", formatted, fragmented.Capture)
    } else if 0 < len(fragmented.Format) {
      formatted = fmt.Sprintf("%v([^/.]+?)", formatted)
    } else {
      formatted = fmt.Sprintf("%v([^/]+?)", formatted)
    }

    formatted = fmt.Sprintf("%v)", formatted)

    if 0 < len(fragmented.Optional) {
      formatted = fmt.Sprintf("%v%v", formatted, fragmented.Optional)
    }

    compiled = strings.Replace(compiled, fragmented.Definition, formatted, -1)
    route.Keys = append(route.Keys, fragmented.Name)
  }

  compiled = ReplaceSlashes.ReplaceAllString(compiled, "\\$1")
  compiled = ReplaceWildcard.ReplaceAllString(compiled, "(.*)")
  route.Regexp = regexp.MustCompile(fmt.Sprintf(`^%v$`, compiled))

  return
}

func GenerateFragmentedPathParameter(segments []string) (fragment FragmentedPathParameter) {
  fragment.Definition = segments[0]
  fragment.Slash = segments[1]
  fragment.Format = segments[2]
  fragment.Name = segments[3]
  fragment.Capture = segments[4]
  fragment.Optional = segments[5]
  return
}

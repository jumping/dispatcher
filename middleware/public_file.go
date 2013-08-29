// Package middleware provides common middleware functions
// for the dispatcher package.
package middleware

import (
  "fmt"
  "io/ioutil"
  "mime"
  "net/http"
  "os"
  "path"
)

import (
  "github.com/chuckpreslar/dispatcher"
)

// ServePublicFilesFrom accepts a `directory` argument where public
// files (i.e. javascript, css, and image files) can be found
// and returns a function to serve files stored in that `directory`.
// If a request is encounted matching a file found within the public
// folder, the `Content-Type` header is set on the response to
// a MIME type matching the file's extention and the files content
// is written and the function returns true to halt further attempts
// to serve the request. If no file is found, the function returns false
// to allow other middleware or a potential dispatcher Route handler
// to serve the request.
func ServePublicFilesFrom(directory string) dispatcher.Middleware {

  return func(res http.ResponseWriter, req *http.Request) bool {
    location := path.Join(directory, req.URL.Path)
    file, err := os.Open(location)

    if nil != err {
      return false
    } else if stat, err := file.Stat(); nil != err || stat.IsDir() {
      return false
    }

    data, err := ioutil.ReadFile(location)

    if nil != err {
      return false
    }

    // Determing the MIME type of the file located at `location`.
    typ := mime.TypeByExtension(path.Ext(location))

    // Write the Content-Type header of the public file.
    header := res.Header()
    header.Add("Content-Type", typ)

    if _, err := fmt.Fprintf(res, "%s", data); nil != err {
      return false
    }

    return true
  }
}

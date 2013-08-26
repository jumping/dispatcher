# dispatcher

Router dispatcher for Go web applications.

[![Build Status](https://drone.io/github.com/chuckpreslar/dispatcher/status.png)](https://drone.io/github.com/chuckpreslar/dispatcher/latest)

## Installation

With Google's [Go](http://www.golang.org) installed on your machine:

    $ go get -u github.com/chuckpreslar/dispatcher

## Usage

Getting started with Dispatcher is simple.  The following is an example of how to use the package.

```go
package main

// Standard library imports.
import(
  "net/http"
)

// Package imports.
import(
  "github.com/chuckpreslar/dispatcher"
)

func Index(w http.ResponseWriter, r *http.Request) {
  // Index route handler.
}

func ViewPosts(w http.ResponseWriter, r *http.Request) {
  // ViewPosts route handler.  
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
  // CreatePosts route handler.  
}

func

func main() {
  // Create an new Router instance from the dispatcher package.
  router := dispatcher.NewRouter()
  
  router.
    // Your servers index route.
    Get("/", Index).
    // The `:id` route parameter is made optional, will respond to `/posts` and `/posts/1`.
    Get("/posts/:id?", ViewPosts). 
    // Responds to HTTP PUT requests that match the `/posts` path.
    Put("/post", CreatePost)
  
  // Start your server.
  http.ListenAndServe(":3000", router)
}

```

### Middleware

It is possible to provide middleware for your router to use with each HTTP request.  It is the responsibility of the middleware to let the Router know if it has handled the request or not by returning `true` or `false` respectively.  Middleware handlers are called regardless of the requests HTTP method.

``` go
func SampleMiddlewareHandler(w http.ResponseWriter, r *http.Request) bool {
  // Do something.
  return false
}

func main() {
  // ...
  router.RegisterMiddleware(SampleMiddlewareHandler)
  http.ListenAndServe(":3000", router)
}
```


## Documentation

View godoc or visit [godoc.org](http://godoc.org/github.com/chuckpreslar/dispatcher).

    $ godoc dispatcher

## License

> The MIT License (MIT)

> Copyright (c) 2013 Chuck Preslar

> Permission is hereby granted, free of charge, to any person obtaining a copy
> of this software and associated documentation files (the "Software"), to deal
> in the Software without restriction, including without limitation the rights
> to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
> copies of the Software, and to permit persons to whom the Software is
> furnished to do so, subject to the following conditions:

> The above copyright notice and this permission notice shall be included in
> all copies or substantial portions of the Software.

> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
> IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
> FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
> AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
> LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
> OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
> THE SOFTWARE.

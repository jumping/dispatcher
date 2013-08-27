# dispatcher

Router dispatcher for Go web applications.

[![Build Status](https://drone.io/github.com/chuckpreslar/dispatcher/status.png)](https://drone.io/github.com/chuckpreslar/dispatcher/latest)

## Installation

With Google's [Go](http://www.golang.org) installed on your machine:

    $ go get -u github.com/chuckpreslar/dispatcher

## Usage

### Introduction

To get started using dispatcher, simply:

```go
package main

import(
    "net/http"
)

import(
    "github.com/chuckpreslar/dispatcher"
)

func SomePathHandler(w http.ResponseWriter, r *http.Request) {
    // Handle the request.
}

func main() {
    // Initialize a new Router instance.
    router := dispatcher.NewRouter()
    
    // Register a HandlerFunc.
    router.Match("/some/path", SomePathHandler)
    
    // Start your server.
    http.ListenAndServe(":3000" /* Port to listen on */, router)
}
```

The Router's Match method will match a request based only on the HTTP Request's URL Path, ignoring the method.

### Routing based on HTTP Method

Dispatcher supports route dispatching for most* supported HTTP methods.  To have a handler respond to a specific HTTP method, simply call the corresponding receiver method on the Router you've created passing it the path to match and the handler to use when a request with a matching path is encountered.
```go
//...

func main() {
    router := dispatcher.NewRouter()
    
    router.
        Get("/get/something", GetSomethingHandler).
        Put("/put/something", PutSomethingHandler).
        Post("/post/something", PostSomethingHandler).
        Delete("/delete/something", DeleteSomethingHandler)
        
    http.ListenAndServe(":3000" /* Port to listen on */, router)
}
```

### Path Matching

__Match Explicit Path__

```go
    // Matches route `/posts`
    router.Match("/posts", PostsHandler)
```

__Match Path with Required Parameters__

```go
    // Matches route `/posts/2012`, `/posts/2013`, etc.
    router.Match("/posts/:year", IndividualPostsHandler)
```

__Match Path with Optional Parameters__

```go
    // Matches route `/posts/2013` and `/posts/2013/january`
    router.Match("/posts/:id/:option?", IndividualPostsHandler)
```

__Match Wildcard__

```go    
    // Matches any route starting with `/posts` (i.e. `/posts/`, `/posts/2013`, `/posts/2013/january`)
    router.Match("/posts/*", WildcardPostsHandler)
```

### Accessing Path Parameters
    
### Middleware

Route middleware is registered as follows:

```go
    //...
    router.RegisterMiddleware(func(res http.ResponseWriter, req *http.Request) bool {
        // Middleware handler.
        return true || false
    })
```

Dispatcher attempts to call each piece of registered middleware with every request.  If the middleware handler returns true, Dispatcher assumes that the request was handled by the middleware and it no longer needs to attempt to find a registered Route and handler for the request.  If the middleware returns false, the next registered middleware handler runs or an attempt to find a registered Route and handler is made.

__TODO:__
* Finalize route parameter retrieval.
* Finalize public asset serving middleware.
* Finalize session support middleware.

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

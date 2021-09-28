// Package adapter is a simple mechanism for chaining "middleware" handlers.
// A common pattern for injecting data into a handler involves returning an
// http.Handler from a func:
//
// func InjectableHandler(data interface{}) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// do something with data
// 		w.WriteHeader(http.StatusOK)
// 	})
// }
//
// However, chaining many of these together can be difficult as passing handlers
// to each other can lead to callback hell. This package proposes a different
// pattern that, while slightly strange at first, adds readability:
//
// func InjectableAdapter(data interface{}) adapter.Adapter {
// 	return adapter.Adapter(func(h http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			// do something with data
// 			h.ServeHTTP(w, r)
// 		})
// 	})
// }
//
// At first blush returning a func that returns a func feels awkward, but allows
// for readable invocation where the intended "last handler" is called out:
//
// func HandleRequest(data interface{}) http.Handler {
// 	return adapter.Adapt(
// 		FinalHandler(data), // returns an http.Handler
// 		FirstHandler(data), // returns an Adapter
// 		SecondHandler(data), // returns an Adapter
// 	)
// }
//
// All of these funcs allow for injecting data at each step. This allows for the
// one-time work on that data before being bound to the returned func. Allowing
// for runtime/startup work to be done only once, while binding the results to
// the request. Further, if there are errors each handler can simply not call
// the next handler.
//
// Ripped straight from https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81
//
// Feel free to import this, but it might be much simpler to copy-n-paste as "A
// little copying is better than a little dependency."
package main

import "net/http"

// Adapter is a type that takes a http.Handler and return a http.Handler. This
// allows http.Handlers to be chained. Each http.Handler must invoke the
// previous.
type Adapter func(http.Handler) http.Handler

// Adapt applies all the of the given adapters (left-to-right, source code
// order) to the given http.Handler. Each Adapter must invoke the previous.
func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	// call the adapters in source code order; code should be readable
	for i := len(adapters) - 1; i >= 0; i-- {
		h = adapters[i](h)
	}

	return h
}

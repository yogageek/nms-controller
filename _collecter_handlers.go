package middleware

import (
	"fmt"
	"net/http"
	"time"
)

// func Mware(next http.Handler) http.Handler {
//   return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//     //  if r.RequestURI is in mwarePaths {
//     //     // do the middleware
//     //  }

//   }
// }

func timeHandler(format string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		tm := time.Now().Format(format)
		w.Write([]byte("The time is: " + tm))
	}
	return http.HandlerFunc(fn)
}

func timeHandler2(format string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tm := time.Now().Format(format)
		w.Write([]byte("The time is: " + tm))
	}
}

func checkSecurityA(next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		header := req.Header.Get("Super-Duper-Safe-Security")
		if header != "password" {
			fmt.Fprint(res, "Invalid password")
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		next(res, req)
	}
}

type httpHandlerFunc func(http.ResponseWriter, *http.Request)

func checkSecurity(next httpHandlerFunc) httpHandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// header := req.Header.Get("token")
		// if header != config.Token {
		// 	res.WriteHeader(http.StatusUnauthorized)
		// 	res.Write([]byte("401 - Unauthorized"))
		// 	return
		// }
		next(res, req)
	}
}

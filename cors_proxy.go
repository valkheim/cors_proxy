package main

import (
	"crypto/tls"
	"io"
	"net/http"
)

func main() {
	(&http.Server{
		Addr: ":8000",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.URL.Scheme = "https"
			r.URL.Host = "api.local"
			r.Header.Set("Origin", "app.local")

			http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
			res, err := http.DefaultTransport.RoundTrip(r)
			if err != nil {
				return
			}

			for k, vv := range res.Header {
				for _, v := range vv {
					w.Header().Add(k, v)
				}
			}

			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(res.StatusCode)
			io.Copy(w, res.Body)
			res.Body.Close()
		}),
	}).ListenAndServe()
}

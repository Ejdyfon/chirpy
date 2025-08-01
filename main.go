package main

import "net/http"

func customHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Add("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(200)
	rw.Write([]byte("OK"))

}

func main() {
	srv := http.NewServeMux()
	srv.HandleFunc("/healthz", customHandler)
	srv.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	httpsrv := http.Server{}
	httpsrv.Handler = srv
	httpsrv.Addr = ":8080"
	httpsrv.ListenAndServe()

}

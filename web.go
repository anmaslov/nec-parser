package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type CallFilter struct {
	limit int
	skip int
	tp string
	stantion string
	phone string
	called string
	start string
	end string
}

func filterCalls(w http.ResponseWriter, r *http.Request) {

	callF := CallFilter{}
	var err error
	if len(r.URL.RawQuery) > 0 {
		limitStr := r.URL.Query().Get("limit")
		callF.limit, err = strconv.Atoi(limitStr)
		if err != nil {
			callF.limit = 100
		}

		offsetStr := r.URL.Query().Get("skip")
		callF.skip, err = strconv.Atoi(offsetStr)
		if err != nil {
			callF.skip = 0
		}

		callF.stantion = r.URL.Query().Get("stantion")
		callF.phone = r.URL.Query().Get("phone")

		callF.start = r.URL.Query().Get("start")
		callF.end = r.URL.Query().Get("end")
		callF.called = r.URL.Query().Get("called")
		callF.tp = r.URL.Query().Get("tp")
	}else{
		callF.limit = 100
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	ci := getCalls(&callF)
	js, err := json.Marshal(ci)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Write(js)
}

func startServer() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		http.ServeFile(w, r, "static/index.html")
	})

	http.HandleFunc("/find", filterCalls)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

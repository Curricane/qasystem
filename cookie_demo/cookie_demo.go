package main

import (
	"fmt"
	"net/http"
)

func indexHandle(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:   "indexCookie",
		Value:  "helloworld",
		MaxAge: 360000,
		Domain: "localhost",
	}
	http.SetCookie(w, cookie)
	fmt.Printf("set index cookie: %#v\n", cookie)
}

func getCookieHandle(w http.ResponseWriter, r *http.Request) {
	cookies := r.Cookies()
	for idx, cookie := range cookies {
		fmt.Printf("get index:%d cookie: %#v\n", idx, cookie)
	}
}

func main() {
	http.HandleFunc("/", indexHandle)
	http.HandleFunc("/cookie/", getCookieHandle)

	http.ListenAndServe(":9090", nil)
}

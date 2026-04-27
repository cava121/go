package main

import (
	"fmt"
	"net/http"
)

var gitCommit string

func main() {
    http.HandleFunc("/debug/info", getConfig)

	fmt.Println("Сервер запущен")
    http.ListenAndServe(":8090", nil)
}

func getConfig(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Last commit: ", gitCommit);
}

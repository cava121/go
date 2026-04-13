package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

var gitCommit string // Значение подставится при сборке

func main() {
	fmt.Println(gitCommit);
    http.HandleFunc("/debug/info", getConfig)
fmt.Println("Сервер запуще1")
    http.ListenAndServe(":8090", nil)
}

func getConfig(w http.ResponseWriter, req *http.Request) {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			fmt.Printf("Git Commit Hash: %s\n", setting.Value)
		}
	}
}

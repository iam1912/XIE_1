package main

import (
	"log"
	"net/http"

	"github.com/iam1912/XIE_1/control"
)

func main() {
	http.HandleFunc("/index/xjh", control.ViewHandler)
	http.HandleFunc("/edit/xjh", control.EditHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

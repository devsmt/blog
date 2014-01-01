package main

import (
	"log"
	"fmt"
	"net/http"
)

func internalServerErr(w http.ResponseWriter, err error) {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintln(w, "500 Internal Server Error")
        log.Println("internal server error:", err)
}

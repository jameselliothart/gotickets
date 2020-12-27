package main

import (
	"net/http"

	"github.com/jameselliothart/gotickets/infrastructure"
)

func main() {
	infrastructure.Startup()
	http.ListenAndServe(":8000", nil)
}

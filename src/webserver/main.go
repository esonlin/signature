package main

import (
	"fmt"
	"net/http"

	"signature/router"
)

func main() {

	// 监听端口7007
	address := fmt.Sprintf(":%d", 7777)
	fmt.Printf("Listening and serving HTTP on %s\n", address)
	http.ListenAndServe(address, router.Router())
}

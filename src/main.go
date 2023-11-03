package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
    http.HandleFunc("/", getRoot)

    err := http.ListenAndServe(":42069", nil)

    if errors.Is(err, http.ErrServerClosed) {
        fmt.Printf("server closed\n")
    } else if err != nil {
        os.Exit(1)
        fmt.Printf("erro starting server: %s\n", err)
    }
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is my website!\n")

}


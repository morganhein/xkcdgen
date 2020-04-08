package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/morganhein/xkcdgen"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		result, err := xkcdgen.Generate()
		if err != nil {
			_, _ = fmt.Fprintf(w, "error:, %v", err)
			return
		}
		fmt.Println(result)
		_, _ = fmt.Fprintf(w, "%s", result)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}

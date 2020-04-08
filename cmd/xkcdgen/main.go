package main

import (
	"fmt"
	"log"

	"github.com/morganhein/xkcdgen"
)

func main() {
	result, err := xkcdgen.Generate()
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(result)
}

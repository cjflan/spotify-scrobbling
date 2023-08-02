package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cjflan/spotify-scrobbling/controllers"
)

func main() {
	fmt.Print("Hello World")
	http.HandleFunc("/callback", controllers.Callback)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go func() {
		fmt.Print("Starting Server!")
		err := http.ListenAndServe("localhost:6669", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
}

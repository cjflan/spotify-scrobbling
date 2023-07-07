package main

import (
  "fmt"
  "strings"
  "net/http"
)

func main() {
		url := "https://accounts.spotify.com/api/token"
		payload := strings.NewReader(`{
				"grant_type": "client_credentials",
				"client_id": "6d29d24f3b19439a93c732a2832e5341",
				"client_secret": "e32ba7c0b58f448f9dc96859caae7fe2" 
		}`)
		
		client := &http.Client{}
		req, err := http.NewRequest("POST", url, payload)
		
		if err != nil {
				fmt.Println(err)
				return
		}

		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp, err := client.Do(req)

		if err != nil {
				fmt.Println(err)
				return
		}
		fmt.Println(resp)
		return
}



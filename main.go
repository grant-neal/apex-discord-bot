package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	api_key := os.Getenv("API_KEY")
	resp, err := http.Get("https://api.mozambiquehe.re/bridge?version=5&platform=PC&player=" + os.Args[1] + "&auth=" + api_key)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
}

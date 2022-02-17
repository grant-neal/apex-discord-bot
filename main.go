package main

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	val := os.Getenv("API_KEY")
	fmt.Println(val)
}

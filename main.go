package main

import (
	"fmt"
	"mainmodule/helper"
	"mainmodule/service"

	"log"
	"os"
)

func main() {
	result := helper.Reduce(2, 1)
	fmt.Println("result::", result)

	service.Get("https://haokur.com")

	log.Print("hello this is a message")
	log.Printf("her name is %s, her age is %d,her weight is %f", "jack", 18, 55.66)

	file, _ := os.Create("output.txt")
	defer file.Close()
	fmt.Fprintln(file, "hello,file")
}

package main

import (
	"fmt"
	"mainmodule/helper"
	"mainmodule/service"
)

func main() {
	result := helper.Reduce(2, 1)
	fmt.Println("result::", result)

	service.Get("https://haokur.com")
}

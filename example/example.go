package main

import (
	"fmt"
	"github.com/mtsanderson/huego"
)

func main() {
	b := huego.NewHueBridge("192.168.1.19")

	fmt.Println(b)
}

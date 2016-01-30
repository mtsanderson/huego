package main

import (
	"fmt"
	"github.com/mtsanderson/huego"
	"strconv"
)

func main() {
	b := huego.NewHueBridge("192.168.1.19")

	/*
		l := b.Getlight(17)

		fmt.Println(l.State)
		l.State.On = false
		l.SetLightState()
	*/

	groups := b.GetLightGroups()

	for _, g := range groups {
		if g.Name == "notify" {
			for _, l := range g.Lights {
				lid, err := strconv.Atoi(l)
				if err != nil {
					panic(err)
				}
				light := b.Getlight(lid)
				light.On(false)
				fmt.Println(light.State)
			}
		}
	}
}

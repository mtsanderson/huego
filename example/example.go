package main

import (
	"fmt"
	"github.com/mtsanderson/huego"
	"math/rand"
	"strings"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	b := huego.NewHueBridge("192.168.1.19")

	l := b.Getlights()

	done := make(chan bool)

	for _, v := range l {
		if strings.Contains(v.Name, "Living") {
			if v.State.On {
				go lightparty(v)
			}
		}
	}
	<-done
}

func lightparty(l huego.Light) {

	for {
		n := rand.Intn(64000)
		l.Hue(uint16(n))
		fmt.Printf("Set light %d to %d hue!\n", l.Id, n)
		time.Sleep(time.Microsecond * 500)
	}

}

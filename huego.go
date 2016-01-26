package huego

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
)

func perror(err error) {
	// Panic on error
	if err != nil {
		panic(err)
	}
}

type LightState struct {
	On        bool       `json:"on,omitempty"`
	Bri       uint8      `json:"bri"`
	Hue       uint16     `json:"hue"`
	Sat       uint8      `json:"sat"`
	Xy        [2]float32 `json:"xy"`
	ct        uint16     `json:"ct"`
	Alert     string     `json:"alert"`
	Effect    string     `json:"effect"`
	Colormode string     `json:"colormode"`
	Reachable bool       `json:"reachable"`
}

type Light struct {
	Id               int              `json:"id"`
	State            LightState       `json:"state"`
	Type             string           `json:"type"`
	Name             string           `json:"name"`
	Modelid          string           `json:"modelid"`
	Manufacturername string           `json:"manufacturername"`
	Uniqueid         string           `json:"uniqueid"`
	Swversion        string           `json:"swversion"`
	Pointsymbol      LightPointsymbol `json:"pointsymbol"`
}

type LightPointsymbol struct {
	Num1 string `json:"1"`
	Num2 string `json:"2"`
	Num3 string `json:"3"`
	Num4 string `json:"4"`
	Num5 string `json:"5"`
	Num6 string `json:"6"`
	Num7 string `json:"7"`
	Num8 string `json:"8"`
}

type Bridge struct {
	ip         string
	auth_token string
}

func NewHueBridge(ip string) *Bridge {
	// Constructs a new hue bridge

	b := new(Bridge)
	b.ip = ip

	u, err := user.Current()
	perror(err)

	file_path := fmt.Sprintf("%s/.huego", u.HomeDir)

	var fh *os.File
	fh, err = os.Open(file_path)
	if os.IsNotExist(err) {
		fh, err = os.Create(file_path)
		perror(err)
		b.register()
	} else {
		data, err := ioutil.ReadFile(file_path)
		perror(err)

		decoder := json.NewDecoder(bytes.NewReader(data))
		var jsondata map[string]interface{}
		err = decoder.Decode(&jsondata)
		perror(err)

		b.auth_token = jsondata["username"].(string)
	}
	fh.Close()

	return b
}

func (b *Bridge) request(method string, url string, data []byte) *http.Response {
	var resp *http.Response
	var err error

	switch method {
	case "GET":
		resp, err = http.Get(url)
		perror(err)
	case "POST":
		buf := bytes.NewReader(data)
		resp, err = http.Post(url, "application/json", buf)
		perror(err)
	}

	return resp
}

func (b *Bridge) register() {
	// This method registers the library with the Hue Bridge

	reqdata := map[string]string{"devicetype": "huego"}
	jsondata, err := json.Marshal(reqdata)
	perror(err)

	url := fmt.Sprintf("http://%s/api", b.ip)

	resp := b.request("POST", url, jsondata)
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var data []map[string]interface{}
	err = decoder.Decode(&data)
	perror(err)

	for _, line := range data {
		for key, val := range line {
			if key == "error" {
				fmt.Println("Error!")
			} else if key == "success" {
				fmt.Println("Success!")
				fmt.Println(val)
			}
		}
	}
}

func (b *Bridge) getlight(id int) Light {
	// This method will return a Light Object
	url := fmt.Sprintf("http://%s/api/%s/lights/%d", b.ip, b.auth_token, id)
	resp := b.request("GET", url, nil)
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var light Light
	err := decoder.Decode(&light)
	perror(err)

	light.Id = id

	return light
}

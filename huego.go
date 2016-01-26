package huego

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"strconv"
)

func perror(err error) {
	// Panic on error
	if err != nil {
		panic(err)
	}
}

type Light struct {
	Bridge           Bridge
	Id               int              `json:"id,omitempty"`
	State            LightState       `json:"state,omitempty"`
	Type             string           `json:"type,omitempty"`
	Name             string           `json:"name,omitempty"`
	Modelid          string           `json:"modelid,omitempty"`
	Manufacturername string           `json:"manufacturername,omitempty"`
	Uniqueid         string           `json:"uniqueid,omitempty"`
	Swversion        string           `json:"swversion,omityempty"`
	Pointsymbol      LightPointsymbol `json:"pointsymbol,omitempty"`
}

type LightState struct {
	On        bool       `json:"on,omitempty"`
	Bri       uint8      `json:"bri,omitempty"`
	Hue       uint16     `json:"hue,omitempty"`
	Sat       uint8      `json:"sat,omitempty"`
	Xy        [2]float32 `json:"xy,omitempty"`
	ct        uint16     `json:"ct,omitempty"`
	Alert     string     `json:"alert,omitempty"`
	Effect    string     `json:"effect,omitempty"`
	Colormode string     `json:"colormode,omitempty"`
	Reachable bool       `json:"reachable,omitempty"`
}

type LightPointsymbol struct {
	Num1 string `json:"1,omitempty"`
	Num2 string `json:"2,omitempty"`
	Num3 string `json:"3,omitempty"`
	Num4 string `json:"4,omitempty"`
	Num5 string `json:"5,omitempty"`
	Num6 string `json:"6,omitempty"`
	Num7 string `json:"7,omitempty"`
	Num8 string `json:"8,omitempty"`
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
	case "PUT":
		client := &http.Client{}
		buf := bytes.NewReader(data)
		req, err := http.NewRequest("PUT", url, buf)
		if err != nil {
			panic(err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err = client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
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

func (b *Bridge) Getlight(id int) Light {
	// This method will return a Light Object
	url := fmt.Sprintf("http://%s/api/%s/lights/%d", b.ip, b.auth_token, id)
	resp := b.request("GET", url, nil)
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var light Light
	err := decoder.Decode(&light)
	perror(err)

	light.Id = id
	light.Bridge = *b

	return light
}

func (b *Bridge) Getlights() []Light {
	url := fmt.Sprintf("http://%s/api/%s/lights", b.ip, b.auth_token)
	resp := b.request("GET", url, nil)
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var data map[string]interface{}
	err := decoder.Decode(&data)
	if err != nil {
		panic(err)
	}

	var lights []Light

	for k, v := range data {
		// TODO: make this more efficient!!
		var light Light
		jsondata, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(jsondata, &light)
		if err != nil {
			panic(err)
		}

		id, err := strconv.Atoi(k)
		if err != nil {
			panic(err)
		}

		light.Id = id
		light.Bridge = *b

		lights = append(lights, light)
	}

	return lights
}

func (l *Light) On(state bool) {
	l.State.On = state
	url := fmt.Sprintf("http://%s/api/%s/lights/%d/state", l.Bridge.ip, l.Bridge.auth_token, l.Id)

	data := map[string]bool{"on": state}
	jdata, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	l.Bridge.request("PUT", url, jdata)
}

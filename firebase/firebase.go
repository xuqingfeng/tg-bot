package firebase
import (
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"strings"
	"heroku.com/tg-bot/common"
)

type Request struct {
	Request int64 `json:"request"`
}

var config common.Config

func init(){

	config = common.GetConfig()
}

func GetRequest() (request int64) {

	request = 1

	client := &http.Client{}

	req, err := http.NewRequest("GET", config.FirebaseUrl+"request.json", nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("firebase: get data fail")
	}else {
		if resp.Body == nil {
			log.Println("firebase: data is null")
			// create request
			iniRequest := Request{1}
			reqInByte, _ := json.Marshal(iniRequest)
			log.Printf("request %#v", string(reqInByte))
			req4Post, err := http.NewRequest("POST", config.FirebaseUrl+"request.json", strings.NewReader(string(reqInByte)))
			resp4Post, err := client.Do(req4Post)
			defer resp4Post.Body.Close()
			if err != nil {
				log.Fatal("firebase: data create fail")
			}else {
				log.Println(resp4Post.Header)
				body, _ := ioutil.ReadAll(resp4Post.Body)
				log.Printf("resp4Post body: %v", string(body))
				log.Print("firebase: data create success")
			}
		}else {
			body, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			var req Request
			json.Unmarshal(body, &req)
			request = req.Request
		}

	}

	return
}

func IncRequest() (ok bool) {

	client := &http.Client{}

	ok = false
	originRequest := GetRequest()
	currentRequest := originRequest + 1
	currentReq := Request{currentRequest}
	reqInByte, _ := json.Marshal(currentReq)
	req, err := http.NewRequest("PUT", config.FirebaseUrl+"request.json", strings.NewReader(string(reqInByte)))
	client.Do(req)
	if err != nil {
		log.Fatal("firebase: update request fail")
	}else {
		log.Println("firebase: update request success")
		ok = true
	}
	return
}
package biubot
import (
	"io/ioutil"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"heroku.com/tg-bot/util"
)

type Config struct {
	Domain     string
	Api_url    string
	Token      string
	Google_top int
}

/*
https://core.telegram.org/bots/api
*/
type Update struct {
	Update_id int
	Message   Message
}

type User struct {
	Id         int
	First_name string
	Last_name  string
	Username   string
}

type GroupChat struct {
	Id    int
	Title string
}

type Message struct {
	Message_id            int
	From                  User
	Date                  int
	Chat                  User // or group?
	Forward_from          User
	Forward_date          int
	Reply_to_message      *Message
	Text                  string
	Audio                 Audio
	Document              Document
	Photo                 []PhotoSize
	Sticker               Sticker
	Video                 Video
	Contact               Contact
	Location              Location
	New_chat_participant  User
	Left_chat_participant User
	New_chat_title        string
	New_chat_photo        []PhotoSize
	Delete_chat_photo     bool
	Group_chat_created    bool
	Caption               string
}

type PhotoSize struct {
	File_id   string
	Width     int
	Height    int
	File_size int
}

type Audio struct {
	File_id   string
	Duration  int
	Mime_type string
	File_size int
}

type Document struct {
	File_id   string
	Thumb     PhotoSize
	File_name string
	Mime_type string
	File_size int
}

type Sticker struct {
	File_id   string
	Width     int
	Height    int
	Thumb     PhotoSize
	File_size int
}

type Video struct {
	File_id   string
	Width     int
	Height    int
	Duration  int
	Thumb     PhotoSize
	Mime_type string
	File_size int
}

type Contact struct {
	Phone_number string
	First_name   string
	Last_name    string
	User_id      int
}

type Location struct {
	Longitude float64
	Latitude  float64
}

type UserProfilePhotos struct {
}
type ReplyKeyboardMarkup struct {
}
type ReplyKeyboardHide struct {
}
type ForceReply struct {
}

var config Config

func init() {

	configData, err := ioutil.ReadFile("./config.json")
	util.PanicIf(err)
	json.Unmarshal(configData, &config)
}

func GetConfig() (Config) {
	return config
}

func GetWebHookUrl() (string) {

	return "/"+config.Token
}

func GetUpdates() (updates []Update) {

	resp, err := http.Get(config.Api_url+config.Token+"/getUpdates")
	util.PanicIf(err)
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	type Data struct {
		Ok     bool
		Result []Update
	}
	var data Data
	json.Unmarshal(body, &data)

	updates = data.Result
	return
}

func GetLargestUpdateId(updates []Update) (int) {

	var largestId = 0
	for _, update := range updates {
		if update.Update_id > largestId {
			largestId = update.Update_id
		}
	}

	return largestId
}

func SendMessage(chat_id int, text string) {

	type PostData struct {
		Chat_id int `json:"chat_id"`
		Text    string `json:"text"`
	}
	
	resp, err := http.PostForm(config.Api_url+config.Token+"/sendMessage", url.Values{"chat_id": {strconv.Itoa(chat_id)}, "text": {text}})

	defer resp.Body.Close()
	util.PanicIf(err)
}


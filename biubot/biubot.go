package biubot
import (
	"io/ioutil"
	"encoding/json"
	"heroku.com/tg-bot/util"
	"net/http"
	"net/url"
	"strconv"
	"heroku.com/tg-bot/common"
	"io"
)

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

type File struct {
	File_id   string
	File_size int
	File_path string
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

var config common.Config
func init() {
	config = common.GetConfig()
}

func GetWebHookUrl() (string) {

	// fix bug
	return "/"+config.Token
}

func GetUpdates() (updates []Update) {

	resp, err := http.Get(config.ApiUrl+config.Token+"/getUpdates")
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

//func GetUpdates() (updates string){
//
//	resp, err := http.Get(config.Api_url+config.Token+"/getUpdates")
//	util.PanicIf(err)
//	body, err := ioutil.ReadAll(resp.Body)
//	updates = string(body)
//	resp.Body.Close()
//	return
//}


/*
	strconv.Itoa()
	not
	string()
*/
func SendMessage(chat_id int, text string) {

	resp, err := http.PostForm(config.ApiUrl+config.Token+"/sendMessage", url.Values{"chat_id": {strconv.Itoa(chat_id)}, "text": {text}})

	defer resp.Body.Close()
	util.PanicIf(err)
}

func SendLocation(chat_id int, longitude float64, latitude float64) {

	// 'f'?
	resp, err := http.PostForm(config.ApiUrl+config.Token+"/sendLocation", url.Values{"chat_id": {strconv.Itoa(chat_id)}, "latitude": {strconv.FormatFloat(latitude, 'f', 6, 64)}, "longitude": {strconv.FormatFloat(longitude, 'f', 6, 64)}})
	defer resp.Body.Close()
	util.PanicIf(err)
}

func SendPhoto(chat_id int, photo_id string) {

	resp, err := http.PostForm(config.ApiUrl+config.Token+"/sendPhoto", url.Values{"chat_id": {strconv.Itoa(chat_id)}, "photo": {photo_id}})
	defer resp.Body.Close()
	util.PanicIf(err)
}

func GetFile(file_id string) (filePath string) {

	filePath = ""
	resp, err := http.PostForm(config.ApiUrl+config.Token+"/getFile", url.Values{"file_id": {file_id}})
	if err != nil {
		return
	}
	type FileData struct {
		Ok bool
		Result File
	}
	if resp.StatusCode == 200 {
		respInByte, _ := ioutil.ReadAll(resp.Body)
		var fileData FileData
		err := json.Unmarshal(respInByte, &fileData)
		defer resp.Body.Close()

		if err != nil {
			return
		}else{
			filePath = fileData.Result.File_path
		}
	}

	return
}

func GetFileData(filePath string) (r io.Reader){

	resp, _ := http.Get(config.FileApiUrl + config.Token + "/" + filePath)
	if resp.StatusCode == 200 {
		r = resp.Body
	}
	return
}


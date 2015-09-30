package common
import (
	"io/ioutil"
	"encoding/json"
)

type Config struct {
	Domain          string `json:"DOMAIN"`
	ApiUrl          string `json:"API_URL"`
	FileApiUrl      string `json:"FILE_API_URL"`
	Token           string `json:"TOKEN"`
	FirebaseUrl     string `json:"FIREBASE_URL"`
	GoogleTop       int `json:"GOOGLE_TOP"`
	QiniuAccessKey  string `json:"QINIU_ACCESS_KEY"`
	QiniuSecretKey  string `json:"QINIU_SECRET_KEY"`
	QiniuDomain     string `json:"QINIU_DOMAIN"`
	QiniuBucketName string `json:"QINIU_BUCKET_NAME"`
}

var config Config

func init() {

	configData, _ := ioutil.ReadFile("./config.json")
	json.Unmarshal(configData, &config)
}

func GetConfig() (Config) {

	return config
}

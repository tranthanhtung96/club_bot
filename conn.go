package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"runtime"
	"time"

	"github.com/mattermost/mattermost-server/v5/model"
)

const (
	// UserEmail is the email of user loged in as a bot
	UserEmail = "tranthanhtung@ginno.com"
	// UserPassword is the password of this account
	UserPassword = "965148thanhtung"

	// ServerAddr is the address of the mattermost server
	ServerAddr = "https://chat.ginno.com"
	// ChannelID is the id of the club channel
	ChannelID = "rz7c6rzkpi84jgi6ywbt7it35r"
	// AdminID is the id of admin 's chat channel
	AdminID = "k3otob78efn5zeg936fhfoqsce"
)

var client *model.Client4
var botUser *model.User
var dailyPosts []*model.Post

// resDir is a directory contains yesterday posts
var resDir string

// LoginAsTheBotUser log in to the mattermost account
func LoginAsTheBotUser() {
	client = model.NewAPIv4Client(ServerAddr)
	if user, resp := client.Login(UserEmail, UserPassword); resp.Error != nil {
		println("Login unsuccessfully")
		os.Exit(1)
	} else {
		println("Login successfully")
		botUser = user
	}
}

// GetDailyPosts get all daily posts
func GetDailyPosts() {
	loc, _ := time.LoadLocation("Asia/Saigon")
	year, month, day := time.Now().Date()
	todayEpoch := time.Date(year, month, day, 0, 0, 0, 0, loc).Unix() * 1000
	// todayEpoch := int64(1588698000000)
	yesterdayEpoch := todayEpoch - 86400000

	if posts, resp := client.GetPostsSince(ChannelID, yesterdayEpoch); resp.Error != nil {
		println("Get daily posts unsuccessfully")
		os.Exit(1)
	} else {
		dailyPosts = posts.ToSlice()

		_, curGoFile, _, _ := runtime.Caller(0)
		curDir := path.Dir(curGoFile)
		year, month, day := time.Unix(yesterdayEpoch/1000, 0).Date()
		resDir = curDir + fmt.Sprintf("/%04d-%s-%02d/", year, month.String(), day)
		os.RemoveAll(resDir)
		os.MkdirAll(resDir, 0777)
		os.RemoveAll(curDir + "/phrases.txt")

		var phrasesFile1, phrasesFile2 *os.File
		var phrasesErr error

		if phrasesFile1, phrasesErr = os.OpenFile(resDir+"phrases.txt", os.O_CREATE|os.O_WRONLY, 0777); phrasesErr != nil {
			println("Write to phrases.txt unsuccessfully")
		} else {
			defer phrasesFile1.Close()
		}
		if phrasesFile2, phrasesErr = os.OpenFile(curDir+"/phrases.txt", os.O_CREATE|os.O_WRONLY, 0777); phrasesErr != nil {
			println("Write to phrases.txt unsuccessfully")
		} else {
			defer phrasesFile2.Close()
		}
		phrasesFile2.Write([]byte("> "))
		for _, post := range dailyPosts {
			if post.CreateAt > todayEpoch || post.DeleteAt != 0 {
				continue
			}

			re := regexp.MustCompile(`(?m)^\s*[#*]+\s{0,1}(.*)(\n|$)`)
			matches := re.FindAllStringSubmatch(post.Message, -1)
			if matches != nil {
				phraseData := ""
				for _, match := range matches {
					phraseData += match[1] + "\n"
				}
				phrasesFile1.Write([]byte(phraseData))
				phrasesFile2.Write([]byte(phraseData))
				phrasePath := resDir + post.UserId
				phraseName := post.Id + ".txt"
				if err := SaveFile(phrasePath, phraseName, []byte(phraseData)); err != nil {
					println("Save file " + post.Id + " unsuccessfully")
				}
			}

			if post.FileIds != nil {
				for _, fileIDm := range post.FileIds {
					if imgByte, _ := client.GetFile(fileIDm); resp.Error != nil {
						println("Get file " + fileIDm + " unsuccessfully")
					} else {
						imgPath := resDir + post.UserId
						imgType := http.DetectContentType(imgByte)
						imgName := fileIDm + "." + imgType[6:]
						if err := SaveFile(imgPath, imgName, imgByte); err != nil {
							println("Save file " + fileIDm + " unsuccessfully")
						}
					}
				}
			}
		}
		phrasesFile1.Close()
	}
}

// SaveFile save file
func SaveFile(path string, name string, data []byte) error {
	var err error
	var out *os.File
	if _, err = os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, os.FileMode(int(0777))); err != nil {
			return err
		}
	}
	if out, err = os.Create(path + "/" + name); err != nil {
		return err
	}

	if _, err = io.Copy(out, bytes.NewReader(data)); err != nil {
		return err
	}
	defer out.Close()
	return nil
}

// PostResult reads phrases.txt and creat a post to the channel
func PostResult() {
	_, curGoFile, _, _ := runtime.Caller(0)
	curDir := path.Dir(curGoFile)
	if phrasesFile, phrasesErr := os.OpenFile(curDir+"/phrases.txt", os.O_RDONLY, 0777); phrasesErr != nil {
		println("Read the phrases.txt unsuccessfully")
	} else {
		defer phrasesFile.Close()
		bytes, _ := ioutil.ReadAll(phrasesFile)

		post := &model.Post{}
		post.ChannelId = ChannelID
		post.Message = string(bytes)

		if _, resp := client.CreatePost(post); resp.Error != nil {
			println("Post to the group 's channel unsuccessfully")
		}
	}
}

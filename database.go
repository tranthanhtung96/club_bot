package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
)

// ClubMems is the number of the english club members
var ClubMems int

// DBClub contains the club members' info
type DBClub []struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// DBEngClub is the DBClub for the english club
var DBEngClub DBClub

// DBfile is the name of the db file
var DBfile = "clubDB.json"

// DBLoadFromFile loads DBClub.json to DBEngClub
func DBLoadFromFile() {
	if jsonFile, err := os.Open(DBfile); err != nil {
		println("Open db unsuccessfully")
		os.Exit(1)
	} else {
		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &DBEngClub)
		defer jsonFile.Close()
	}
}

// DBCheckOff checks all the images in the folder "images" and return whenever member off or not
func DBCheckOff() {
	// _, curGoFile, _, _ := runtime.Caller(0)
	// curDir := path.Dir(curGoFile)
	offList := ""

	for _, mem := range DBEngClub {
		chkFolder := resDir + mem.ID
		_, err := os.Stat(chkFolder)
		files, _ := ioutil.ReadDir(chkFolder)

		if os.IsNotExist(err) || len(files) < 2 {
			println("name: " + mem.Username + "\tid: " + mem.ID)
			offList += mem.Username + ", "
		} else {
			txtExist := false
			for _, file := range files {
				if strings.Contains(file.Name(), ".txt") {
					txtExist = true
					break
				}
			}
			if !txtExist {
				println("name: " + mem.Username + "\tid: " + mem.ID)
				offList += mem.Username
			}
		}
	}

	post := &model.Post{}
	post.ChannelId = AdminID
	post.Message = offList
	if _, resp := client.CreatePost(post); resp.Error != nil {
		println("Post to the group 's channel unsuccessfully")
	}
}

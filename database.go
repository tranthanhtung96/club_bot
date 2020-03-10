package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

// ClubMems is the number of the english club members
var ClubMems int

// DBClub contains the club members' info
type DBClub []struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	NoOffs   int    `json:"nooffs"`
	Sec      int    `json:"sec"`
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

// DBReset reset the database
func DBReset() {
	DBLoadFromFile()
	for i := range DBEngClub {
		DBEngClub[i].NoOffs = 0
	}
	file, _ := json.MarshalIndent(DBEngClub, "", "\t")
	_ = ioutil.WriteFile(DBfile, file, 0644)

}

// DBSetOff increase the number of offline of a mem who has a given ID
func DBSetOff(id string) {
	for i, mem := range DBEngClub {
		if mem.ID == id {
			DBEngClub[i].NoOffs++
			file, _ := json.MarshalIndent(DBEngClub, "", "\t")
			_ = ioutil.WriteFile(DBfile, file, 0644)
			return
		}
	}
}

// DBCheckOff checks all the images in the folder "images" and return whenever member off or not
func DBCheckOff() {
	// _, curGoFile, _, _ := runtime.Caller(0)
	// curDir := path.Dir(curGoFile)
	for _, mem := range DBEngClub {
		chkFolder := resDir + mem.ID
		_, err := os.Stat(chkFolder)
		files, _ := ioutil.ReadDir(chkFolder)

		if os.IsNotExist(err) || len(files) < 2 {
			DBSetOff(mem.ID)
			println("name: " + mem.Username + "\tid: " + mem.ID)
		} else {
			txtExist := false
			for _, file := range files {
				if strings.Contains(file.Name(), ".txt") {
					txtExist = true
					break
				}
			}
			if !txtExist {
				DBSetOff(mem.ID)
				println("name: " + mem.Username + "\tid: " + mem.ID)
			}
		}
	}
}

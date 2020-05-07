package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
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

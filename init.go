package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/user"
	"path/filepath"
)

// Return json path in userfolder
func initFilePath() (initFile string, err error) {
	usr, err := user.Current()
	initFile = filepath.Join(usr.HomeDir, "imgcvrt_init.json")
	return initFile, err
}

// read data from userfolder
func readInit(initFile string) (*InitData, error) {
	data := InitData{}
	files, err := ioutil.ReadFile(initFile)
	if err != nil {
		return &data, err
	}
	err = json.Unmarshal([]byte(files), &data)
	return &data, err
}

// make json data
func makeInit(initFile string) (*InitData, error) {
	data := InitData{}
	data.OiioTool = ""
	data.Ffmpeg = ""
	data.Ocio = ""
	fmt.Println("Enter oiioTool path: (ex:/usr/local/bin/oiioTool)")
	fmt.Scanln(&data.OiioTool)
	fmt.Println("Enter ffmpeg path: (ex:/usr/local/bin/ffmpeg)")
	fmt.Scanln(&data.Ffmpeg)
	fmt.Println("Enter ocio path: (ex:/Users/user/OpenColorIO-Configs/aces_1.0.3/config.ocio)")
	fmt.Scanln(&data.Ocio)
	files, _ := json.MarshalIndent(data, "", " ")
	err := ioutil.WriteFile(initFile, files, 0644)
	return &data, err
}

// loadInit is get init data from userfolder.
// if not found, make one
func LoadInit() (data *InitData, err error) {
	initFile, err := initFilePath()
	if err != nil {
		return data, err
	}
	if Exists(initFile) {
		data, err = readInit(initFile)
	} else {
		data, err = makeInit(initFile)
	}
	return data, err
}

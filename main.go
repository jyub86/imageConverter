package main

import (
	"log"
)

func main() {
	initData, err := LoadInit()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("OpenImageIO : %s | FFmpeg : %s | OCIO : %s", initData.OiioTool, initData.Ffmpeg, initData.Ocio)
	data := FlagParser()
	data, err = data.Analyze()
	if err != nil {
		log.Fatal(err)
	}
	data, err = data.MakeCmd(initData)
	if err != nil {
		log.Fatal(err)

	}
	err = data.Render(initData)
	if err != nil {
		log.Fatal(err)
	}
}

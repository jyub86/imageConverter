package main

import (
	"errors"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// analyze argument
func (i Item) Analyze() (*Item, error) {
	// check argument data
	if i.Input == "" {
		return &i, errors.New("input not found")
	}
	if i.Output == "" {
		return &i, errors.New("output not found")
	}
	if i.Incolor != "" && i.Outcolor == "" {
		return &i, errors.New("outcolor not found")
	}
	if i.Incolor == "" && i.Outcolor != "" {
		return &i, errors.New("incolor not found")
	}
	// analyze and fill data
	// ext and data type
	i.InputExt = strings.ToLower(filepath.Ext(i.Input))
	switch i.InputExt {
	case ".jpg", ".jpeg", ".png", ".dpx", ".exr", ".tif", ".tiff":
		i.InputType = "image"
	case ".mov", ".mp4":
		i.InputType = "video"
	}
	i.OutputExt = strings.ToLower(filepath.Ext(i.Output))
	switch i.OutputExt {
	case ".jpg", ".jpeg", ".png", ".dpx", ".exr":
		i.OutputType = "image"
	case ".mov", ".mp4":
		i.OutputType = "video"
	}
	// color convert require proxy img path. using png
	if i.Incolor != "" {
		dir := filepath.Dir(i.Output)
		if i.InputType == "image" {
			baseName := filepath.Base(i.Input)
			names := strings.Split(baseName, FindSeqPad(i.Input))
			proxyFolder := filepath.Join(dir, names[0]+"_proxy")
			i.Proxy = filepath.Join(proxyFolder, names[0]+FindSeqPad(i.Input)+".png")
		}
		if i.OutputType == "image" {
			baseName := filepath.Base(i.Output)
			names := strings.Split(baseName, FindSeqPad(i.Output))
			upDir := filepath.Dir(dir)
			proxyFolder := filepath.Join(upDir, names[0]+"_proxy")
			i.Proxy = filepath.Join(proxyFolder, names[0]+FindSeqPad(i.Output)+".png")
		}
	}
	// set frames
	if i.Frames != "" {
		re, _ := regexp.Compile("([0-9]+)-([0-9]+$)")
		results := re.FindStringSubmatch(i.Frames)
		if results == nil {
			return &i, errors.New("frames not found")
		}
		i.FrameIn = results[1]
		// get Frame range
		last, err := strconv.Atoi(results[2])
		if err != nil {
			return &i, err
		}
		first, err := strconv.Atoi(results[1])
		if err != nil {
			return &i, err
		}
		i.FrameRange = strconv.Itoa(last - first + 1)
	}
	// set timecode
	if i.Timecode != "" {
		re, _ := regexp.Compile("([0-9]{2}:[0-9]{2}:[0-9]{2})-([0-9]{2}:[0-9]{2}:[0-9]{2}$)")
		results := re.FindStringSubmatch(i.Timecode)
		if results == nil {
			return &i, errors.New("timecode not found")
		}
		i.StartTime = results[1]
		i.EndTime = results[2]
	}
	return &i, nil
}

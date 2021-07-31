package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// struct for ocio, oiiotool, ffmpeg path
type InitData struct {
	Ocio     string `json:"ocio"`
	OiioTool string `json:"oiiotool"`
	Ffmpeg   string `json:"ffmpeg"`
}

// item struct for input, output data
type Item struct {
	// input, output
	Input  string `json:"input"`
	Output string `json:"output"`
	// Frame
	Frames     string `json:"frames"`
	FrameIn    string `json:"framein"`
	FrameRange string `json:"framerange"`
	// Timecode
	Timecode  string `json:"timecode"`
	StartTime string `json:"starttime"`
	EndTime   string `json:"endtime"`
	// resize
	Resize string `json:"resize"`
	// colorconvert
	Incolor  string `json:"incolor"`
	Outcolor string `json:"outcolor"`
	// video
	Fps    string `json:"fps"`
	Codec  string `json:"codec"`
	PixFmt string `json:"pixfmt"`
	// ext
	InputExt  string `json:"inputext"`
	OutputExt string `json:"outputext"`
	// media type
	InputType  string `json:"inputtype"`
	OutputType string `json:"outputtype"`
	// proxy path
	Proxy string `json:"proxy"`
	// command line
	FirstCmd  []string `json:"firstcmd"`
	SecondCmd []string `json:"secondcmd"`
}

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
func loadInit() (data *InitData, err error) {
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

// flag parser
func flagParser() *Item {
	data := Item{}
	flag.StringVar(&data.Input, "input", "", "input file path ex) input.%04d.jpg")
	flag.StringVar(&data.Frames, "frames", "", "input frames ex) 1001-1010")
	flag.StringVar(&data.Timecode, "timecode", "", "start & end timecode ex) 00:00:30-00:01:30")
	flag.StringVar(&data.Incolor, "incolor", "", "input colorspace ex) 'ACES - ACES2065-1'")
	flag.StringVar(&data.Outcolor, "outcolor", "", "output colorspace ex) 'Output - Rec.709'")
	flag.StringVar(&data.Resize, "resize", "", "output resize ex) 1920x1080")
	flag.StringVar(&data.Output, "output", "", "output file path ex) output.mov")
	flag.StringVar(&data.Fps, "fps", "24", "output video fps ex) 24")
	flag.StringVar(&data.Codec, "codec", "libx264", "output video codec ex) libx264")
	flag.StringVar(&data.PixFmt, "pixfmt", "yuv420p", "output video pixel format ex) yuv420p")
	flag.Parse()
	if len(os.Args) < 3 {
		fmt.Println("Usage: imgConverter -input input.%04d.jpg -frames 1001-1010 -incolor 'ACES - ACES2065-1' -outcolor 'Output - Rec.709' -resize 1920x1080 -out output.mov")
		flag.PrintDefaults()
		os.Exit(0)
	}
	return &data
}

// analyze argument
func (i Item) analyze() (*Item, error) {
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

// make cmd line
// type                          | option                | firstCmd | secondCmd
// img to img                    | frames, color, resize | oiio     |
// img to video w colorconvert   | frames, resize        | ffmpeg   |
// img to video w/o colorconvert | frames, resize        | oiio     | ffmpeg
// video to video                | frames, resize        | ffmpeg   |
// video to img w colorconvert   | frames, resize        | ffmpeg   |
// video to img w/o colorconvert | frames, resize        | ffmpeg   | oiio
func (i Item) makeCmd(d *InitData) (*Item, error) {
	switch {
	case i.InputType == "image" && i.OutputType == "image":
		i.FirstCmd = append(i.FirstCmd, d.OiioTool, i.Input)
		if i.Frames != "" {
			i.FirstCmd = append(i.FirstCmd, "--frames", i.Frames)
		}
		if i.Incolor != "" {
			i.FirstCmd = append(i.FirstCmd, "--colorconvert", i.Incolor, i.Outcolor)
		}
		if i.Resize != "" {
			i.FirstCmd = append(i.FirstCmd, "-resize", i.Resize)
		}
		i.FirstCmd = append(i.FirstCmd, "-o", i.Output)
	case i.InputType == "image" && i.OutputType == "video" && i.Incolor != "":
		i.FirstCmd = append(i.FirstCmd, d.OiioTool, i.Input, "--colorconvert", i.Incolor, i.Outcolor)
		if i.Frames != "" {
			i.FirstCmd = append(i.FirstCmd, "--frames", i.Frames)
		}
		if i.Resize != "" {
			i.FirstCmd = append(i.FirstCmd, "-resize", i.Resize)
		}
		i.FirstCmd = append(i.FirstCmd, "-o", i.Proxy)
		i.SecondCmd = append(i.SecondCmd, d.Ffmpeg, "-loglevel", "error", "-y", "-r", i.Fps, "-f", "image2", "-start_number", i.FrameIn, "-frames:v", i.FrameRange, "-i", i.Proxy, "-c:v", i.Codec, "-pix_fmt", "yuv420p", "-vf", "pad=ceil(iw/2)*2:ceil(ih/2)*2", i.Output)
	case i.InputType == "image" && i.OutputType == "video" && i.Incolor == "":
		i.FirstCmd = append(i.FirstCmd, d.Ffmpeg, "-loglevel", "error", "-y", "-r", i.Fps, "-f", "image2", "-start_number", i.FrameIn, "-frames:v", i.FrameRange, "-i", i.Input, "-c:v", i.Codec, "-pix_fmt", "yuv420p")
		if i.Resize != "" {
			i.FirstCmd = append(i.FirstCmd, "-vf", "'scale="+i.Resize+"'")
		} else {
			i.FirstCmd = append(i.FirstCmd, "-vf", "pad=ceil(iw/2)*2:ceil(ih/2)*2")
		}
		i.FirstCmd = append(i.FirstCmd, i.Output)
	case i.InputType == "video" && i.OutputType == "video":
		i.FirstCmd = append(i.FirstCmd, d.Ffmpeg, "-loglevel", "error", "-y", "-i", i.Input, "-c:v", i.Codec, "-pix_fmt", "yuv420p")
		if i.Frames != "" {
			i.FirstCmd = append(i.FirstCmd, "-start_number", i.FrameIn, "-frames:v", i.FrameRange)
		}
		if i.Timecode != "" {
			i.FirstCmd = append(i.FirstCmd, "-ss", i.StartTime, "-t", i.EndTime)
		}
		if i.Resize != "" {
			i.FirstCmd = append(i.FirstCmd, "-vf", "'scale="+i.Resize+"'")
		} else {
			i.FirstCmd = append(i.FirstCmd, "-vf", "pad=ceil(iw/2)*2:ceil(ih/2)*2")
		}
		i.FirstCmd = append(i.FirstCmd, i.Output)
	case i.InputType == "video" && i.OutputType == "image" && i.Incolor != "":
		i.FirstCmd = append(i.FirstCmd, d.Ffmpeg, "-loglevel", "error", "-y", "-i", i.Input, "-r", i.Fps)
		if i.Frames != "" {
			i.FirstCmd = append(i.FirstCmd, "-start_number", i.FrameIn, "-frames:v", i.FrameRange)
		}
		if i.Timecode != "" {
			i.FirstCmd = append(i.FirstCmd, "-ss", i.StartTime, "-t", i.EndTime)
		}
		if i.Resize != "" {
			i.FirstCmd = append(i.FirstCmd, "-vf", "scale="+i.Resize)
		} else {
			i.FirstCmd = append(i.FirstCmd, "-vf", "pad=ceil(iw/2)*2:ceil(ih/2)*2")
		}
		i.FirstCmd = append(i.FirstCmd, i.Proxy)
		i.SecondCmd = append(i.SecondCmd, d.OiioTool, i.Proxy, "--colorconvert", i.Incolor, i.Outcolor, "-o", i.Output)
	case i.InputType == "video" && i.OutputType == "image" && i.Incolor == "":
		i.FirstCmd = append(i.FirstCmd, d.Ffmpeg, "-loglevel", "error", "-y", "-i", i.Input, "-r", i.Fps)
		if i.Frames != "" {
			i.FirstCmd = append(i.FirstCmd, "-start_number", i.FrameIn, "-frames:v", i.FrameRange)
		}
		if i.Timecode != "" {
			i.FirstCmd = append(i.FirstCmd, "-ss", i.StartTime, "-t", i.EndTime)
		}
		if i.Resize != "" {
			i.FirstCmd = append(i.FirstCmd, "-vf", "scale="+i.Resize)
		} else {
			i.FirstCmd = append(i.FirstCmd, "-vf", "pad=ceil(iw/2)*2:ceil(ih/2)*2")
		}
		i.FirstCmd = append(i.FirstCmd, i.Output)
	}
	return &i, nil
}

// execute cmd
func (i Item) render(d *InitData) error {
	// make proxy dir
	if i.Proxy != "" && !Exists(filepath.Dir(i.Proxy)) {
		err := os.MkdirAll(filepath.Dir(i.Proxy), os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
	// make output Dir
	if !Exists(filepath.Dir(i.Output)) {
		err := os.MkdirAll(filepath.Dir(i.Output), os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println(strings.Join(i.FirstCmd, " "))
	cmd := exec.Command(i.FirstCmd[0], i.FirstCmd[1:]...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "OCIO="+d.Ocio)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	if i.SecondCmd != nil {
		log.Println(strings.Join(i.SecondCmd, " "))
		cmd := exec.Command(i.SecondCmd[0], i.SecondCmd[1:]...)
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "OCIO="+d.Ocio)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}
		if err := cmd.Wait(); err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

// find seq pad return. ex) .1001 or .%04d
func FindSeqPad(str string) string {
	re, _ := regexp.Compile(".+([_.][0-9]+)(.[a-zA-Z]+$)")
	results := re.FindStringSubmatch(str)
	if results != nil {
		return results[1]
	}
	re, _ = regexp.Compile(".+([_.]%[0-9]+d)(.[a-zA-Z]+$)")
	results = re.FindStringSubmatch(str)
	if results != nil {
		return results[1]
	}
	return ""
}

// check file exists
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func main() {
	initData, err := loadInit()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("load init data")
	data := flagParser()
	data, err = data.analyze()
	if err != nil {
		log.Fatal(err)
	}
	data, err = data.makeCmd(initData)
	if err != nil {
		log.Fatal(err)

	}
	err = data.render(initData)
	if err != nil {
		log.Fatal(err)
	}
}

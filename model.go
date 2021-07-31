package main

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

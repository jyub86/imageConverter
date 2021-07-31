package main

// make cmd line
// type                          | option                | firstCmd | secondCmd
// img to img                    | frames, color, resize | oiio     |
// img to video w colorconvert   | frames, resize        | ffmpeg   |
// img to video w/o colorconvert | frames, resize        | oiio     | ffmpeg
// video to video                | frames, resize        | ffmpeg   |
// video to img w colorconvert   | frames, resize        | ffmpeg   |
// video to img w/o colorconvert | frames, resize        | ffmpeg   | oiio
func (i Item) MakeCmd(d *InitData) (*Item, error) {
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
			i.FirstCmd = append(i.FirstCmd, "-vf", "scale="+i.Resize)
		} else {
			i.FirstCmd = append(i.FirstCmd, "-vf", "pad=ceil(iw/2)*2:ceil(ih/2)*2")
		}
		i.FirstCmd = append(i.FirstCmd, i.Output)
	case i.InputType == "video" && i.OutputType == "video":
		i.FirstCmd = append(i.FirstCmd, d.Ffmpeg, "-loglevel", "error", "-y", "-r", i.Fps, "-i", i.Input, "-c:v", i.Codec, "-pix_fmt", "yuv420p")
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
	case i.InputType == "video" && i.OutputType == "image" && i.Incolor != "":
		i.FirstCmd = append(i.FirstCmd, d.Ffmpeg, "-loglevel", "error", "-y", "-r", i.Fps, "-i", i.Input)
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
		i.FirstCmd = append(i.FirstCmd, d.Ffmpeg, "-loglevel", "error", "-y", "-r", i.Fps, "-i", i.Input)
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

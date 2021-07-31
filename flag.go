package main

import (
	"flag"
	"fmt"
	"os"
)

// flag parser
func FlagParser() *Item {
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

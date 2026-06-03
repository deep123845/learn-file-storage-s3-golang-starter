package main

import (
	"bytes"
	"encoding/json"
	"math"
	"os/exec"
)

type Probe struct {
	Streams []struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	}
}

func getVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	var outDimensions Probe
	err = json.Unmarshal(out.Bytes(), &outDimensions)
	if err != nil {
		return "", err
	}

	ratio := float64(outDimensions.Streams[0].Width) / float64(outDimensions.Streams[0].Height)
	horizontalRatio := 16.0 / 9.0
	verticalRatio := 9.0 / 16.0

	if math.Abs(ratio-horizontalRatio) <= 0.05 {
		return "16:9", nil
	} else if math.Abs(ratio-verticalRatio) <= 0.05 {
		return "9:16", nil
	} else {
		return "other", nil
	}
}

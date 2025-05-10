package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

func getVideoAspectRatio(filepath string) (string, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-print_format", "json",
		"-show_streams",
		filepath,
	)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("ffprobe error: %v", err)
	}

	var videoData struct {
		Streams []struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"streams"`
	}

	err = json.Unmarshal(out.Bytes(), &videoData)
	if err != nil {
		return "", fmt.Errorf("could not parse ffprobe output: %v", err)
	}

	if len(videoData.Streams) == 0 {
		return "", fmt.Errorf("could not parse ffprobe output: %v", err)
	}

	ratio := float64(videoData.Streams[0].Width) / float64(videoData.Streams[0].Height)

	ratio_string := "other"
	if ratio >= 1.7 && ratio <= 1.8 {
		ratio_string = "16:9"
	} else if ratio >= 0.5 && ratio <= 0.6 {
		ratio_string = "9:16"
	}

	return ratio_string, nil
}

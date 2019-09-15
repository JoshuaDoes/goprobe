package goprobe //ffprobe for Go

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
)

// Probe stores data returned from ffprobe
type Probe struct {
	Streams []*Stream
	Format  *Format `json:"format"`
}

// Stream stores data about a media's stream
type Stream struct {
	Index              int    `json:"index"`
	CodecName          string `json:"codec_name"`
	CodecLongName      string `json:"codec_long_name"`
	Profile            string `json:"profile"`
	CodecType          string `json:"codec_type"`
	CodecTimeBase      string `json:"codec_time_base"`
	CodecTagString     string `json:"codec_tag_string"`
	CodecTag           string `json:"codec_tag"`
	Width              int    `json:"width"`
	Height             int    `json:"heigth"`
	CodecWidth         int    `json:"codec_width"`
	CodecHeight        int    `json:"codec_height"`
	HasBFrames         int    `json:"has_b_frames"`
	SampleAspectRatio  string `json:"sample_aspect_ratio"`
	DisplayAspectRatio string `json:"display_aspect_ratio"`
	PixelFormat        string `json:"pix_fmt"`
	Level              int    `json:"level"`
	ColorRange         string `json:"color_range"`
	ColorSpace         string `json:"color_space"`
	ColorTransfer      string `json:"color_transfer"`
	ColorPrimaries     string `json:"color_primaries"`
	ChromaLocation     string `json:"chroma_location"`
	Refs               int    `json:"refs"`
	IsAVC              string `json:"is_avc"`
	NalLengthSize      string `json:"nal_length_size"`
	FrameRateRatio     string `json:"r_frame_rate"`
	FrameRateAverage   string `json:"avg_frame_rate"`
	TimeBase           string `json:"time_base"`
	StartPts           int    `json:"start_pts"`
	StartTime          string `json:"start_time"`
	DurationTS         int    `json:"duration_ts"`
	Duration           string `json:"duration"`
	Bitrate            string `json:"bit_rate"`
	BitsPerRawSample   string `json:"bits_per_raw_sample"`
	FrameCount         string `json:"nb_frames"`

	Disposition *Disposition `json:"disposition"`
	Tags        *Tags        `json:"tags"`
}

// Disposition stores data about a media's disposition
type Disposition struct {
	Default         int `json:"default"`
	Dub             int `json:"dub"`
	Original        int `json:"original"`
	Comment         int `json:"comment"`
	Lyrics          int `json:"lyrics"`
	Karaoke         int `json:"karaoke"`
	Forced          int `json:"forced"`
	HearingImpaired int `json:"hearing_impaired"`
	VisualImpaired  int `json:"visual_impaired"`
	CleanEffects    int `json:"clean_effects"`
	AttachedPic     int `json:"attached_pic"`
	TimedThumbnails int `json:"timed_thumbnails"`
}

// Tags stores data about a media's tags
type Tags struct {
	Title            string `json:"title"`
	Artist           string `json:"artist"`
	Album            string `json:"album"`
	Track            string `json:"track"`
	Comment          string `json:"comment"`
	Genre            string `json:"genre"`
	Date             string `json:"date"` //Usually the year
	MajorBrand       string `json:"major_brand"`
	MinorVersion     string `json:"major_version"`
	CompatibleBrands string `json:"compatible_brands"`
	CreationTime     string `json:"creation_time"`
	Language         string `json:"language"`
	HandlerName      string `json:"handler_name"`
}

// Format stores data about a media's format
type Format struct {
	Filename       string `json:"filename"`
	NumStreams     int    `json:"nb_streams"`
	NumPrograms    int    `json:"nb_programs"`
	FormatName     string `json:"format_name"`
	FormatLongName string `json:"format_long_name"`
	StartTime      string `json:"start_time"`
	Duration       string `json:"duration"`
	Size           string `json:"size"`
	Bitrate        string `json:"bit_rate"`
	ProbeScore     int    `json:"probe_score"`

	Tags *Tags `json:"tags"`
}

// ProbeMedia returns probe data from ffprobe
func ProbeMedia(filename string) (*Probe, error) {
	var buffer bytes.Buffer
	var probe *Probe

	ffprobe := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", "-show_error", "-show_chapters", filename)
	ffprobe.Stdout = &buffer

	err := ffprobe.Start()
	if err != nil {
		return nil, errors.New("goprobe: could not start ffprobe process")
	}

	err = ffprobe.Wait()
	if err != nil {
		return nil, errors.New("goprobe: ffprobe ended unexpectedly: " + fmt.Sprintf("%v", err))
	}

	err = json.Unmarshal(buffer.Bytes(), &probe)
	if err != nil {
		return nil, errors.New("goprobe: could not unmarshal JSON: " + fmt.Sprintf("%v", err))
	}

	if probe.Format == nil {
		probe.Format = &Format{}
	}
	if probe.Format.Tags == nil {
		probe.Format.Tags = &Tags{}
	}

	return probe, nil
}

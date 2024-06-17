package main

import (
	"encoding/json"
	"errors"
	"image"
)

var (
	ErrNotSupportJsonType = errors.New("not support json type")
	ErrNotSupportFileType = errors.New("not support file type")
)

type JsonSize struct {
	W int `json:"w"`
	H int `json:"h"`
}

type JsonRect struct {
	W int `json:"w"`
	H int `json:"h"`
	X int `json:"x"`
	Y int `json:"y"`
}

type JsonTextures struct {
	Image string `json:"image"`
	Foramt string `json:"format"`
	scale int `json:"scale"`
	size JsonSize `json:"size"`
	Frames []*JsonFrameV3 `json:"frames"`
}

type JsonMetaData struct {
	Version string `json:"version"`
	App string `json:"app"`
	SmartUpdate string `json:"smartupdate"`
}

type JsonStruct struct {
	Textures []*JsonTextures `json:"textures"`
	Meta *JsonMetaData `json:"meta"`
}

type JsonFrameV3 struct {
	Frame            *JsonRect `json:"frame"`
	Rotated          bool      `json:"rotated"`
	Trimmed          bool      `json:"trimmed"`
	SpriteSourceSize *JsonRect `json:"spriteSourceSize"`
	SourceSize       *JsonSize `json:"sourceSize"`
	Filename         string    `json:"filename"`
}

func dumpJson(c *DumpContext) error {

	jsonStruct := JsonStruct{}
	err := json.Unmarshal(c.FileContent, &jsonStruct)
	if err != nil {
		return err
	}

	if jsonStruct.Meta == nil {
		return ErrNotSupportJsonType
	}

	part := c.AppendPart()

	part.ImageFile = jsonStruct.Textures[0].Image

	frames := map[string]*JsonFrameV3{}

	for _, v := range jsonStruct.Textures[0].Frames {
		frames[v.Filename] = v
	}

	for k, v := range frames {
		f := v.Frame
		s := v.SourceSize
		part.Frames[k] = &Frame{
			Rect:         image.Rect(f.X, f.Y, f.X+f.W, f.Y+f.H),
			OriginalSize: image.Point{s.W, s.H},
			Rotated:      ifelse(v.Rotated, 90, 0),
			Offset:       image.Point{-v.SpriteSourceSize.X / 2, -v.SpriteSourceSize.Y / 2}, //plist offset in center, json in left-top
		}
	}

	return nil
}

package youtubets

import (
	"encoding/xml"
	"fmt"
	"io"
)

func FetchTransscript(video_id string, lang string, name string) error {

	var ts Transcript
	if err := xml.Unmarshal(bytes, data); err != nil {
		fmt.Println("XML Unmarshal error:", err)
		return err
	}
	for _, text := range data.Text {
		fmt.Println(text + "  ")
	}
	return nil
}

func request(io.Writer) error {
	return nil
}

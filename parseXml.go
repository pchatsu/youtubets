package youtubets

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

type Transcript struct {
	Text []string `xml:"text"`
}

func main() {
	bytes, err := ioutil.ReadFile("/Users/sunagawa/text.xml")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	data := new(Transcript)
	if err := xml.Unmarshal(bytes, data); err != nil {
		fmt.Println("XML Unmarshal error:", err)
		return
	}
	for _, text := range data.Text {
		fmt.Println(text + "  ")
	}
}

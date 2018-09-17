package youtubets

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

var (
	Domain = "http://video.google.com"
)

type Video struct {
	ID string
}

type TranscriptList struct {
	List []Transcript `xml:"track"`
}

type Transcript struct {
	VideoID string
	Name    string `xml:"name,attr"`
	Lang    string `xml:"lang_code,attr"`
	Default bool   `xml:"lang_default,attr"`
}

type Text struct {
	Lines []string `xml:"text"`
}

func (v *Video) TranscriptList() ([]Transcript, error) {
	req, err := http.NewRequest("GET", Domain+"/timedtext", nil)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	q := url.Values{}
	q.Add("v", v.ID)
	q.Add("type", "list")
	req.URL.RawQuery = q.Encode()

	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request track list")
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response %d in requesting transcript", res.StatusCode)
	}

	ts, err := NewTranscriptList(v.ID, res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get transcript list")
	}
	return ts.List, nil
}

func NewTranscriptList(videoID string, xmlReader io.Reader) (*TranscriptList, error) {
	var tl TranscriptList
	var buf bytes.Buffer

	if _, err := buf.ReadFrom(xmlReader); err != nil {
		return nil, errors.Wrap(err, "read transcript list XML error")
	}

	if buf.Len() <= 0 {
		return nil, errors.New("has empty transcript list XML")
	}

	if err := xml.Unmarshal(buf.Bytes(), &tl); err != nil {
		return nil, errors.Wrap(err, "invalid XML")
	}

	for i := range tl.List {
		tl.List[i].VideoID = videoID
	}

	return &tl, nil
}

func (t *Transcript) Text() (*Text, error) {
	req, err := http.NewRequest("GET", Domain+"/timedtext", nil)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	q := url.Values{}
	q.Add("v", t.VideoID)
	q.Add("lang", t.Lang)
	q.Add("name", t.Name)
	req.URL.RawQuery = q.Encode()

	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request transcript")
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response %d in requesting transcript", res.StatusCode)
	}

	return NewText(res.Body)
}

func NewText(xmlReader io.Reader) (*Text, error) {
	var text Text
	var buf bytes.Buffer

	if _, err := buf.ReadFrom(xmlReader); err != nil {
		return nil, errors.Wrap(err, "can't read transcript XML")
	}

	if buf.Len() <= 0 {
		return nil, errors.New("has empty transcript XML")
	}

	if err := xml.Unmarshal(buf.Bytes(), &text); err != nil {
		return nil, errors.Wrap(err, "failed to parse XML")
	}

	return &text, nil
}

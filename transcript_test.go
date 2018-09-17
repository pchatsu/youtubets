package youtubets_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/pchatsu/youtubets"
)

func TestVideo_TranscriptList(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/timedtext" {
				t.Fatalf("invalid List API Path %s", r.URL.Path)
			}
			if r.URL.Query().Get("type") != "list" {
				t.Fatalf("invalid type parameter")
			}
			if r.URL.Query().Get("v") != "test1234" {
				t.Fatalf("invalid v parameter")
			}

			w.Header().Set("content-Type", "text/xml")
			fmt.Fprintf(w, listXML)
			return
		},
	))
	defer ts.Close()
	youtubets.Domain = ts.URL

	type fields struct {
		ID string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []youtubets.Transcript
		wantErr bool
	}{
		{"#1 normal case",
			fields{"test1234"},
			[]youtubets.Transcript{{VideoID: "test1234", Name: "test_en", Lang: "en", Default: true}, {VideoID: "test1234", Lang: "ja"}},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &youtubets.Video{
				ID: tt.fields.ID,
			}
			got, err := v.TranscriptList()
			if (err != nil) != tt.wantErr {
				t.Errorf("Video.TranscriptList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Video.TranscriptList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTranscriptList(t *testing.T) {
	type args struct {
		videoID   string
		xmlReader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *youtubets.TranscriptList
		wantErr bool
	}{
		{
			"#1 normal case",
			args{"test1234", bytes.NewBuffer([]byte(listXML))},
			&youtubets.TranscriptList{List: []youtubets.Transcript{{VideoID: "test1234", Name: "test_en", Lang: "en", Default: true}, {VideoID: "test1234", Lang: "ja"}}},
			false,
		},
		{
			"#2 empty xml error case",
			args{"test11111", bytes.NewBuffer([]byte(""))},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := youtubets.NewTranscriptList(tt.args.videoID, tt.args.xmlReader)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTranscriptList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTranscriptList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTranscript_Text(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/timedtext" {
				t.Fatalf("invalid List API Path %s", r.URL.Path)
			}
			if r.URL.Query().Get("lang") != "en" {
				t.Fatalf("invalid type parameter")
			}
			if r.URL.Query().Get("name") != "test_en" {
				t.Fatalf("invalid type parameter")
			}
			if r.URL.Query().Get("v") != "test1234" {
				t.Fatalf("invalid v parameter")
			}

			w.Header().Set("content-Type", "text/xml")
			fmt.Fprintf(w, transcriptXML)
			return
		},
	))
	defer ts.Close()
	youtubets.Domain = ts.URL

	type fields struct {
		VideoID string
		Name    string
		Lang    string
		Default bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    *youtubets.Text
		wantErr bool
	}{
		{"#1 normal case",
			fields{VideoID: "test1234", Name: "test_en", Lang: "en"},
			&youtubets.Text{Lines: []string{"hello", "world.", ">> It's a test."}},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := &youtubets.Transcript{
				VideoID: tt.fields.VideoID,
				Name:    tt.fields.Name,
				Lang:    tt.fields.Lang,
				Default: tt.fields.Default,
			}
			got, err := ts.Text()
			if (err != nil) != tt.wantErr {
				t.Errorf("Transcript.Text() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Transcript.Text() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewText(t *testing.T) {
	type args struct {
		xmlReader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *youtubets.Text
		wantErr bool
	}{
		{
			"#1 normal case",
			args{bytes.NewBuffer([]byte(transcriptXML))},
			&youtubets.Text{Lines: []string{"hello", "world.", ">> It's a test."}},
			false,
		},
		{
			"#2 empty xml error case",
			args{bytes.NewBuffer([]byte(""))},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := youtubets.NewText(tt.args.xmlReader)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewText() = %v, want %v", got, tt.want)
			}
		})
	}
}

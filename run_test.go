package youtubets_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pchatsu/youtubets"
)

func TestCmd_Run(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.String() {
			case "/timedtext?lang=en&name=test_en&v=test0001":
				w.Header().Set("content-Type", "text/xml")
				fmt.Fprintf(w, transcriptXML)
				return
			case "/timedtext?type=list&v=test0001":
				w.Header().Set("content-Type", "text/xml")
				fmt.Fprintf(w, listXML)
				return
			case "/timedtext?type=list&v=test0002":
				w.Header().Set("content-Type", "text/xml")
				fmt.Fprintf(w, emptyListXML)
				return
			default:
				w.Header().Set("content-Type", "text/xml")
				fmt.Fprintf(w, "")
				return
			}
		},
	))
	defer ts.Close()
	youtubets.Domain = ts.URL

	type args struct {
		args []string
		opt  *youtubets.Option
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"#1 normal case", args{[]string{"test0001"}, &youtubets.Option{Lang: "en", Name: "test_en"}}, "hello world.\n>> It's a test.\n", false},
		{"#2 normal case get default lang", args{[]string{"test0001"}, &youtubets.Option{}}, "hello world.\n>> It's a test.\n", false},
		{"#3 normal case get list", args{[]string{"test0001"}, &youtubets.Option{List: true}}, "name \"test_en\" lang \"en\"\nname \"\" lang \"ja\"\n", false},
		{"#4 normal case get empty list", args{[]string{"test0002"}, &youtubets.Option{List: true}}, "has no transcript\n", false},
		{"#5 error case no args", args{[]string{}, &youtubets.Option{List: true}}, "args must have one video id\n", true},
		{"#6 error case empty id", args{[]string{""}, &youtubets.Option{List: true}}, "args must have one video id\n", true},
		{"#7 error case too many args", args{[]string{"test0001", "test0002"}, &youtubets.Option{List: true}}, "args must have one video id\n", true},
		{"#8 error case no transcript", args{[]string{"test0002"}, &youtubets.Option{}}, "has no transcript\n", true},
		{"#9 error case no transcript", args{[]string{"test0001"}, &youtubets.Option{Lang: "ja", Name: "foo"}}, "has empty transcript XML\n", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var inBuf bytes.Buffer
			var outBuf bytes.Buffer
			var errBuf bytes.Buffer
			c := &youtubets.Cmd{
				Stdin:  &inBuf,
				Stdout: &outBuf,
				Stderr: &errBuf,
			}
			if err := c.Run(tt.args.args, tt.args.opt); (err != nil) != tt.wantErr {
				t.Errorf("Cmd.Run() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				if tt.wantErr {
					if errBuf.String() != tt.want {
						t.Errorf("Cmd.Run() Stderr = %v, want %v", errBuf.String(), tt.want)
					}
					if outBuf.String() != "" {
						t.Errorf("Cmd.Run() Stdout = %v, want empty", outBuf.String())
					}
				} else {
					if outBuf.String() != tt.want {
						t.Errorf("Cmd.Run() Stdout = %v, want %v", outBuf.String(), tt.want)
					}
					if errBuf.String() != "" {
						t.Errorf("Cmd.Run() Stderr = %v, want empty", errBuf.String())
					}

				}
			}
		})
	}
}

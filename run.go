package youtubets

import (
	"fmt"
	"html"
	"io"

	"github.com/pkg/errors"
)

type Option struct {
	Lang string
	Name string
	List bool
}

type Cmd struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func (c *Cmd) Run(args []string, opt *Option) error {
	if len(args) != 1 || args[0] == "" {
		err := errors.New("args must have one video id")
		fmt.Fprintln(c.Stderr, err.Error())
		return err
	}
	videoID := args[0]

	if opt.List {
		return c.displayList(videoID)
	}

	if opt.Lang == "" && opt.Name == "" {
		return c.displayDefaultTranscript(videoID)
	}

	ts := Transcript{VideoID: videoID, Name: opt.Name, Lang: opt.Lang}
	return c.displayTransScript(&ts)
}

func (c *Cmd) displayList(videoID string) error {
	v := Video{videoID}
	transcriptList, err := v.TranscriptList()
	if err != nil {
		fmt.Fprintln(c.Stderr, err.Error())
		return err
	}

	if len(transcriptList) <= 0 {
		fmt.Fprintln(c.Stdout, "has no transcript") // is not error
		return nil
	}

	for _, transcript := range transcriptList {
		fmt.Fprintf(c.Stdout, `name "%s" lang "%s"`+"\n", transcript.Name, transcript.Lang)
	}
	return nil
}

func (c *Cmd) displayTransScript(ts *Transcript) error {
	text, err := ts.Text()
	if err != nil {
		fmt.Fprintln(c.Stderr, err.Error())
		return err
	}

	for _, s := range text.Lines {
		s := html.UnescapeString(s)
		if len(s) == 0 {
			continue
		}

		switch s[len(s)-1:] {
		case ".":
			fmt.Fprintln(c.Stdout, s)
		case " ":
			fmt.Fprint(c.Stdout, s)
		default:

			fmt.Fprint(c.Stdout, s+" ")
		}
	}

	return nil
}

func (c *Cmd) displayDefaultTranscript(videoID string) error {
	v := Video{videoID}
	transcriptList, err := v.TranscriptList()
	if err != nil {
		fmt.Fprintln(c.Stderr, err.Error())
		return err
	}

	if len(transcriptList) <= 0 {
		err := errors.New("has no transcript")
		fmt.Fprintln(c.Stderr, err.Error())
		return err
	}

	for _, ts := range transcriptList {
		if ts.Default {
			return c.displayTransScript(&ts)
		}
	}

	return c.displayTransScript(&transcriptList[0])
}

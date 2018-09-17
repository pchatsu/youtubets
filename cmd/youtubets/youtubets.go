package main

import (
	"flag"
	"os"

	"github.com/pchatsu/youtubets"
)

var (
	Version  string
	Revision string
)

var (
	lang = flag.String("lang", "en", "target language")
	name = flag.String("name", "", "target track name")
	list = flag.Bool("l", false, "display enable transcripts")
)

func main() {
	flag.Parse()
	if err := Run(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

func Run() error {
	cmd := youtubets.Cmd{Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stdout}
	opt := youtubets.Option{Lang: *lang, Name: *name, List: *list}
	if err := cmd.Run(flag.Args(), &opt); err != nil {
		return err
	}
	return nil
}

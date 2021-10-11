package tmpl

import (
	"errors"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/term"
)

type ContextType string

const (
	ENV  ContextType = "env"
	JSON ContextType = "json"
)

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }

func getEnvironContext() (context map[string]string) {
	context = map[string]string{}
	for _, line := range os.Environ() {
		parts := strings.SplitN(line, "=", 2)
		context[parts[0]] = parts[1]
	}
	return
}

func getContext(s string, t ContextType) (interface{}, error) {
	if s == "-" {
		return getEnvironContext(), nil
	}
	return nil, errors.New("invalid context")
}

func getTemplateReader(s string) (io.ReadCloser, error) {
	if s == "-" {
		if !term.IsTerminal(int(os.Stdin.Fd())) {
			return io.NopCloser(os.Stdin), nil
		} else {
			return nil, errors.New("no stdin")
		}
	}

	return os.Open(s)
}

func getResultWriter(s string) (io.WriteCloser, error) {
	if s == "-" {
		return nopCloser{os.Stdout}, nil
	}

	return os.Create(s)
}

func Execute(dst io.Writer, src io.Reader, context interface{}) error {
	tp := template.New("tmpl")

	raw, err := ioutil.ReadAll(src)
	if err != nil {
		return err
	}
	tp.Parse(string(raw))

	if err := tp.Execute(dst, context); err != nil {
		return err
	}

	return nil
}

type Template struct {
	Source      string
	Destination string
	Context     string
	ContextType ContextType
}

func (t Template) Execute() error {
	context, err := getContext(t.Context, t.ContextType)
	if err != nil {
		return err
	}
	src, err := getTemplateReader(t.Source)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := getResultWriter(t.Destination)
	if err != nil {
		return err
	}
	defer dst.Close()
	if err := Execute(dst, src, context); err != nil {
		return err
	}

	return nil
}

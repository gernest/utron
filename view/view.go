package view

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// View is an interface for rendering templates.
type View interface {
	Render(out io.Writer, name string, data interface{}) error
}

// SimpleView implements View interface, but based on golang templates.
type SimpleView struct {
	viewDir string
	tmpl    *template.Template
}

//NewSimpleView returns a SimpleView with templates loaded from viewDir
func NewSimpleView(viewDir string) (View, error) {
	info, err := os.Stat(viewDir)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("utron: %s is not a directory", viewDir)
	}
	s := &SimpleView{
		viewDir: viewDir,
		tmpl:    template.New(filepath.Base(viewDir)),
	}
	return s.load(viewDir)
}

// load loads templates from dir. The templates should be valid golang templates
//
// Only files with extension .html, .tpl, .tmpl will be loaded. references to these templates
// should be relative to the dir. That is, if  dir is foo, you don't have to refer to
// foo/bar.tpl, instead just use bar.tpl
func (s *SimpleView) load(dir string) (View, error) {

	// supported is the list of file extensions that will be parsed as templates
	supported := map[string]bool{".tpl": true, ".html": true, ".tmpl": true}

	werr := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		extension := filepath.Ext(path)
		if _, ok := supported[extension]; !ok {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		// We remove the directory name from the path
		// this means if we have directory foo, with file bar.tpl
		// full path for bar file foo/bar.tpl
		// we trim the foo part and remain with /bar.tpl
		//
		// NOTE we don't account for the opening slash, when dir ends with /.
		name := path[len(dir):]

		name = filepath.ToSlash(name)

		name = strings.TrimPrefix(name, "/") // case  we missed the opening slash

		name = strings.TrimSuffix(name, extension) // remove extension

		t := s.tmpl.New(name)
		
		if _, err = t.Parse(string(data)); err != nil {
			return err
		}
		return nil
	})

	if werr != nil {
		return nil, werr
	}

	return s, nil
}

// Render executes template named name, passing data as context, the output is written to out.
func (s *SimpleView) Render(out io.Writer, name string, data interface{}) error {
	return s.tmpl.ExecuteTemplate(out, name, data)
}

package lib

import (
	"errors"
	"fmt"
	"github.com/xuzhuoxi/infra-go/osxu"
	"io/ioutil"
	"text/template"
)

var templateMap map[string]*Template

func init() {
	templateMap = make(map[string]*Template)
}

type Template struct {
	Name string
	temp *template.Template
}

func (temp *Template) CloneTemplate() *template.Template {
	if nil != temp.temp {
		return nil
	}
	clone, _ := temp.temp.Clone()
	return clone
}

func LoadTemplate(TempFile string) (*Template, error) {
	if temp, ok := templateMap[TempFile]; ok {
		return temp, nil
	}
	if !osxu.IsExist(TempFile) {
		return nil, errors.New(fmt.Sprintf("Templete File Not Found: \"%s\"", TempFile))
	}
	body, err := ioutil.ReadFile(TempFile)
	if nil != err {
		return nil, err
	}
	text := string(body)
	temp, err := template.New(TempFile).Parse(text)
	if nil != err {
		return nil, err
	}
	rs := &Template{Name: TempFile, temp: temp}
	templateMap[TempFile] = rs
	return rs, nil
}

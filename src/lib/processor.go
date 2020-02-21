package lib

import (
	"bytes"
	"io/ioutil"
	"os"
)

func DoProcess(proxy *ExcelProxy, tempPath string, target string) error {
	temp, err := LoadTemplate(tempPath)
	if nil != err {
		return err
	}
	var bs bytes.Buffer
	err = temp.Template.Execute(&bs, proxy)
	if nil != err {
		return err
	}
	//fmt.Println("啦啦啦", string(bs.Bytes()))
	err = ioutil.WriteFile(target, bs.Bytes(), os.ModePerm)
	if nil != err {
		return err
	}
	return nil
}

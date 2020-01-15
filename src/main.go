package main

import (
	"bytes"
	"fmt"
	"github.com/xuzhuoxi/ConstantTableProcessor/src/lib"
	"github.com/xuzhuoxi/infra-go/logx"
	"io/ioutil"
	"os"
)

func main() {
	//初始化Logger
	logger := logx.NewLogger()
	logger.SetConfig(logx.LogConfig{Type: logx.TypeConsole, Level: logx.LevelAll})

	var err error
	//处理启动参数
	fg, err := lib.ParseFlag()
	if nil != err {
		logger.Errorln(err)
		return
	}

	//加载配置文件
	config, err := lib.LoadConfig(fg.ConfigFile)
	if nil != err {
		logger.Errorln(err)
		return
	}
	config.MakeDetailed(fg.BasePath)
	fmt.Println(config)

	//根据配置进行处理
	for _, processor := range config.Processor {
		excelProxy := &lib.ExcelProxy{}
		err = excelProxy.LoadSheets(processor.Source, processor.SheetPrefix, processor.NickRow)
		if nil != err {
			logger.Errorln(err)
			continue
		}
		err = excelProxy.MergedRowsByFilter(processor.StartRow, func(row *lib.ExcelRow) bool {
			return !row.Empty()
		})
		if nil != err {
			logger.Errorln(err)
			continue
		}
		for _, process := range processor.Process {
			temp, err := lib.LoadTemplate(process.Temp)
			if nil != err {
				logger.Errorln(err)
				break
			}
			out := bytes.NewBuffer(nil)
			err = temp.Template.Execute(out, excelProxy)
			if nil != err {
				logger.Errorln(err)
				break
			}
			ioutil.WriteFile(process.Target, out.Bytes(), os.ModePerm)
		}
	}
}

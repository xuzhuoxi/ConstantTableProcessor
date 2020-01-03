package main

import (
	"github.com/xuzhuoxi/ConstantTableProcessor/src/lib"
	"github.com/xuzhuoxi/infra-go/logx"
)

func main() {
	logger := logx.NewLogger()
	logger.SetConfig(logx.LogConfig{Type: logx.TypeConsole, Level: logx.LevelAll})
	fg, err := lib.ParseFlag()
	if nil != err {
		logger.Errorln(err)
		return
	}
	config, err := lib.LoadConfig(fg.ConfigFile)
	if nil != err {
		logger.Errorln(err)
		return
	}
	config.MakeDetailed(fg.BasePath)
	runAt(config)
}

func runAt(config *lib.Config) {

}

package lib

import (
	"flag"
	"errors"
	"github.com/xuzhuoxi/infra-go/osxu"
	"strings"
	"fmt"
	"io/ioutil"
	"github.com/json-iterator/go"
)

type Process struct {
	Temp   string //使用模板文件名
	Target string //输出文件名
}

type Processor struct {
	Source      string //来源文件路径，支持文件和文件夹，文件夹不支持递归子文件夹。
	SheetPrefix string //参与Sheet的名称前缀，可选。若没有则使用SheetPrefix值。
	StartRow    int    //开始行号，可选。若没有则使用StartRow值。
	Process     []Process
}

type Config struct {
	TempFolder   string //模板目录，可选。若无，Temp值要求使用绝对路径。
	SourceFolder string //来源目录，可选。若无，Source值要求使用绝对路径。
	TargetFolder string //目标目录，可选。若无，Target值要求使用绝对路径。
	SheetPrefix  string //参与Sheet的名称前缀，可选。若无，每个Processor必须配置SheetPrefix值。
	StartRow     int    //开始行号，可选。若无，每个Processor必须配置StartRow值。
	Processor    [] Processor
}

func (c *Config) MakeDetailed(BasePath string) error {
	if len(c.Processor) == 0 {
		return errors.New("No Processor! ")
	}
	if c.TempFolder != "" {
		c.TempFolder = osxu.FormatDirPath(BasePath + c.TempFolder)
	}
	if c.SourceFolder != "" {
		c.SourceFolder = osxu.FormatDirPath(BasePath + c.SourceFolder)
	}
	if c.TargetFolder != "" {
		c.TargetFolder = osxu.FormatDirPath(BasePath + c.TargetFolder)
	}
	for index := range c.Processor {
		if c.Processor[index].SheetPrefix == "" && c.SheetPrefix == "" {
			return errors.New("SheetPrefix Undefined! ")
		}
		if c.Processor[index].SheetPrefix == "" {
			c.Processor[index].SheetPrefix = c.SheetPrefix
		}

		if c.Processor[index].StartRow == 0 && c.StartRow == 0 {
			return errors.New("StartRow Undefined! ")
		}
		if c.Processor[index].StartRow == 0 {
			c.Processor[index].StartRow = c.StartRow
		}
		if c.StartRow < 1 {
			return errors.New("StartRow Should be greater than 0! ")
		}

		if len(c.Processor[index].Process) == 0 {
			return errors.New("No Process! ")
		}
		for index2 := range c.Processor[index].Process {
			c.Processor[index].Process[index2].Temp = c.TempFolder + c.Processor[index].Process[index2].Temp
			c.Processor[index].Process[index2].Target = c.TargetFolder + c.Processor[index].Process[index2].Target
		}
	}
	return nil
}

type Flag struct {
	BasePath   string
	ConfigFile string
}

// -base 	可选	自定义基目录	字符串路径，文件夹或文件,"./"开头视为相对路径
// -config 	配置文件，默认为config.json
func ParseFlag() (fg *Flag, err error) {
	base := flag.String("base", "", "Input Base Path! ")
	config := flag.String("config", "config.json", "Config File! ")
	flag.Parse()

	BasePath := *base
	if "" == BasePath || "." == BasePath || "./" == BasePath {
		BasePath = osxu.RunningBaseDir()
	} else if strings.Index(BasePath, "./") == 0 {
		BasePath = osxu.RunningBaseDir() + BasePath
	}
	if nil == config || "" == *config {
		return nil, errors.New("Config Not Define! ")
	}
	ConfigFile := BasePath + *config
	return &Flag{BasePath: BasePath, ConfigFile: ConfigFile}, nil
}

func LoadConfig(configFile string) (config *Config, err error) {
	if !osxu.IsExist(configFile) {
		return nil, errors.New(fmt.Sprintf("Config \"%s\" is not exist!", configFile))
	}
	cfgBody, err := ioutil.ReadFile(configFile)
	if nil != err {
		return nil, errors.New(fmt.Sprintf("Config \"%s\" is not exist!", configFile))
	}
	config = &Config{}
	err = jsoniter.Unmarshal(cfgBody, config)
	if nil != err {
		return nil, err
	}
	return config, nil
}

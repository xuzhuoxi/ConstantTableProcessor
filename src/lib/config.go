package lib

import (
	"errors"
	"flag"
	"fmt"
	"github.com/json-iterator/go"
	"github.com/xuzhuoxi/infra-go/osxu"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Process struct {
	Temp   string //使用模板文件名。支持多个模板路径，可用英文逗号\",\"分割。
	Target string //输出文件名
}

func (p Process) String() string {
	return fmt.Sprintf("Process{Temp=\"%s\", Target=\"%s\"}", p.Temp, p.Target)
}

type Processor struct {
	Source      string //来源文件路径，支持文件和文件夹，文件夹不支持递归子文件夹。支持多个路径，可用英文逗号\",\"分割。
	SheetPrefix string //参与Sheet的名称前缀，可选。若没有则使用SheetPrefix值。
	NickRow     int    //别名行号，可选。若无，NameRow值。
	StartRow    int    //开始行号，可选。若没有则使用StartRow值。
	Process     []Process
}

func (p Processor) String() string {
	return fmt.Sprintf("Processor{Source=\"%s\", SheetPrefix=\"%s\", NickRow=\"%d\", StartRow=\"%d\", Process=\"%s\"}",
		p.Source, p.SheetPrefix, p.NickRow, p.StartRow, fmt.Sprint(p.Process))
}

type Config struct {
	TempFolder   string //模板目录，可选。若无，Temp值等于命令参数-base对应值。
	SourceFolder string //来源目录，可选。若无，Source值等于命令参数-base对应值。
	TargetFolder string //目标目录，可选。若无，Target值等于命令参数-base对应值。
	SheetPrefix  string //参与Sheet的名称前缀，可选。若无，每个Processor必须配置SheetPrefix值。
	NickRow      int    //别名行号，可选。若无，使用Excel列号。
	StartRow     int    //开始行号，可选。若无，每个Processor必须配置StartRow值。
	Processor    []Processor
}

func (c *Config) String() string {
	return fmt.Sprintf("Config{\n TempFolder=%s\n SourceFolder=%s\n TargetFolder=%s\n SheetPrefix=%s\n NickRow=%d\n StartRow=%d\n Processor=%s\n}",
		c.TempFolder, c.SourceFolder, c.TargetFolder, c.SheetPrefix, c.NickRow, c.StartRow, fmt.Sprint(c.Processor))
}

// 对原始配置数据进行详细化处理
// 1.补充默认参数
// 2.详细化路径
func (c *Config) MakeDetailed(BasePath string) error {
	if len(c.Processor) == 0 {
		return errors.New("No Processor! ")
	}

	c.TempFolder = linkPath(c.TempFolder, BasePath)
	c.SourceFolder = linkPath(c.SourceFolder, BasePath)
	c.TargetFolder = linkPath(c.TargetFolder, BasePath)

	for index := range c.Processor {
		// Source处理
		c.Processor[index].Source = linkPaths(c.Processor[index].Source, c.SourceFolder)

		// SheetPrefix处理
		if c.Processor[index].SheetPrefix == "" && c.SheetPrefix == "" {
			return errors.New("SheetPrefix Undefined! ")
		}
		if c.Processor[index].SheetPrefix == "" {
			c.Processor[index].SheetPrefix = c.SheetPrefix
		}

		// NickRow处理
		if c.Processor[index].NickRow == 0 {
			c.Processor[index].NickRow = c.NickRow
		}
		if c.NickRow < 0 {
			return errors.New("NickRow Should be greater than 0! ")
		}

		// StartRow处理
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
			c.Processor[index].Process[index2].Temp = linkPaths(c.Processor[index].Process[index2].Temp, c.TempFolder)
			c.Processor[index].Process[index2].Target = linkPath(c.Processor[index].Process[index2].Target, c.TargetFolder)
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
	base := flag.String("base", "./", "Input Base Path! ")
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

// 加载配置文件
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

// 补全路径
// 绝对路径忽略
func linkPaths(paths string, baseFolder string) string {
	pathArr := strings.Split(paths, ",")
	rs := ""
	for index, p := range pathArr {
		linkedPath := linkPath(p, baseFolder)
		if index == len(pathArr)-1 {
			rs = rs + linkedPath
		} else {
			rs = rs + linkedPath + ","
		}
	}
	return rs
}

func linkPath(path string, baseFolder string) string {
	baseFolder = osxu.FormatDirPath(baseFolder)
	if "" == path || "." == path || "./" == path {
		return baseFolder
	}
	if filepath.IsAbs(path) {
		return osxu.FormatPath(path)
	}
	return osxu.FormatPath(baseFolder + path)
}

# ConstantTableProcessor
ConstantTableProcessor 可以根据Excel表格内容进行模板导出。通常可用于代码生成、配置生成、数据文件生成等。

[English](/README.md)/中文

## 兼容性
go 1.11

## 依赖性

- infra-go(库依赖) [https://github.com/xuzhuoxi/infra-go](https://github.com/xuzhuoxi/infra-go)

- goxc(编译依赖) [https://github.com/laher/goxc](https://github.com/laher/goxc) 

## 使用方法

### 下载

根据系统平台，下载对应的执行文件。[这里下载](/download.md)

### 配置

#### 配置文件

配置文件格式为json文件，可参考[config.json](/test/config.json)

- "TempFolder": 

	模板目录，可选。若无，Temp值等于命令参数-base对应值。
	
- "SourceFolder": 
	来源目录，可选。若无，Source值等于命令参数-base对应值。
- "TargetFolder": 
	目标目录，可选。若无，Target值等于命令参数-base对应值。
- "SheetPrefix": 
	参与Sheet的名称前缀，可选。若无，每个Processor必须配置SheetPrefix值。
- "NickRow": 
	别名行号，可选。若无，使用Excel列号。
- "StartRow": 
	开始行号，可选。若无，每个Processor必须配置StartRow值。
- "Processor.Source": 
	来源文件路径，支持文件和文件夹，文件夹不支持递归子文件夹。支持多个路径，可用英文逗号","分割。
- "Processor.SheetPrefix": 
	参与Sheet的名称前缀，可选。若没有则使用SheetPrefix值。
- "Processor.NickRow": 
	别名行号，可选。若无，NameRow值。
- "Processor.StartRow": 
	开始行号，可选。若没有则使用StartRow值。
- "Processor.Process.Temp": 
	使用模板文件名。支持多个模板路径，可用英文逗号","分割。
- "Processor.Process.Target": 
	输出文件路径。
	
**注意**

- 数据源只支持Excel文件，包括xlsx文件和xls文件。
- Processor和Process为数组，执行的处理数为Processor的数量乘以Process的数量。

#### 模板文件

模板文件格式为go语言模板，文档说明如下:

[https://golang.google.cn/pkg/text/template/](https://golang.google.cn/pkg/text/template/)

模板接收的数据结构体为[lib.ExcelProxy](/src/lib/excel.go)

开放属性及行为如下：

ExcelProxy:
```
type ExcelProxy struct {
	Sheets   []*ExcelSheet
	DataRows []*ExcelRow
}

func (ep *ExcelProxy) GetSheet(sheet string) (es *ExcelSheet, err error)

// Open to templates
// 通过Sheet的名称和坐标取值，坐标格式：B4
func (ep *ExcelProxy) ValueAtAxis(sheet string, axis string) (value string, err error)

```

ExcelSheet:
```
type ExcelSheet struct {
	SheetIndex int
	SheetName  string

	Axis []string
	Nick []string
	Rows []*ExcelRow

	RowLength int //行数
	ColLength int //列数
}

func (es *ExcelSheet) GetDataRows(startIndex int) (rows []*ExcelRow)

// Open to templates
// 通过坐标取值，坐标格式：B4
func (es *ExcelSheet) ValueAtAxis(axis string) (value string, err error)

```

ExcelRow:
```
type ExcelRow struct {
	Index  int
	Length int

	Axis []string
	Nick []string
	Row  []string
}

// Open to templates
// 通过索引号取值，索引号从0开始
func (er *ExcelRow) ValueAtIndex(index int) (value string, err error)

// Open to templates
// 通过别名取值
func (er *ExcelRow) ValueAtNick(nick string) (value string, err error)

// Open to templates
// 通过列名取值
func (er *ExcelRow) ValueAtAxis(axis string) (value string, err error)

```

### 执行

使用命令行执行，如：

`ConstantTableProcessor -base=路径 -config=配置文件路径`

- -base: 
  指定一个基础目录，如果不设置，则使用当前运行目录，即“./”。
  
- -config: 
  指定一个配置文件的路径。建议使用相对路径，实际的绝对路径结合-base生成。
  
**建议**:

写成脚本，方便执行。如：[run.sh](/test/run.sh)

## 联系作者

xuzhuoxi 

<xuzhuoxi@gmail.com> 或 <mailxuzhuoxi@163.com>

## 开源许可证
ConstantTableProcessor 源代码基于[MIT许可证](/LICENSE)进行开源。



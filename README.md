# ExcelExportTool
You can export the template according to the contents of the Excel table. It can be used for code generation, configuration generation, data file generation, etc.

English/[中文](/README_zh.md)

## Compatibility
go 1.11

## Dependence

- infra-go(Library dependencies) [https://github.com/xuzhuoxi/infra-go](https://github.com/xuzhuoxi/infra-go)

- excelize(Library dependencies) [https://github.com/360EntSecGroup-Skylar/excelize](https://github.com/360EntSecGroup-Skylar/excelize)

- goxc(Compile dependencies) [https://github.com/laher/goxc](https://github.com/laher/goxc) 

## Usage

### Download

Download the corresponding execution file according to the system platform.[Download here](https://github.com/xuzhuoxi/ExcelExportTool/releases)

### Configuration

#### Config File

The configuration file format is json file, you can refer [config.json](/test/config.json).

- "TempFolder": 

	Template directory, optional. If not, the Temp value is equal to the corresponding value of the command parameter -base.
	
- "SourceFolder": 
	
	Source directory, optional. If not, the Source value is equal to the corresponding value of the command parameter -base.
	
- "TargetFolder": 
	
	Destination directory, optional. If not, the Target value is equal to the value of the command parameter -base.
	
- "SheetPrefix": 
	
	Participating sheet name prefix, optional. If not, each Processor must be configured with a SheetPrefix value.
	
- "NickRow": 
	
	Nickname line number, optional. If not, use Excel column number.
	
- "StartRow": 
	
	Start line number, optional. If not, each Processor must be configured with a StartRow value.
	
- "Processor.Source": 
	
	Source file path. Files and folders are supported. Folders do not support recursive subfolders. Support multiple paths, can be separated by comma ",".
	
- "Processor.SheetPrefix": 
	
	Including sheet name prefix, optional. If not, the SheetPrefix value is used.
	
- "Processor.NickRow": 
	
	Nickname line number, optional. If not, the NickRow value is used.
	
- "Processor.StartRow": 
	
	Start line number, optional. If not, the StartRow value is used.
	
- "Processor.Process.Temp": 
	
	Use template file name. Supports multiple template paths, which can be separated by commas ",".
	
- "Processor.Process.Target": 
	
	Output file path.
	
**Note**

- The data source only supports Excel files, including xlsx files and xls files.

- "Processor" and "Process" are arrays, and the number of processes executed is the number of processors multiplied by the number of processes.

#### Template file

The template file format is go language template, the address of the document is as follows:

[https://golang.google.cn/pkg/text/template/](https://golang.google.cn/pkg/text/template/)

The data structure received by the template is: [lib.ExcelProxy](/src/lib/excel.go)

Open attributes and functions are as follows：

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

Use the command line, such as：

`ExcelExportTool -base=Dir -config=filepath`

- -base: 
  
  Specify a base directory. If not set, the current running directory is used, that is, "./".
  
- -config: 
  
  Specify the path of a configuration file. It is recommended to use relative paths, and the actual absolute path combined with -base generation.
  
**Suggestion**:

Write scripts for easy execution. Such as：[run.sh](/test/run.sh)

## Contact

xuzhuoxi 

<xuzhuoxi@gmail.com> or <mailxuzhuoxi@163.com>

## License
ExcelExportTool source code is available under the MIT [License](/LICENSE).



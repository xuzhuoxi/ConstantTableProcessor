package lib

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/xuzhuoxi/infra-go/mathx"
	"github.com/xuzhuoxi/infra-go/osxu"
	"github.com/xuzhuoxi/infra-go/slicex"
	"os"
	"strconv"
	"strings"
)

type ExcelRow struct {
	Index  int
	Length int

	Axis []string
	Nick []string
	Row  []string
}

func (er *ExcelRow) String() string {
	return fmt.Sprintf("ExcelRow{Index=%d, Length=%d, Axis=%s, Nick=%s, Row=%s}", er.Index, er.Length,
		fmt.Sprint(er.Axis), fmt.Sprint(er.Nick), fmt.Sprint(er.Row))
}

// Open to templates
// 通过索引号取值，索引号从0开始
func (er *ExcelRow) ValueAtIndex(index int) (value string, err error) {
	if index < 0 || index >= er.Length {
		return "", errors.New(fmt.Sprintf("Index(%d) out of range! ", index))
	}
	return er.Row[index], nil
}

// Open to templates
// 通过别名取值
func (er *ExcelRow) ValueAtNick(nick string) (value string, err error) {
	index, ok := slicex.IndexString(er.Nick, nick)
	if !ok {
		return "", errors.New(fmt.Sprintf("Nick(%s) is not exist! ", nick))
	}
	return er.Row[index], nil
}

// Open to templates
// 通过列名取值
func (er *ExcelRow) ValueAtAxis(axis string) (value string, err error) {
	index, ok := slicex.IndexString(er.Axis, axis)
	if !ok {
		return "", errors.New(fmt.Sprintf("Axis(%s) is not exist! ", axis))
	}
	return er.Row[index], nil
}

//-----------------

type ExcelSheet struct {
	SheetIndex int
	SheetName  string

	Axis []string
	Nick []string
	Rows []*ExcelRow

	RowLength int //行数
	ColLength int //列数
}

func (es *ExcelSheet) String() string {
	strRows := ""
	for _, r := range es.Rows {
		strRows = strRows + "\t" + fmt.Sprint(r) + "\n"
	}
	return fmt.Sprintf("ExcelSheet{Index=%d, Name=%s, Axis=%s, Nick=%s, RowLen=%d, ColLen=%d,\nRow=\n%s}",
		es.SheetIndex, es.SheetName, fmt.Sprint(es.Axis), fmt.Sprint(es.Nick), es.RowLength, es.ColLength, strRows)
}

func (es *ExcelSheet) SetNick(nick []string) {
	es.Nick = nick
	for _, r := range es.Rows {
		r.Nick = nick
	}
}

// Open to templates
// 通过坐标取值，坐标格式：B4
func (es *ExcelSheet) ValueAtAxis(axis string) (value string, err error) {
	colIndex, rowIndex, err := ParseAxis(axis)
	if nil != err {
		return "", err
	}
	return es.Rows[rowIndex].Row[colIndex], nil
}

//-----------------

type ExcelProxy struct {
	Sheets []*ExcelSheet
}

func (ep *ExcelProxy) GetSheet(sheet string) (es *ExcelSheet, err error) {
	for _, s := range ep.Sheets {
		if s.SheetName == sheet {
			return s, nil
		}
	}
	return nil, errors.New("No Sheet is " + sheet)
}

// 加载SourcePath指定的一个或多个Excel文件。
// SourcePath支持多路径模式，用","分隔。
// SourcePath支持文件夹，不支持递归。
func (ep *ExcelProxy) LoadSheets(SourcePath string, SheetPrefix string, NickRow int) error {
	excels, err := LoadExcels(SourcePath)
	if nil != err {
		return err
	}
	for _, excel := range excels {
		sheets, err := LoadSheets(excel, SheetPrefix, NickRow)
		if nil != err {
			return err
		}
		ep.Sheets = append(ep.Sheets, sheets...)
	}
	return nil
}

//-------------------

// 加载路径下的Excel文件，多个路径用","分割
// 支持文件夹路径
func LoadExcels(Path string) (excels []*excelize.File, err error) {
	paths := strings.Split(strings.TrimSpace(Path), ",")
	if len(paths) == 0 {
		return nil, errors.New("Path Empty:" + Path)
	}
	var filePaths []string
	for _, path := range paths {
		fp := osxu.FormatPath(path)
		if osxu.IsFolder(fp) {
			files, err := osxu.GetFolderFileList(fp, false, func(fileInfo os.FileInfo) bool {
				if fileInfo.IsDir() {
					return false
				}
				name := fileInfo.Name()
				if osxu.CheckExtensionName(name, ".xls") || osxu.CheckExtensionName(name, ".xlsx") {
					return true
				}
				return false
			})
			if nil != err {
				continue
			}
			for _, fi := range files {
				filePaths = append(filePaths, fi.FullPath())
			}
		} else {
			if osxu.CheckExtensionName(fp, ".xls") || osxu.CheckExtensionName(fp, ".xlsx") {
				filePaths = append(filePaths, fp)
			}
		}
	}
	if len(filePaths) == 0 {
		return nil, errors.New("Path Empty:" + Path)
	}
	for _, filePath := range filePaths {
		excel, err := LoadExcel(filePath)
		if err != nil {
			return nil, err
		}
		excels = append(excels, excel)
	}
	return excels, nil
}

// 加载Excel文件，过滤无Sheet情况
func LoadExcel(FileName string) (excel *excelize.File, err error) {
	excelFile, err := excelize.OpenFile(FileName)
	if nil != err {
		return nil, err
	}
	if excelFile.SheetCount <= 0 {
		return nil, errors.New("No Sheets! ")
	}
	return excelFile, nil
}

// 通过SheetPrefix作为限制加载Sheet
// 指定NickRow所在行为别名,NickRow=0时，使用列号作为别名
func LoadSheets(excelFile *excelize.File, SheetPrefix string, NickRow int) (sheets []*ExcelSheet, err error) {
	var names []string
	var indexs []int
	for index, name := range excelFile.GetSheetMap() {
		if strings.Index(name, SheetPrefix) == 0 {
			names = append(names, name)
			indexs = append(indexs, index)
		}
	}
	if len(names) == 0 {
		return nil, nil
	}
	for i, n := range names {
		rows, err := excelFile.GetRows(n)
		if nil != err {
			return nil, err
		}
		es := &ExcelSheet{SheetIndex: indexs[i], SheetName: n}
		es.RowLength = len(rows)
		if es.RowLength > 0 {
			es.ColLength = len(rows[0])
			es.Axis = GenAxis(es.ColLength)
			for rowIndex, row := range rows {
				er := &ExcelRow{Index: rowIndex, Length: es.ColLength, Axis: es.Axis, Row: row}
				es.Rows = append(es.Rows, er)
			}
		}
		if NickRow > 0 {
			es.SetNick(rows[NickRow-1])
		} else {
			es.SetNick(es.Axis)
		}
		sheets = append(sheets, es)
	}
	return sheets, nil
}

// 3 => [A, B, C]
func GenAxis(length int) []string {
	rs := make([]string, length, length)
	for index := 0; index < length; index += 1 {
		rs[index] = mathx.System10To26(index + 1)
	}
	return rs
}

// "A1" => 0, 0, nil
func ParseAxis(axis string) (colIndex int, rowIndex int, err error) {
	Axis := strings.ToUpper(strings.TrimSpace(axis))
	bs := []byte(Axis)
	var c, r []byte
	for index, b := range bs {
		if !(b >= 'A' && b <= 'Z') {
			c, r = bs[:index], bs[index:]
		}
	}
	if nil == c && nil == r {
		return 0, 0, errors.New("Axis Error:" + axis)
	}
	colNum := mathx.System26To10(string(c))
	rowNum, err := strconv.Atoi(string(r))
	if nil != err {
		return 0, 0, err
	}
	return colNum - 1, rowNum - 1, nil
}

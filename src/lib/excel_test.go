package lib

import (
	"fmt"
	"github.com/xuzhuoxi/infra-go/osxu"
	"testing"
)

func TestLoadExcel(t *testing.T) {
	path := osxu.GetRunningDir() + "/Source/const.xlsx"
	excelProxy := &ExcelProxy{}
	err := excelProxy.LoadSheets(path, "Const_", 0)
	if nil != err {
		fmt.Println(err)
		return
	}

	err = excelProxy.MergedRowsByFilter(2, func(row *ExcelRow) bool {
		return !row.Empty()
	})
	//err = excelProxy.MergedRows(2)
	if nil != err {
		fmt.Println(err)
		return
	}

	fmt.Println(excelProxy.Sheets)
	fmt.Println(excelProxy.DataRows, len(excelProxy.DataRows))
}

//func TestLoadSheets(t *testing.T) {
//	path := osxu.RunningBaseDir() + "Source/const.xlsx"
//	excel, err := LoadExcel(path)
//	if nil != err {
//		fmt.Println(err)
//		return
//	}
//	excel.LoadSheets("Const_", 2)
//	for _, s := range excel.Sheets {
//		fmt.Println(*s)
//	}
//}

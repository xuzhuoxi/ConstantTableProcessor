package lib

import (
	"fmt"
	"github.com/xuzhuoxi/infra-go/osxu"
	"testing"
)

func TestLoadExcel(t *testing.T) {
	path := osxu.RunningBaseDir() + "Source/const.xlsx"
	excel, err := LoadExcel(path)
	if nil != err {
		fmt.Println(err)
		return
	}
	excel.LoadSheets("Count_", 0)
	rows, err := excel.ExcelFile.GetRows("Const_A")
	if nil != err {
		fmt.Println(err)
		return
	}
	ln := len(rows)
	index := 0
	for index < ln {
		for _, v := range rows[index] {
			fmt.Print(v, ",")
		}
		fmt.Println()
		index += 1
	}
}

func TestLoadSheets(t *testing.T) {
	path := osxu.RunningBaseDir() + "Source/const.xlsx"
	excel, err := LoadExcel(path)
	if nil != err {
		fmt.Println(err)
		return
	}
	excel.LoadSheets("Const_", 2)
	for _, s := range excel.Sheets {
		fmt.Println(*s)
	}
}

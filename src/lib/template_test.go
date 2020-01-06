package lib

import (
	"os"
	"testing"
)

type ab struct {
	A string
	B string
	C string
}

func TestLoadTemplate(t *testing.T) {
	path := "/Users/xuzhuoxi/go/src/github.com/xuzhuoxi/ConstantTableProcessor/test/Temp/const.tmp"
	temp, _ := LoadTemplate(path)

	datas := []ab{{A: "A1", B: "B1", C: "B1"}, {A: "A2", B: "B2"}, {A: "A3", B: "B3"}, {A: "A4", B: "B4"}}
	//a:={B:"B1"}
	//a["A"] = "A1"
	//datas2 := make(map[string]string)
	//datas2["A1"] = "B1"
	//datas2["A2"] = "B2"
	//datas2["A3"] = "B3"
	//datas2["A4"] = "B4"

	temp.temp.Execute(os.Stdout, datas)
}

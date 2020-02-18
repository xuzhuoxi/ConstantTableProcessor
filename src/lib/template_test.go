package lib

import (
	"fmt"
	"github.com/xuzhuoxi/infra-go/osxu"
	"os"
	"testing"
)

type ab struct {
	A string
	B string
	C string
}

func TestLoadTemplate(t *testing.T) {
	path := osxu.RunningBaseDir() + "Temp/const.tmp"
	temp, _ := LoadTemplate(path)

	datas := []ab{{A: "A1", B: "B1", C: "B1"}, {A: "A2", B: "B2"}, {A: "A3", B: "B3"}, {A: "A4", B: "B4"}}
	//a:={B:"B1"}
	//a["A"] = "A1"
	//datas2 := make(map[string]string)
	//datas2["A1"] = "B1"
	//datas2["A2"] = "B2"
	//datas2["A3"] = "B3"
	//datas2["A4"] = "B4"

	temp.Execute(os.Stdout, datas)
}

func TestLoadTemplates(t *testing.T) {
	baseDir := osxu.RunningBaseDir()
	path1 := baseDir + "Temp/const.tmp"
	path2 := baseDir + "Temp/const2.tmp"
	temp, _ := LoadTemplates([]string{path1, path2})

	datas := []ab{{A: "A1", B: "B1", C: "B1"}, {A: "A2", B: "B2"}, {A: "A3", B: "B3"}, {A: "A4", B: "B4"}}
	//a:={B:"B1"}
	//a["A"] = "A1"
	//datas2 := make(map[string]string)
	//datas2["A1"] = "B1"
	//datas2["A2"] = "B2"
	//datas2["A3"] = "B3"
	//datas2["A4"] = "B4"
	fmt.Println(temp, temp.Template)
	temp.ExecuteTemplate(os.Stdout, datas)
}

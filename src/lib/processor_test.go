package lib

import (
	"fmt"
	"github.com/xuzhuoxi/infra-go/osxu"
	"os"
	"testing"
	"text/template"
)

var tempPath = osxu.RunningBaseDir() + "test/Temp/const.tmp"

//func TestProcessor(t *testing.T) {
//	temp, _ := LoadTemplate(tempPath)
//
//	//excel, err := LoadExcel(sourcePath)
//	//if nil != err {
//	//	fmt.Println(err)
//	//	return
//	//}
//	//excel.LoadSheets("Const_", 0)
//	//fmt.Println(excel.Sheets)
//	//Template.Template.Execute(os.Stdout, excel.Sheets[0].Rows)
//
//	mapData := make(map[string]string)
//	mapData["A"] = "A1"
//	mapData["B"] = "B1"
//	mapData["C"] = "C1"
//
//	mapData2 := make([]string, 3, 3)
//	mapData2[0] = "A1"
//	mapData2[1] = "B1"
//	mapData2[2] = "C1"
//
//	fmt.Println(len(mapData), len(mapData2))
//
//	temp.Template.Execute(os.Stdout, mapData)
//}

type Person struct {
	Name     string
	Children []string
}

func (p *Person) GetChildrenCount() int {
	return len(p.Children)
}

func (p *Person) GetChild(index int) string {
	return p.Children[index]
}

func TestTemplateArray(t *testing.T) {
	p := []*Person{
		{
			Name:     "Susy",
			Children: []string{"Bob", "Herman", "Sherman"},
		},
		{
			Name:     "Norman",
			Children: []string{"Rachel", "Ross", "Chandler"},
		},
	}

	str := `{{$people := .}}
			{{range $i, $pp := $people}}
    			Name: {{$pp.Name}} ChildrenCount: {{$pp.GetChildrenCount}}
    			Children:
      				{{range $j, $c := $pp.Children}}
					Child {{$j}}: {{$c}}
      				{{end}}       
				Index: {{$pp.GetChild 0}}
			{{end}}`

	tp := template.Must(template.New("abc").Parse(str))
	err := tp.Execute(os.Stdout, p)
	if err != nil {
		fmt.Println(err)
	}
}

func TestTemplateMap(t *testing.T) {
	pm := make(map[string]*Person)
	pm["Susy"] = &Person{Name: "Susy", Children: []string{"Bob", "Herman", "Sherman"}}
	pm["Norman"] = &Person{Name: "Norman", Children: []string{"Rachel", "Ross", "Chandler"}}

	str := `{{$people := .}}
			{{range $i, $pp := $people}}
    			Name: {{$pp.Name}}
    			Children:
      				{{range $j, $c := $pp.Children}}
					Child {{$j}}: {{$c}}
      				{{end}}                   
			{{end}}`

	tp := template.Must(template.New("abc").Parse(str))
	err := tp.Execute(os.Stdout, pm)
	if err != nil {
		fmt.Println(err)
	}
}

func TestTemplateMapSimple(t *testing.T) {
	pm := make(map[string]string)
	pm["Susy"] = "Susy1"
	pm["Norman"] = "Norman1"

	str := `{{range $i, $pp := .}}
    			Key: {{$i}}
    			Value: {{$pp}} 
			{{end}}`

	tp := template.Must(template.New("abc").Parse(str))
	err := tp.Execute(os.Stdout, pm)
	if err != nil {
		fmt.Println(err)
	}
}

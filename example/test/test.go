package main

import "fmt"

func main() {
	//aa := `{"Item": [{"WayBillNumber": "YT2121621236013143", "data": }]}`
	//testJsonDom, _ := dom.NewJsonDom(strings.NewReader(aa))
	//itemDataList := testJsonDom.Xpath("//Item")
	//for _, itemData := range itemDataList{
	//	wm := itemData.XpathOne("//WayBillNumber").InnerText()
	//	fmt.Println(wm)
	//}

	type aaa struct {
		bbb bool
	}

	ccc := aaa{}
	fmt.Println(ccc)
}
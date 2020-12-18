package AppleProductMonitor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestCrawler_parse(t *testing.T) {
	c := &Crawler{}

	b, err := ioutil.ReadFile("test_resources/sample.html")
	if err != nil {
		t.Fatalf("err : %s", err)
	}
	source := ioutil.NopCloser(bytes.NewReader(b))
	p, _ := c.parse(source)

	j, err := json.MarshalIndent(p, "", "\t")
	if err != nil {
		t.Fatalf("err : %s", err)
	}
	fmt.Println(string(j))
}

//
//func TestCrawler_compare(t *testing.T) {
//	c := &Crawler{}
//
//	a := `{"arr":[
//	{
//		"Group": "iPhone",
//		"Product": "iPhone 12 Pro Max",
//		"Model": "A2411",
//		"NCC": "CCAI204G0440T0"
//	},
//	{
//		"Group": "iPhone",
//		"Product": "iPhone 12 mini",
//		"Model": "A2399",
//		"NCC": "CCAI204G0430T0"
//	},
//	{
//		"Group": "iPhone",
//		"Product": "iPhone 12 Pro",
//		"Model": "A2407",
//		"NCC": "CCAI204G0420T7"
//	},
//	{
//		"Group": "iPhone",
//		"Product": "iPhone 12",
//		"Model": "A2403",
//		"NCC": "CCAI204G0410T4"
//	}]}`
//
//	b := `{"arr":[
//	{
//		"Group": "iPhone",
//		"Product": "iPhone 12 Pro Max",
//		"Model": "A2411",
//		"NCC": "CCAI204G0440T0"
//	},
//	{
//		"Group": "iPhone",
//		"Product": "iPhone 12 mini",
//		"Model": "A2399",
//		"NCC": "CCAI204G0430T0"
//	},
//	{
//		"Group": "iPhone",
//		"Product": "iPhone 12 Pro Test",
//		"Model": "A2487",ch
//		"NCC": "CCAI204G0421T7"
//	},
//	{
//		"Group": "iPhone",
//		"Product": "iPhone 12",
//		"Model": "A2403",
//		"NCC": "CCAI204G0410T4"
//	}]}`
//
//	//c.compare(a, b)
//}
//
//func TestCrawler_Run(t *testing.T) {
//
//	c := Crawler{}
//	ctx := context.Background()
//	ch := c.Run(ctx)
//
//	for {
//		select {
//		case event := <-ch:
//			t.Logf("%+v", event)
//		}
//	}
//}

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

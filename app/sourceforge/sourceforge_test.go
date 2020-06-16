package sourceforge

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/jdxj/notice/client"
)

type Root struct {
	XMLName xml.Name `xml:"root"`

	Hello  *Hello  `xml:"hello"`
	Hello2 *Hello2 `xml:"hello2"`
}

type Hello struct {
	XMLName xml.Name `xml:"hello"`

	Num   int    `xml:"num,attr"`
	World *World `xml:"world"`
	Pen   *Pen   `xml:"pen"`
}

type Hello2 struct {
	XMLName xml.Name `xml:"hello2"`
	Pen1    string   `xml:",attr"`
	Pen2    string   `xml:",cdata"`
}

type World struct {
	XMLName xml.Name `xml:"world"`
	Apple   string   `xml:"apple,attr"`
	Pen     string   `xml:",innerxml"`
}

type Pen struct {
	Data string `xml:",cdata"`
}

func TestXMLEleCDATA(t *testing.T) {
	h2 := &Hello2{
		Pen1: "apple",
		Pen2: "pen",
	}
	r := &Root{
		Hello2: h2,
	}
	data, err := xml.MarshalIndent(r, "", "  ")
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	fmt.Printf("%s\n", data)
}

func TestMultiCDATA(t *testing.T) {

	h := &Hello2{}

	data, err := xml.MarshalIndent(h, "", "  ")
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	fmt.Printf("%s\n", data)
}

func TestParseMultiCATA(t *testing.T) {
	data := []byte(`<hello2>
    <![CDATA[apple]]>
    <![CDATA[pen]]>
</hello2>`)

	h2 := &Hello2{}
	if err := xml.Unmarshal(data, h2); err != nil {
		t.Fatalf("%s\n", err)
	}

	fmt.Printf("h2: %#v\n", *h2)

	fmt.Printf("pen1: %#v\n", h2.Pen1)
	fmt.Printf("pen2: %#v\n", h2.Pen2)
}

func TestXML(t *testing.T) {
	w := &World{
		Apple: "pen",
		Pen:   "apple",
	}
	h := &Hello{
		Num:   5,
		World: w,
	}
	r := &Root{
		Hello: h,
	}

	data, err := xml.MarshalIndent(r, "", "  ")
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	fmt.Printf("%s\n", data)
}

func TestRSSURL(t *testing.T) {
	c := http.Client{}

	req, err := client.NewRequestUserAgent(http.MethodGet, "", nil)
	if err != nil {
		t.Fatalf("%s\n", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	fmt.Printf("get data ok\n")

	rss := &RSS{}
	err = xml.Unmarshal(data, rss)
	if err != nil {
		t.Fatalf("%s\n", err)
	}

	//fmt.Printf("channel: %#v\n", *rss.Channel)
	//
	//for i, item := range rss.Channel.Items {
	//	fmt.Printf("item-%d: %#v\n", i, *item)
	//}

	data, err = xml.MarshalIndent(rss, "", "  ")
	if err != nil {
		t.Fatalf("%s\n", err)
	}

	fmt.Printf("marshal result:\n%s", data)
}

func TestReadCDATA(t *testing.T) {
	p := &Pen{Data: "fadfi"}

	h := &Hello{
		Num: 0,
		Pen: p,
	}
	data, _ := xml.MarshalIndent(h, "", "  ")
	fmt.Printf("%s\n", data)

	h2 := &Hello{}
	xml.Unmarshal(data, h2)
	fmt.Printf("pen: %s\n", h2.Pen.Data)
}

func TestRaphael_SendUpdate(t *testing.T) {
	r := NewSourceforge("", nil)

	r.UpdateItem()
	r.SendUpdate()
}

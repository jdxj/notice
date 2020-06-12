package raphael

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

	Hello *Hello `xml:"hello"`
}

type Hello struct {
	XMLName xml.Name `xml:"hello"`

	Num   int    `xml:"num,attr"`
	World *World `xml:"world"`
	Pen   *Pen   `xml:"pen"`
}

type World struct {
	XMLName xml.Name `xml:"world"`
	Apple   string   `xml:"apple,attr"`
	Pen     string   `xml:",innerxml"`
}

type Pen struct {
	Data string `xml:",cdata"`
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

	req, err := client.NewRequestUserAgent(http.MethodGet, RSSExRomURL, nil)
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
	r := NewRaphael()

	r.UpdateItem()
	r.SendUpdate()
}

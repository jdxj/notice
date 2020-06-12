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
}

type World struct {
	XMLName xml.Name `xml:"world"`
	Apple   string   `xml:"apple,attr"`
	Pen     string   `xml:",innerxml"`
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

	req, err := client.NewRequestUserAgent(http.MethodGet, RSSURL, nil)
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

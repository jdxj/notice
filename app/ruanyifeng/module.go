package ruanyifeng

import (
	"encoding/xml"
	"fmt"
)

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Xmlns   string   `xml:"xmlns,attr"`
	Title   string   `xml:"title"`

	Links     []*Link    `xml:"link"`
	ID        string     `xml:"id"`
	Updated   string     `xml:"updated"`
	Subtitle  string     `xml:"subtitle"`
	Generator *Generator `xml:"generator"`

	Entries []*Entry `xml:"entry"`
}

type Link struct {
	XMLName xml.Name `xml:"link"`
	Rel     string   `xml:"rel,attr"`
	Type    string   `xml:"type,attr"`
	Href    string   `xml:"href,attr"`
}

type Generator struct {
	XMLName xml.Name `xml:"generator"`
	Uri     string   `xml:"uri,attr"`
	Data    string   `xml:",innerxml"`
}

type Entry struct {
	XMLName xml.Name `xml:"entry"`
	Title   string   `xml:"title"`

	Link      *Link  `xml:"link"`
	ID        string `xml:"id"`
	Published string `xml:"published"`
	Updated   string `xml:"updated"`
	Summary   string `xml:"summary"`

	Author   *Author   `xml:"author"`
	Category *Category `xml:"category"`
	Content  *Content  `xml:"content"`
}

type Author struct {
	XMLName xml.Name `xml:"author"`
	Name    string   `xml:"name"`
	Uri     string   `xml:"uri"`
}

type Category struct {
	XMLName xml.Name `xml:"category"`
	Term    string   `xml:"term,attr"`
	Scheme  string   `xml:"scheme,attr"`
}

type Content struct {
	XMLName xml.Name `xml:"content"`
	Type    string   `xml:"type,attr"`
	Lang    string   `xml:"lang,attr"`
	Base    string   `xml:"base,attr"`

	Data []byte `xml:",cdata"`
}

func unmarshalFeed(data []byte) (*Entry, error) {
	feed := &Feed{}
	if err := xml.Unmarshal(data, feed); err != nil {
		return nil, err
	}
	if feed.Entries == nil || len(feed.Entries) <= 0 {
		return nil, fmt.Errorf("not found entry update")
	}
	return feed.Entries[0], nil
}

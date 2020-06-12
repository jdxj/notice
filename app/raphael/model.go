package raphael

import "encoding/xml"

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Content string   `xml:"content,attr"`
	Files   string   `xml:"files,attr"`
	Media   string   `xml:"media,attr"`
	Doap    string   `xml:"doap,attr"`
	SF      string   `xml:"sf,attr"`
	Version string   `xml:"version,attr"`

	Channel *Channel `xml:"channel"`
}

type Channel struct {
	XMLName xml.Name `xml:"channel"`
	Files   string   `xml:"files,attr"`
	Media   string   `xml:"media,attr"`
	Doap    string   `xml:"doap,attr"`
	SF      string   `xml:"sf,attr"`

	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description *struct {
		Data string `xml:",cdata"`
	} `xml:"description"`
	PubDate        string `xml:"pubDate"`
	ManagingEditor string `xml:"managingEditor"`
	Docs           string `xml:"docs"`

	Items []*Item `xml:"item"`
}

type Item struct {
	XMLName xml.Name `xml:"item"`
	Title   *struct {
		Data string `xml:",cdata"`
	} `xml:"title"`
	Link        string `xml:"link"`
	Guid        string `xml:"guid"`
	PubDate     string `xml:"pubDate"`
	Description *struct {
		Data string `xml:",cdata"`
	} `xml:"description"`

	SfFileId  *SfFileId  `xml:"sf-file-id"`
	ExtraInfo *ExtraInfo `xml:"extra-info"`
	Content   *Content   `xml:"content"`
}

type SfFileId struct {
	XMLName xml.Name `xml:"sf-file-id"`
	Files   string   `xml:"files,attr"`
	Data    string   `xml:",innerxml"`
}

type ExtraInfo struct {
	XMLName xml.Name `xml:"extra-info"`
	Files   string   `xml:"files,attr"`
	Data    string   `xml:",innerxml"`
}

type Content struct {
	XMLName  xml.Name `xml:"content"`
	Media    string   `xml:"media,attr"`
	Type     string   `xml:"type,attr"`
	URL      string   `xml:"url,attr"`
	FileSize int      `xml:"filesize,attr"`

	Hash *Hash `xml:"hash"`
}

type Hash struct {
	XMLName xml.Name `xml:"hash"`
	Algo    string   `xml:"algo,attr"`
	Data    string   `xml:",innerxml"`
}

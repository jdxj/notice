package rss

type entity struct {
	ID  int
	URL string
}

func (entity) TableName() string {
	return "rss_urls"
}

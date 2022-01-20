package crawler

import "github.com/gocolly/colly/v2"

type Interface interface {
	Crawler() colly.HTMLCallback
}

type Visitor interface {
	Visit(u string) error
}

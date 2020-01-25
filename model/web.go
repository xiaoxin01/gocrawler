package model

// Web web to crawl
type Web struct {
	Name         string
	Enabled      bool
	URL          string
	ListSelector string
	MinFields    int
	Schedule     *string
	PageCursor   *PageCursor
	ItemKey      *string
	Collection   *string
	Subscribe    bool
	Fields       map[string]Field
	Headers      map[string]string
	Visited      map[string]bool
	VisitedItems map[string]bool
}

// PageCursor visit page by identity
type PageCursor struct {
	URLFormat string
	Start     int
	End       int
}

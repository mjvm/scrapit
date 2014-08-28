package scrapit


type Scraper interface {
	Initialize() (chan *Scrapit, error)
	ScrapMe(chan*Scrapit)
	MatchResource(string) (bool)
	Uname() string
	ScraperType() SourceType
}


type SourceType uint16


const (
	_ = iota
	SOURCE_HTML
	SOURCE_XML
	SOURCE_PDF
)


type Scrapit struct {
	Source    string
	Tags      map[string]string
	sentences []string
	sourceType SourceType
	data []byte
}

//

type scrapInput interface {
	InputGenerator(chan *Scrapit)
}

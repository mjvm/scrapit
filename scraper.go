package scrapit


type Fetcher interface {
	Initialize() (chan *Scrapit, error)
	FetchMe(chan*Scrapit)
	Uname() string
	FetcherType() SourceType
}

type SourceType uint8

const (
	_ = iota
	SOURCE_WEBPAGE
)

type Scrapit struct {
	Source    string
	Tags      map[string]string
	sentences []string
	SourceType SourceType
}

//

type scrapInput interface {
	InputGenerator(chan *Scrapit)
}

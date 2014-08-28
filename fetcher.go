package scrapit

type Fetcher interface {
	Initialize() (chan *Scrapit, error)
	MatchResource(string) (bool)
	FetchMe(chan *Scrapit)
	Uname() (string)
}

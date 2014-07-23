package scrapit


type webscraper struct {
	input chan *Scrapit
	uname string
}


func WebScraper() (Fetcher) {
	w := new(webscraper)
	w.uname = "webscraper"
	return w
}

func (w *webscraper) Initialize() (c chan *Scrapit, e error) {
	c = make(chan *Scrapit, 100)
	w.input = c
	return c, nil
}

func (w *webscraper) FetchMe(o chan *Scrapit) {
	wk := NewWikiScraper()
	for ;; {
		d := <- w.input		
		wk.GoGoGadget(d, o)
	}
}

func (w *webscraper) Uname() (string) {
	return w.uname
}

func (w *webscraper) FetcherType() (SourceType) {
	return SOURCE_WEBPAGE
}

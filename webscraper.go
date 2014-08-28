package scrapit


type webscraper struct {
	input chan *Scrapit
	uname string
}


func WebScraper() (Scraper) {
	w := new(webscraper)
	w.uname = "webscraper"
	return w
}


func (w *webscraper) Initialize() (c chan *Scrapit, e error) {
	c = make(chan *Scrapit, 100)
	w.input = c
	return c, nil
}


func (w *webscraper) ScrapMe(o chan *Scrapit) {
	wk := NewWikiScraper()
	for ;; {
		d := <- w.input
		wk.GoGoGadget(d, o)
	}
}


func (w *webscraper) Uname() (string) {
	return w.uname
}


func (w *webscraper) MatchResource(r string) (bool) {
	return (true)
}


func (w *webscraper) ScraperType() (SourceType) {
	return SOURCE_HTML
}

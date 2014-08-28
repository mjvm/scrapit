package scrapit

import (
	"log"
	"fmt"
)

type sourceScrap struct {
	s Scraper
	c chan *Scrapit
}

type sourceFetch struct {
	f Fetcher
	c chan *Scrapit
}


type router struct {
	fetchMod []*sourceFetch
	scrapMod map[SourceType][]*sourceScrap
	storeMod []chan *Scrapit
	comms [3]chan *Scrapit
}


func InitRouter() (r *router) {
	r = new(router)
	r.scrapMod = make(map[SourceType][]*sourceScrap)
	for i:=0; i<3; i++ {
		r.comms[i] = make(chan *Scrapit, 100)
	}
	return
}


func (r *router) AddFetcher(f Fetcher) (e error) {
	if c, err := f.Initialize(); err != nil {
		return fmt.Errorf("Unable to initialize source", f.Uname())
	} else {
		r.fetchMod = append(r.fetchMod, &sourceFetch{f: f, c: c});
	}
	go f.FetchMe(r.comms[1])
	return
} 


func (r *router) AddScraper(f Scraper) (e error) {
	s := f.ScraperType()
	if c, err := f.Initialize(); err != nil {
		return fmt.Errorf("Unable to initialize source", f.Uname())
	} else {
		r.scrapMod[s] = append(r.scrapMod[s], &sourceScrap{s: f, c: c});
	}
	go f.ScrapMe(r.comms[2])
	return
}


func (r *router) AddStorage(f Storage) (e error) {
	if c, err := f.Initialize(); err != nil {
		return fmt.Errorf("Unable to initialize storage", f.Uname())
	} else {
		r.storeMod = append(r.storeMod, c)
	}
	go f.Process()
	return (nil)
}


func (r *router) RouteMe() {
	for ;; {
		select {
		case elem := <- r.comms[0]:
			processed := false
			for _, s := range r.fetchMod {
				if s.f.MatchResource(elem.Source) {
					s.c <- elem
					processed = true
					break;
				}
			}
			if !processed {
				log.Println("no registered module(Fetcher) can handle this source:", elem.Source)
			}
			// fetch
		case elem := <- r.comms[1]:
			// scrap
			if t, ok := r.scrapMod[elem.sourceType]; ok {
				processed := false
				for _, s := range t {
					if s.s.MatchResource(elem.Source) {
						s.c <- elem
						processed = true
						break;
					}
				}
				if !processed {
					log.Println("no registered module(Scraper) can handle this source:", elem.Source)
				}
			} else {
				log.Println("source type has no scraper module for (!) source:", elem.Source)
			}
		case elem := <- r.comms[2]:
			// store
			for _, v := range r.storeMod {
				v <- elem
			}
		}
	}
}


func (r *router) InputChan() (chan *Scrapit) {
	return r.comms[0]
}

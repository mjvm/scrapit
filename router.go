package scrapit

import (
	"log"
	"fmt"
)

type router struct {
	h map[SourceType]chan *Scrapit
	in chan *Scrapit
	out chan *Scrapit
}


func InitRouter() (r *router) {
	r = new(router)
	r.h = make(map[SourceType]chan *Scrapit)
	r.in = make(chan *Scrapit, 100)
	r.out = make(chan *Scrapit, 100)
	return
}


func (r *router) AddFetcher(f Fetcher) (e error) {
	s := f.FetcherType()
	if _, ok := r.h[s]; ok {
		return fmt.Errorf("Duplicate source", s)
	}
	if c, err := f.Initialize(); err != nil {
		return fmt.Errorf("Unable to initialize source", f.Uname())
	} else {
		r.h[s] = c
	}
	go f.FetchMe(r.out)
	return
}


func (r *router) RouteMe() {
	for ;; {
		select {
		case elem := <- r.in:
			if s, ok := r.h[elem.SourceType]; ok {
				s <- elem
			} else {
				log.Println("Unknown source type")
			}
		case elem := <- r.out:
			for _, v := range elem.sentences {
				log.Println(v)
			}
		}
	}
}


func (r *router) InputChan() (chan *Scrapit) {
	return r.in
}

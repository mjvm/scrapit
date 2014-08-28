package scrapit

import (
	"strings"
	"net/http"
	"regexp"
	"log"
)


type httpFetcher struct {
	input chan *Scrapit
	re map[SourceType]*regexp.Regexp
}


func HttpFetcher() (Fetcher) {
	h := new(httpFetcher)
	return (h)
}


func (h *httpFetcher) Initialize() (chan *Scrapit, error) {
	c := make(chan *Scrapit, 100)
	h.input = c
	h.re = make(map[SourceType]*regexp.Regexp)
	h.re[SOURCE_HTML], _ = regexp.Compile(`text/html`)
	h.re[SOURCE_XML], _ = regexp.Compile(`application/xml`)
	h.re[SOURCE_PDF], _ = regexp.Compile(`application/pdf`)
	return c, nil
}


func (h *httpFetcher) Uname() (string) {
	return "fetch-http"
}


func (h *httpFetcher) MatchResource(r string) (bool) {
	if strings.HasPrefix(r, "http://") {
		return (true)
	}
	if strings.HasPrefix(r, "https://") {
		return (true)
	}
	return (false)
}


func (h *httpFetcher) fetchAndClassify(s *Scrapit) (bool) {
	resp, err := http.Get(s.Source)
	if err != nil {
		log.Println("error crawling ", s.Source, ":", err)
		return (false)
	}
	defer resp.Body.Close()
	if header, ok := resp.Header["Content-Type"]; ok {
		// we understand content-type or dont we?
		for source, re := range h.re {
			if re.Match([]byte(header[0])) {
				s.sourceType = source
				resp.Body.Read(s.data)
				return (true)
			}
		}
	}
	// default is HTML, yeah right
	s.sourceType = SOURCE_HTML
	resp.Body.Read(s.data)
	return (true)
}


func (h *httpFetcher) FetchMe(out chan *Scrapit) {
	for ;; {
		s := <- h.input
		if h.fetchAndClassify(s) {
			out <- s
		}
	}
}

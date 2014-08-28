package main

import(
	"github.com/tust13/scrapit"
	//"time"
)


func fakeEntries(c chan *scrapit.Scrapit) {
	s := &scrapit.Scrapit{
		Source: "http://en.wikipedia.org/wiki/Category:Hairdressing",
		Tags: map[string]string{"Category": "hairdressing"},
	}
	c <- s
	s = &scrapit.Scrapit{
		Source: "https://www.kernel.org/doc/ols/2002/ols2002-pages-479-495.pdf",
		Tags: map[string]string{"Category": "datapron"},
	}
	c <- s
	//for ;; {
	//	time.Sleep(10000)
	//}
}


func main() {	
	r := scrapit.InitRouter()
	r.AddFetcher(scrapit.HttpFetcher())
	r.AddScraper(scrapit.WebScraper())
	r.AddStorage(scrapit.ElasticSearchStorage())
	//r.AddStorage(scrapit.WordFreqPGStorage("/home/ghfg/projs/go/src/github.com/tust13/scrapit/files/sentence_index.json"))
	go fakeEntries(r.InputChan())
	r.RouteMe()
}

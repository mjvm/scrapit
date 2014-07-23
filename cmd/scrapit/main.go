package main

import(
	"github.com/tust13/scrapit"
	"time"
)


func fakeEntries(c chan *scrapit.Scrapit) {
	s := &scrapit.Scrapit{
		Source: "http://en.wikipedia.org/wiki/Category:Hairdressing",
		Tags: map[string]string{"Category": "hairdressing"},
		SourceType: scrapit.SOURCE_WEBPAGE,
	}
	c <- s
	for ;; {
		time.Sleep(10000)
	}
}


func main() {	
	r := scrapit.InitRouter()
	r.AddFetcher(scrapit.WebScraper())
	go fakeEntries(r.InputChan())
	r.RouteMe()
}

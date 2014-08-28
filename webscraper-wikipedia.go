package scrapit


import (
	"net/http"
	"log"
	"regexp"
)

import (
	"code.google.com/p/go.net/html"
	"launchpad.net/xmlpath"
)


//var base_url = string("http://en.wikipedia.org/wiki/")


var divids = map[string]bool{
	"toc": true,
	"siteSub": true,
	"jump-to-nav": true,
}


type wikiscraper struct {
	
}

type wikiscrapit struct {
	Scrapit
	text string
}



func NewWikiScraper() (s *wikiscraper) {
	s = new(wikiscraper)
	return 
}


func (s *wikiscrapit) gettext(n *html.Node, p *html.Node, t bool) bool {
	if n.Type == html.ElementNode {
		b := []byte(n.Data)
		if n.Data == "body" {
			t = true
		} else if (b[0] == 72 || b[0] == 104) && (b[1] >= 46 && b[1] < 57) {
			t = false
		} else if n.Data == "script" {
			return false
		} else if n.Data == "div" {
			for k := range n.Attr {
				if n.Attr[k].Key == "class" && n.Attr[k].Val == "hatnote relarticle mainarticle" {
					return false
				} else if n.Attr[k].Key == "id" {
					if _, ok := divids[n.Attr[k].Val]; ok {
						return false
					}
				}
			}
		} else if n.Data == "sup" {
			return false
		} else if n.Data == "span" {
			giveup := 0
			for k := range n.Attr {
				if n.Attr[k].Key == "class" && n.Attr[k].Val == "mw-headline" {
					giveup++
				} else if n.Attr[k].Key == "id" && n.Attr[k].Val == "See_also" {
					giveup++
				} else if n.Attr[k].Key == "id" && n.Attr[k].Val == "References" {
					giveup++
				} else if n.Attr[k].Key == "id" && n.Attr[k].Val == "External_links" {
					giveup++
				}
			}
			if giveup == 2 {
				return true
			}
		}
	}
	if t && n.Type == html.TextNode {
		s.text += n.Data
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if quit := s.gettext(c, n, t); quit {
			return true
		}
	}
	return false
}


func divideToConquer(text []byte) (r []string) {
	b := make([]byte, 0)
	for i:= 0; i<len(text); i++ {

		switch text[i] {
		case byte(10):
			if len(b) == 0 {
				continue
			}
			r = append(r, string(b))
			b = make([]byte, 0)
		case byte(46):
			if text[i+1] == byte(20) {
				b = append(b, text[i])
				r = append(r, string(b))
				b = make([]byte, 0)
				i++
			}
		case byte(32):
			b = append(b, text[i])
			for ;;i++ {
				if text[i] != byte(32) {
					break
				}
			}
			i--
		case byte(9):
		default:
			b = append(b, text[i])
		}
		
	}
	return
}


func (s *wikiscrapit)extractPageLinks(base_url string) (pages []*wikiscrapit, err error) {
	path, err := xmlpath.Compile("//div[@id=\"mw-pages\"]//li/a/@href")
	if err != nil {
		log.Printf("Error compiling xpath\n")
		return nil, err
	}
	resp, err := http.Get(s.Source)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	rootNode, err := xmlpath.ParseHTML(resp.Body)
	if err != nil {
		return nil, err
	}
	iter := path.Iter(rootNode)
	for iter.Next() {
		node := iter.Node()
		pages = append(pages, &wikiscrapit{
			Scrapit: Scrapit{
				Source: base_url + node.String(),
				Tags: s.Tags,
			},})
		log.Printf("Expanded to %s\n", base_url + node.String())
	}
	return pages, nil
}



func (s *wikiscrapit) ScrapMe() {
	resp, err := http.Get(s.Source)
	if err != nil {
		log.Println("error crawling ", s.Source)
		return
	}
	defer resp.Body.Close()
	
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return
	}
	s.gettext(doc, nil, false)
	for _, v := range divideToConquer([]byte(s.text)) {
		s.sentences = append(s.sentences, v)
	}
	return
}



func (s *wikiscraper) GoGoGadget(sc *Scrapit, o chan *Scrapit) {
	wk := &wikiscrapit{
		Scrapit: *sc,
	}
	re, _ := regexp.Compile(`^(http:\/\/[^\/]+)\/wiki\/Category:`)
	if match := re.FindAllStringSubmatch(wk.Source, -1); match != nil {
		pages, _ := wk.extractPageLinks(match[0][1])
		for _, v := range pages {
			v.ScrapMe()
			o <- &v.Scrapit
		}
	} else {
		wk.ScrapMe()
		o <- &wk.Scrapit
	}
}

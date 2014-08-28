package scrapit

import (
	"log"
)

import (
	elastigo "github.com/mattbaird/elastigo/lib"
)

type EsSession struct {
	conn *elastigo.Conn
	input chan *Scrapit
}


func ElasticSearchStorage() (Storage) {
	s := new(EsSession)
	s.conn = elastigo.NewConn()
	s.conn.Domain = "127.0.0.1"
	return s
}


func (s *EsSession) Initialize() (chan *Scrapit, error) {
	c := make(chan *Scrapit, 100)
	s.input = c
	return c, nil
}


func (s *EsSession) Uname() (string) {
	return "elasticsearch"
}


func (s *EsSession) Process() {
	for {
		sc := <- s.input
		for _, v := range sc.sentences {
			//log.Println(v.Sentence)
			if _, err := s.conn.Index("sentences", "sentence", "", nil, v); err != nil {
				//log.Println("error: ", err)
			}
		}
		log.Println("Elastic searched: ", sc.Source)
	}
}


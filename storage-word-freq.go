package scrapit

import (
	"bufio"
        "database/sql"
	"log"
	"os"
	"strings"

        _ "github.com/lib/pq"
	"github.com/deckarep/golang-set"
)


type pg struct {
	stopwords mapset.Set
	input chan *Scrapit
}


func createStopWordList(stopwordsFile string) (mapset.Set, error) {

	stopwords := mapset.NewSet()

	indexFile, err := os.Open(stopwordsFile)
	if err != nil {
		log.Println("Error open stopwords file")
		return nil, err
	}
	defer indexFile.Close()

	reader := bufio.NewReader(indexFile)
	for {
		word, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		//log.Printf("Add stop word: %s", word)
		stopwords.Add(strings.Trim(string(word), "\n"))
	}
	return stopwords, nil
}


func WordFreqPGStorage(file string) (Storage) {
	p := new(pg)
	stopwords, err := createStopWordList(file)
	if err != nil {
		log.Println("Error creating stop list")
		return (nil)
	}
	p.stopwords = stopwords
	return (p)
}


func (p *pg) Initialize() (chan *Scrapit, error) {
	c := make(chan *Scrapit, 100)
	p.input = c
	return c, nil
}


func (p *pg) Uname() (string) {
	return "postgres-word-freq"
}


func (p *pg) UpdateIndexStats(scraper *Scrapit) {
	//log.Printf("stopwords size: %d", p.stopwords.Cardinality())
	for i := 0; i < len(scraper.sentences); i++ {
		sentence := scraper.sentences[i]
		terms := strings.Split(sentence, " ")
		m := make(map[string]int)
		for j := 0; j < len(terms); j++ {
			term := strings.Trim(terms[j], " ")
			term = strings.Trim(term, "(")
			term = strings.Trim(term, ")")
			term = strings.Trim(term, ",")
			if strings.Contains(term, "mwloaderload") {
				continue
			}
			term = strings.ToLower(term)
			if p.stopwords.Contains(term) {
				//log.Println("drop stop word: %s", term)
				continue
			}
			if len(term) > 64 || len(term) < 3 {
				continue
			}
			counter, _ := m[term]
			m[term] = counter + 1
		}
		UpdateRecord(p.stopwords, scraper.Source, m, scraper.Tags["Category"])
	}
}


func UpdateRecord(stopwords mapset.Set, document string, terms map[string]int, category string) {
	log.Println("processing doc: %s", document)
	db, err := sql.Open("postgres", "user=sg dbname=sg password=sg sslmode=disable")
	if err != nil {
		log.Printf("Error open postgres conn: %#v", err)
	}
	defer db.Close()
	for term, counter := range terms {
		txn, err := db.Begin()
		if err != nil {
			log.Printf("Error starting a transaction %#v\n", err)
		}
		rows, err := txn.Query("SELECT 1 FROM doc_term_frequency WHERE term=$1", term)
		if err != nil {
			log.Printf("Error executing select%#v\n", err)
		}
		row := rows.Next()
		if row == false {
			_, err := txn.Exec("INSERT into doc_term_frequency (term, counter) VALUES ($1, $2)", term, counter)
			if err != nil {
				log.Panicf("Error inserting term %s\n", term)
			}
			_, err = txn.Exec("INSERT into term_frequency (document, term, category, counter) VALUES ($1, $2, $3, $4)",
				document, term, category, counter)
			if err != nil {
				log.Panicf("Error inserting term %s for document: %s error: %#v\n", term, document, err)
			}
		} else {
			_, err := txn.Exec("UPDATE doc_term_frequency SET counter = counter + $1 where term = $2",
				counter, term)
			if err != nil {
				log.Panicf("Error updating term %s error: %#v\n", term, err)
			}
			_, err = txn.Exec("UPDATE term_frequency SET counter = counter + $1 where document = $2 and term = $3 and category = $4",
				counter, document, term, category)
			if err != nil {
				log.Panicf("Error updating term %s for document %s error: %#v\n", term, document, err)
			}
		}

		err = txn.Commit()
		if err != nil {
			log.Fatal(err)
		}
	}
}


func (p *pg) Process() {
	for {
		d := <-p.input
		p.UpdateIndexStats(d)
	}
}

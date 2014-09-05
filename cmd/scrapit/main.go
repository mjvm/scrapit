package main

import (
	"flag"
	"log"
)

import (
	"code.google.com/p/gcfg"
)

import (
	"github.com/tust13/scrapit"
)

type AppConfig struct {
	Host     string
	Port     int
	Debug    bool
	BaseHost string
}

type Config struct {
	Application   AppConfig
	StopWordsFile string
}

func readConfig(configFile string) (*Config, error) {
	var config = new(Config)
	err := gcfg.ReadFileInto(config, configFile)
	if err != nil {
		log.Fatalf("Error reading config file: %s\n", err)
		return nil, err
	}

	if config.Application.Host == "" {
		config.Application.Host = "localhost"
	}

	if config.Application.Port == 0 {
		config.Application.Port = 8080
	}

	return config, nil
}

func fakeEntries(c chan *scrapit.Scrapit) {
	s := &scrapit.Scrapit{
		Source: "http://en.wikipedia.org/wiki/Category:Hairdressing",
		Tags:   map[string]string{"Category": "hairdressing"},
	}
	c <- s
	s = &scrapit.Scrapit{
		Source: "https://www.kernel.org/doc/ols/2002/ols2002-pages-479-495.pdf",
		Tags:   map[string]string{"Category": "datapron"},
	}
	c <- s
	//for ;; {
	//	time.Sleep(10000)
	//}
}

func Init() *Config {
	var configFile, stopWords string
	flag.StringVar(&configFile, "c", "/etc/scrapit.conf", "configuration file")
	flag.StringVar(&stopWords, "s", "/etc/stop_words.txt", "stop words")
	flag.Parse()

	config, err := readConfig(configFile)
	if err != nil {
		log.Panicf("Unable to read config file: %s\n", err)
	}
	config.StopWordsFile = stopWords
	return config
}

func main() {
	config := Init()
	r := scrapit.InitRouter()
	r.AddFetcher(scrapit.HttpFetcher())
	r.AddScraper(scrapit.WebScraper())
	r.AddStorage(scrapit.ElasticSearchStorage())
	r.AddStorage(scrapit.WordFreqPGStorage(config.StopWordsFile))
	go fakeEntries(r.InputChan())
	r.RouteMe()
}

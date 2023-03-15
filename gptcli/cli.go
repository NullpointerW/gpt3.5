package gptcli

import (
	"github.com/sashabaranov/go-openai"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Token struct {
	Context  []openai.ChatCompletionMessage
	LastTime time.Time
}

var (
	apiKey   = "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"
	Cli      = newCli()
	TokenMap = sync.Map{}
)

func tokensCleaner(d time.Duration) {
	log.Println("token cleaner start")
clean:
	its := func(k, v interface{}) bool {
		tk := v.(*Token)
		if time.Now().Sub(tk.LastTime) >= d {
			TokenMap.Delete(k)
			log.Printf("clean token %s", k.(string))
		}
		return true
	}
	TokenMap.Range(its)

	for {
		time.Sleep(d)
		goto clean
	}
}
func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	go tokensCleaner(time.Minute * 30)
}

func newCli() *openai.Client {

	config := openai.DefaultConfig(apiKey)
	proxyUrl, err := url.Parse("http://localhost:7890")
	if err != nil {
		panic(err)
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}
	config.HTTPClient = &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
	return openai.NewClientWithConfig(config)
}

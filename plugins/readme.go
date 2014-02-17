package plugins

import (
	"github.com/google/go-github/github"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"log"
	"net/http"
	"encoding/base64"
)

type Readme struct {
	Html string `json:"html"`
}

func (p *Plugin) rawReadme() []byte {
	client := github.NewClient(nil)
	content, _, err := client.Repositories.GetReadme(p.GithubUser, p.GithubRepo)
	if err != nil {
		panic(err)
	}
	markdown, err := base64.StdEncoding.DecodeString(*content.Content)
	if err != nil {
		panic(err)
	}
	return markdown
}

func HttpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println("HTTP Error:", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return body, nil
}

func (p *Plugin) Readme() *Readme {
	markdown := p.rawReadme()
	if markdown == nil {
		return nil
	}
	html := blackfriday.MarkdownCommon(markdown)
	return &Readme{string(html)}
}

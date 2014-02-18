package plugins

import (
	"github.com/google/go-github/github"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/dickeyxxx/githubmarkdown"
	"encoding/base64"
)

func (p *Plugin) rawReadme() string {
	client := github.NewClient(nil)
	content, _, err := client.Repositories.GetReadme(p.GithubUser, p.GithubRepo)
	if err != nil {
		panic(err)
	}
	markdown, err := base64.StdEncoding.DecodeString(*content.Content)
	if err != nil {
		panic(err)
	}
	return string(markdown)
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

func (p *Plugin) Readme() string {
	markdown := p.rawReadme()
	if markdown == "" {
		return ""
	}
	html, err := githubmarkdown.Parse(markdown)
	if err != nil {
		log.Println("Error parsing markdown:", err)
		return ""
	}
	return string(html)
}

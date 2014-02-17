package plugins

import (
	"github.com/russross/blackfriday"
	"io/ioutil"
	"log"
	"net/http"
)

type Readme struct {
	Html string `json:"html"`
}

func (p *Plugin) rawReadme() []byte {
	url := p.ReadmeUrl()
	readme, _ := HttpGet(url)
	return readme
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

func (*Plugin) ReadmeUrl() string {
	return "https://raw.github.com/tpope/vim-surround/master/README.markdown"
}

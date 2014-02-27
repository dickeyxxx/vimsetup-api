package scraper

import (
	"code.google.com/p/go.net/html"
	"github.com/dickeyxxx/vimsetupapi/cache"
	"github.com/dickeyxxx/vimsetupapi/plugins"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

type Fetcher struct {
	Cache cache.Cacher
}

func (f *Fetcher) Fetch(url string) string {
	cache := f.Cache.Get(url)
	if cache != "" {
		log.Println("Cache hit for:", url)
		return cache
	}
	log.Println("Cache miss for:", url)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	text := string(content)
	f.Cache.Set(url, text, 86400)
	return text
}

func Run(log *log.Logger, mongo *mgo.Session) {
	fetcher := &Fetcher{
		Cache: cache.NewRedisCache(),
	}
	log.Println("Starting scrape...")
	host := "http://vam.mawercer.de/"
	links := make(chan string, 10)
	go getLinks(host, links, fetcher, log)
	var wg sync.WaitGroup
	wg.Add(20)
	go func() {
		for i := 0; i < 20; i++ {
			for l := range links {
				plugin := getPlugin(host+l, fetcher, log)
				savePlugin(plugin, mongo)
			}
			wg.Done()
		}
	}()
	wg.Wait()
	log.Println("Finished scrape")
}

func getPlugin(url string, fetcher *Fetcher, log *log.Logger) *plugins.Plugin {
	log.Println("Getting plugin from", url)
	text := fetcher.Fetch(url)
	name := parseName(text)
	return &plugins.Plugin{Name: name}
}

func parseName(text string) string {
	rp := regexp.MustCompile(`.+Plugin:\s+(?P<name>.+)\n`)
	name := rp.FindStringSubmatch(text)[1]
	name = strings.Split(name, "version")[0]
	return strings.TrimSpace(name)
}

func getLinks(url string, links chan string, fetcher *Fetcher, log *log.Logger) {
	log.Println("Getting plugins from", url)
	body := fetcher.Fetch(url)
	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		panic(err)
	}
	r := regexp.MustCompile(`^\?plugin_info.*$`)
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					if r.MatchString(a.Val) {
						links <- a.Val
					}
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	close(links)
	log.Println("Found all plugins")
}

func savePlugin(plugin *plugins.Plugin, mongo *mgo.Session) {
	pluginCollection(mongo).Upsert(bson.M{"Name": plugin.Name}, bson.M{"Name": plugin.Name})
}

func pluginCollection(mongo *mgo.Session) *mgo.Collection {
	return mongo.DB("").C("plugins")
}

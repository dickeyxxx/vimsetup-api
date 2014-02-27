package main

import (
	"github.com/dickeyxxx/vimsetupapi/version"
	"github.com/dickeyxxx/vimsetupapi/scraper"
	"github.com/dickeyxxx/vimsetupapi/db"
	"github.com/codegangsta/cli"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "vimsetup"
	app.Version = version.Version()
	app.Usage = "administer vimsetup.com"
	app.Commands = []cli.Command{
		cli.Command{
			Name:  "scrape",
			Usage: "scrapes http://vam.mawercer.de/ for new addons",
			Action: func(c *cli.Context) {
				logger := log.New(os.Stdout, "Scrape: ", log.Ldate|log.Ltime|log.Lshortfile)
				scraper.Run(logger, db.MongoSession())
			},
		},
	}

	app.Run(os.Args)
}

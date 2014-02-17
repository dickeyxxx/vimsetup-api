package main

import (
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/cors"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"os"
)

type Plugin struct {
	Name string
	Github string
}

func main() {
	mongo := mongoSession()
	defer mongo.Close()
	Serve(mongo)
}

func Serve(mongo *mgo.Session) {
	m := martini.Classic()
	m.Use(render.Renderer())

	m.Use(cors.Allow(&cors.Options{
		AllowOrigins:     []string{"http://localhost*", "http://vimsetup.com"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
	}))

	m.Get("/plugins", func(r render.Render) {
		plugins := Plugins(mongo)
		r.JSON(200, plugins)
	})

	m.Run()
}

func Plugins(mongo *mgo.Session) []Plugin {
	res := []Plugin{}
	plugins := mongo.DB("").C("plugins")
	err := plugins.Find(bson.M{}).All(&res)
	if err != nil {
		panic(err)
	}
	return res
}

func mongoSession() *mgo.Session {
	uri := os.Getenv("MONGOHQ_URL")
	if uri == "" {
		uri = "mongodb://localhost/vimsetup"
	}
	mongo, err := mgo.Dial(uri)
	if err != nil {
		panic(err)
	}

	return mongo
}

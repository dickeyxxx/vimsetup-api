package main

import (
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/cors"
	"github.com/martini-contrib/encoder"
	"github.com/dickeyxxx/vimsetupapi/plugins"
	"github.com/dickeyxxx/vimsetupapi/db"
	"labix.org/v2/mgo"
	"net/http"
	"github.com/kdar/factorlog"
	"os"
)

func main() {
	mongo := db.MongoSession()
	log := factorlog.New(os.Stdout, factorlog.NewStdFormatter("%{Date} %{Time} %{File}:%{Line} %{Message}"))
	defer mongo.Close()
	Serve(mongo, log)
}

func Serve(mongo *mgo.Session, log factorlog.Logger) {
	m := martini.New()
	m.Map(mongo)
	m.Map(log)
	m.Use(martini.Recovery())
	m.Use(martini.Logger())

	m.Use(func(c martini.Context, w http.ResponseWriter) {
		c.MapTo(encoder.JsonEncoder{}, (*encoder.Encoder)(nil))
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	})

	m.Use(cors.Allow(&cors.Options{
		AllowOrigins:  []string{"http://localhost*", "http://vimsetup.com"},
		AllowMethods:  []string{"GET", "POST"},
		AllowHeaders:  []string{"Origin"},
		ExposeHeaders: []string{"Content-Length"},
	}))

	m.Action(Router().Handle)

	m.Run()
}

func Router() martini.Router {
	router := martini.NewRouter()

	router.Get("/", func(enc encoder.Encoder) (int, []byte) {
		doc := map[string]interface{}{"Foo": "bar"}
		return http.StatusOK, encoder.Must(enc.Encode(doc))
	})

	router.Get("/plugins", func(enc encoder.Encoder, mongo *mgo.Session) (int, []byte) {
		plugins := plugins.All(mongo)
		return http.StatusOK, encoder.Must(enc.Encode(plugins))
	})

	router.Get("/plugins/:author/:name", func(params martini.Params, enc encoder.Encoder, mongo *mgo.Session) (int, []byte) {
		plugin := plugins.FindByName(mongo, params["name"])
		if plugin != nil {
			return http.StatusOK, encoder.Must(enc.Encode(plugin))
		} else {
			return notFound()
		}
	})

	router.Get("/plugins/:author/:name/readme", func(params martini.Params, enc encoder.Encoder, mongo *mgo.Session) (int, []byte) {
		plugin := plugins.FindByName(mongo, params["name"])
		if plugin == nil {
			return notFound()
		}
		readme := plugin.Readme()
		return http.StatusOK, encoder.Must(enc.Encode(readme))
	})

	return router
}

func notFound() (int, []byte) {
	return http.StatusNotFound, []byte{}
}

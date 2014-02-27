package plugins

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Plugin struct {
	Name   string `json:"name"`
	GithubUser string `bson:"github_user" json:"github-user"`
	GithubRepo string `bson:"github_repo" json:"github-repo"`
}

func FindByName(mongo *mgo.Session, name string) *Plugin {
	plugin := Plugin{}
	err := pluginCollection(mongo).Find(bson.M{"name": name}).One(&plugin)
	if err != nil {
		switch err.Error() {
		case "not found":
			return nil
		default:
			panic(err)
		}
	}
	return &plugin
}

func All(mongo *mgo.Session) []Plugin {
	res := []Plugin{}
	err := pluginCollection(mongo).Find(bson.M{}).All(&res)
	if err != nil {
		panic(err)
	}
	return res
}

func pluginCollection(mongo *mgo.Session) *mgo.Collection {
	return mongo.DB("").C("plugins")
}

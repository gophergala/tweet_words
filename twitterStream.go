package tweet_words


import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/robfig/config"
	"gopkg.in/mgo.v2"
	"net/url"
	"time"
)

var TwitterApi *anaconda.TwitterApi
var conf map[string]string

func init() {
	myconf, err := config.ReadDefault("app.properties")
	if err != nil {
		panic(err)
	}
	conf = make(map[string]string)
	var keys = []string{"CONSUMER_KEY", "CONSUMER_SECRET", "ACCESS_TOKEN", "ACCESS_TOKEN_SECRET", "MONGO"}
	for key := range keys {
		conf[keys[key]], err = myconf.String("", keys[key])
		if err != nil {
			panic(err)
		}
	}
	anaconda.SetConsumerKey(conf["CONSUMER_KEY"])
	anaconda.SetConsumerSecret(conf["CONSUMER_SECRET"])
	TwitterApi = anaconda.NewTwitterApi(conf["ACCESS_TOKEN"], conf["ACCESS_TOKEN_SECRET"])
}

func Tweets(query url.Values, timeout time.Duration, quit chan bool) <-chan anaconda.Tweet {
	stream, err := TwitterApi.PublicStreamFilter(query)
	if err != nil {
		panic(err)
	}
	var tweet anaconda.Tweet
	var junk interface{}
	tweetChan := make(chan anaconda.Tweet, 1024) // as much as I like the consumer routine...
	go func() {
		quitter := time.After(timeout)
		for {
			select {
			case junk = <-stream.C:
				switch junk.(type) {
				case anaconda.Tweet:
					tweet = junk.(anaconda.Tweet)
					tweetChan <- tweet
				}
			case <-quitter:
				quit <- true
				return
			}
		}
	}()
	return tweetChan
}

func StoreTweets(query url.Values, timeout time.Duration, collectionName string) (retChan chan bool) {
 retChan = make(chan bool)
 quit := make(chan bool)
 z := Tweets(query, timeout, quit)
 go func() {
 	mgoSession, err := mgo.Dial(conf["MONGO"])
 	if err != nil {
 		panic(err)
 	}
  mgoSession.SetMode(mgo.Monotonic, true)
  defer mgoSession.Close()
	for {
		select {
		case x := <-z:
			newSes := mgoSession.Copy()
			defer newSes.Close()
			col := newSes.DB("test").C(collectionName)
			if (col == nil) {
				panic("unable to get collection")
			}
			err = col.Insert(x)
			if err != nil {
				panic(err)
			}
		case <-quit:
			retChan <- true
			return
		}
	}
 }()
 return
}


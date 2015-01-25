package tweet_words

import (
		"gopkg.in/mgo.v2"
		"gopkg.in/mgo.v2/bson"
		"fmt"
		"log"
		)

type User struct {
        Token string
        Secret string
        Keywords []string
}

var GUser User

func StoreUser(user User) (ret bool) {
		mgoSession, err := mgo.Dial("mongodb://54.188.201.254")
		if err != nil {
			panic(err)
		}
		mgoSession.SetMode(mgo.Monotonic, true)
		defer mgoSession.Close()
		newSes := mgoSession.Copy()
		defer newSes.Close()
		db := newSes.DB("test")
		r, _ := db.CollectionNames()
		fmt.Println(r)
		col := db.C("User")
		if col == nil {
			panic("unable to get collection")
		}
		result := User{}
		fmt.Println(bson.M{"token": user.Token})
		err = col.Find(bson.M{"token": user.Token}).One(&result)
		fmt.Println("Token1:", result)
		if result.Token == "" {
			fmt.Println(user)
			dummy := []string{"bjp"}
			err = col.Insert(&User{user.Token, user.Secret, dummy})
			if err != nil {
				panic(err)
			}
		}
		err = col.Find(bson.M{"token": user.Token}).One(&result)
		if err != nil {
			fmt.Println("sdfsd")
            log.Fatal(err)
        }
        r, _ = db.CollectionNames()
        fmt.Println(r)
        fmt.Println("Token StoreUser:", result)
	return
}

func StoreKeywords(data string) (ret bool) {
	mgoSession, err := mgo.Dial("mongodb://54.188.201.254")
	if err != nil {
		panic(err)
	}
	mgoSession.SetMode(mgo.Monotonic, true)
	defer mgoSession.Close()
	newSes := mgoSession.Copy()
	defer newSes.Close()
	col := newSes.DB("test").C("User")
	if col == nil {
		panic("unable to get collection")
	}
	result := User{}
	err = col.Find(bson.M{"token": GUser.Token}).One(&result)
	if err != nil {
        log.Fatal(err)
    }
    n := len(result.Keywords)
    words := make([]string, n+1);
    copy(words, result.Keywords[0:])
    words[n] = data
    _, err = col.Upsert(bson.M{"token": GUser.Token, "secret" : GUser.Secret}, bson.M{"$set" : bson.M{"keywords": words}})
    if err != nil {
    	fmt.Println(err)
    }
    result1 := User{}
	err = col.Find(bson.M{"token": GUser.Token}).One(&result1)
	if err != nil {
        log.Fatal(err)
    }
    fmt.Println("token:", result1)

	return
}
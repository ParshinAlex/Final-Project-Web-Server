package main

import (
	"fmt"

	"gopkg.in/mgo.v2"
)

// Ses allows to work with the database
var Ses *mgo.Session

// DbUsers is the key variable
var DbUsers *mgo.Collection

// DbArticles is the key variable
var DbArticles *mgo.Collection

// MongoInit initialises Mongo database for upcoming operations
func MongoInit() {
	var err error
	Ses, err = mgo.Dial("mongodb://127.0.0.1")
	if err != nil {
		fmt.Println(err.Error())
	}
	DbUsers = Ses.DB("Project").C("registeredUsers")
	DbArticles = Ses.DB("Project").C("articles")
}

//AddUser adds user to the database registeredUsers
func AddUser(us user) {
	DbUsers.Insert(us)
}

//AddArticle adds article to the database articles
func AddArticle(ar article) {
	DbArticles.Insert(ar)
}

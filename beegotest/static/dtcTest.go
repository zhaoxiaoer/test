package main

import (
	"fmt"

	"gopkg.in/mgo.v2-unstable"
	"gopkg.in/mgo.v2-unstable/bson"
)

type Person struct {
	Id    bson.ObjectId `json:"id"	bson:"_id"`
	Class string        `json:class bson:"_class"`
	Appid string        `json:appid`
	Obdid string        `json:obdid`
	Salt  string        `json:salt`
	Stime string        `json:stime`
	Otime string        `json:otime`
	Data  []interface{} `json:data`
}

func main() {
	fmt.Printf("111\n")
	session, err := mgo.Dial("192.168.1.117")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("sdea").C("obddtcbyte")
	n, err := c.Count()
	if err != nil {
		panic(err)
	}
	fmt.Printf("n: %d\n", n)

	q := c.Find(bson.M{"data.0": "2"}).Sort("-stime").Iter()
	n, err = q.Count()
	if err != nil {
		panic(err)
	}
	fmt.Printf("n: %d\n", n)
}

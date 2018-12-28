package main

import (
	"github.com/anmaslov/smdr"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"regexp"
	"strings"
	"time"
)

const (
	collection = "calls"
)

type DateInfo struct {
	DateStart    time.Time `bson:"date_start"`
	DateEnd    time.Time `bson:"date_end"`
	DateDiff	string
	SecondDiff  float64
}

type CallInfo struct{
	Stantion	string `bson:"stantion"`
	Tp 			string `bson:"tp"`
	TruncOut	string `bson:"trunc_out"`
	TruncInc	string `bson:"trunc_inc"`
	CallId		int    `bson:"call_id"`
	Tenant		string `bson:"tenant"`
	Called		string `bson:"called"`
	Cvt			DateInfo
	Route1		string `bson:"route1"`
	Route2		string `bson:"route2"`
	Phone		string `bson:"phone"`
	PhoneRaw	string `bson:"phone_raw"`
	CallMetering	string `bson:"call_metering"`
}

type Phones struct {
	Id bson.ObjectId `bson:"_id"`
	Ip string `bson:"ip"`
	Port string `bson:"port"`
}

type MongoStore struct {
	session *mgo.Session
}

var mongoStore = MongoStore{}

func initialiseMongo() (session *mgo.Session){

	info := &mgo.DialInfo{
		Addrs:    []string{cfg.Database.Host},
		Timeout:  60 * time.Second,
		Database: cfg.Database.Dbname,
		Username: cfg.Database.Username,
		Password: cfg.Database.Password,
	}

	session, err := mgo.DialWithInfo(info)
	if err != nil {
		log.Fatal("can't connect to mongoDb", err)
	}

	return
}

func fillParam(r *smdr.CDR) CallInfo{

	call := CallInfo{Tp: r.Tp, TruncOut: r.TrunkOut, TruncInc: r.TrunkInc, Called: r.Called}

	call.PhoneRaw = strings.Trim(r.Phone, " ")
	call.Phone = phoneParse(call.PhoneRaw)

	call.Route1 = r.Route1
	call.Route2 = r.Route2

	call.Cvt.DateStart = dateParse(&r.CvsStart)
	call.Cvt.DateEnd = dateParse(&r.CvsEnd)

	diff := call.Cvt.DateEnd.Sub(call.Cvt.DateStart)
	call.Cvt.DateDiff = diff.String()
	call.Cvt.SecondDiff = diff.Seconds()

	call.Tenant = r.Tenant
	call.CallMetering = r.CallMetering

	return call
}

func insertCall(call *CallInfo) error {
	var err error
	c := mongoStore.session.DB(cfg.Database.Dbname).C(collection)

	err = c.Insert(call)
	if err != nil{
		return err
		//log.Fatal("error when trying write to mongoDB", err)
	}

	//log.Println("write to DB success, date end call: ", call.Cvt.DateEnd, " dur:", call.Cvt.DateDiff)
	return nil
}

func getPhones() ([]Phones, error){
	c := mongoStore.session.DB(cfg.Database.Dbname).C("phones")

	query := bson.M{}
	query["enabled"] = bson.M{"$eq": true}

	phones := []Phones{}
	err := c.Find(query).Sort("-_id").All(&phones)

	if err != nil {
		return []Phones{}, err
	}
	return phones, nil
}

func dateParse(c *smdr.Conversation) time.Time{
	strDate := c.Year + "-" + c.Month + "-" + c.Day + "T" + c.Hour + ":" + c.Minute + ":" + c.Second
	dt, err := time.Parse("06-01-02T15:04:05Z07:00", strDate + "+03:00")

	if err != nil {
		log.Println("failed when date parse", err)
	}

	return dt
}

func phoneParse(phone string) string {

	var validID = regexp.MustCompile(`^01[01345][0134567]\d{2}\s{1,}001$`)

	if (validID.MatchString(phone)) {
		return phone[2:6]
	} else {
		return phone
	}
}
package main

import (
	"github.com/anmaslov/smdr"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
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
	CallMetering	string `bson:"call_metering"`
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

	call.Phone = strings.Trim(r.Phone, " ")
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

func insertCall(call *CallInfo) bool {
	var err error
	c := mongoStore.session.DB(cfg.Database.Dbname).C(collection)

	err = c.Insert(call)
	if err != nil{
		log.Fatal("error when trying write to mongoDB", err)
	}

	log.Println("write to DB success, date end call: ", call.Cvt.DateEnd, " dur:", call.Cvt.DateDiff)
	return true
}

func getCalls(filter *CallFilter) []CallInfo{
	// получаем коллекцию
	c := mongoStore.session.DB(cfg.Database.Dbname).C(collection)
	// критерий выборки

	query := bson.M{}

	if len(filter.tp) > 0 {
		query["tp"] = bson.M{"$eq": filter.tp}
	}

	if len(filter.phone) > 0 {
		query["phone"] = bson.M{"$regex": filter.phone}
	}

	if len(filter.stantion) > 0 {
		query["stantion"] = bson.M{"$eq": filter.stantion}
	}

	if len(filter.called) > 0 {
		query["called"] = bson.M{"$regex": filter.called}
	}

	log.Println(query)

	if len(filter.start) > 0 && len(filter.end) > 0 {
		fromDate, fErr := time.Parse("2006-01-02T15:04:05Z07:00", filter.start + "T00:00:00+03:00")
		toDate, tErr := time.Parse("2006-01-02T15:04:05Z07:00", filter.end + "T23:59:59+03:00")

		if fErr == nil || tErr == nil {
			log.Println(fromDate, "-", toDate)
			query["cvt.date_end"] =  bson.M{"$gt": fromDate, "$lt": toDate}
		}
	} /*else {
		//Если фильтр дат не задан - выберем за последний месяц
		now := time.Now()
		currentYear, currentMonth, _ := now.Date()
		currentLocation := now.Location()

		fromDate := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
		toDate := fromDate.AddDate(0, 1, -1)
		query["cvt.date_end"] =  bson.M{"$gt": fromDate, "$lt": toDate}
	}*/

	// объект для сохранения результата
	ci := []CallInfo{}
	err := c.Find(query).Sort("-_id").
		Limit(filter.limit).
		Skip(filter.skip).
		All(&ci)

	if err != nil {
		log.Println(err)
		return []CallInfo{}
	}
	return ci
}

func dateParse(c *smdr.Conversation) time.Time{
	strDate := c.Year + "-" + c.Month + "-" + c.Day + "T" + c.Hour + ":" + c.Minute + ":" + c.Second
	dt, err := time.Parse("06-01-02T15:04:05Z07:00", strDate + "+03:00")

	if err != nil {
		log.Println("failed when date parse", err)
	}

	return dt
}
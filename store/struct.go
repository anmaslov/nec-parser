package store

import (
	"github.com/anmaslov/smdr"
	"gopkg.in/mgo.v2/bson"
	"regexp"
	"strings"
	"time"
)

type DateInfo struct {
	DateStart  time.Time `bson:"date_start"`
	DateEnd    time.Time `bson:"date_end"`
	DateDiff   string
	SecondDiff float64
}

type CallInfo struct {
	Stantion     string `bson:"stantion"`
	Tp           string `bson:"tp"`
	TruncOut     string `bson:"trunc_out"`
	TruncInc     string `bson:"trunc_inc"`
	CallId       int    `bson:"call_id"`
	Tenant       string `bson:"tenant"`
	Called       string `bson:"called"`
	Cvt          DateInfo
	Route1       string `bson:"route1"`
	Route2       string `bson:"route2"`
	Phone        string `bson:"phone"`
	PhoneRaw     string `bson:"phone_raw"`
	CallMetering string `bson:"call_metering"`
}

type Phones struct {
	Id   bson.ObjectId `bson:"_id"`
	Ip   string        `bson:"ip"`
	Port string        `bson:"port"`
}

// FillParam заполенение структуры звонка
func FillParam(r *smdr.CDR) *CallInfo {

	call := &CallInfo{Tp: r.Tp, TruncOut: r.TrunkOut, TruncInc: r.TrunkInc, Called: r.Called}

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

// dateParse парсинг даты
func dateParse(c *smdr.Conversation) time.Time {
	strDate := c.Year + "-" + c.Month + "-" + c.Day + "T" + c.Hour + ":" + c.Minute + ":" + c.Second
	dt, _ := time.Parse("06-01-02T15:04:05Z07:00", strDate+"+03:00")

	return dt
}

// phoneParse парсинг телефона под определенный формат
func phoneParse(phone string) string {
	validID := regexp.MustCompile(`^01[0-5][0-9]\d{2}\s{1,}001$`)

	if validID.MatchString(phone) {
		return phone[2:6]
	}

	return phone
}

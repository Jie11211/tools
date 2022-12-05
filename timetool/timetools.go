package timetool

import (
	"time"
)

const(
	layout="2006-01-02 15:04:05"
)

type TimeTools struct {
}

func NewTimeTools() *TimeTools {
	return &TimeTools{}
}

func (t *TimeTools)SleepToNextMinute(){
	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location()).Add(time.Minute)
	time.Sleep(time.Duration(int(next.Unix()-now.Unix())) * time.Second)
}

func (t *TimeTools)SleepMinute(minute int){
	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location()).Add(time.Minute*time.Duration(minute))
	time.Sleep(time.Duration(int(next.Unix()-now.Unix())) * time.Second)
}

func (t *TimeTools)SleepToNextHour(){
	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location()).Add(time.Hour)
	time.Sleep(time.Duration(int(next.Unix()-now.Unix())) * time.Second)
}

func (t *TimeTools)SleepHour(hour int){
	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location()).Add(time.Hour*time.Duration(hour))
	time.Sleep(time.Duration(int(next.Unix()-now.Unix())) * time.Second)
}

func (t *TimeTools)SleepToNextDay(){
	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(time.Hour*24)
	time.Sleep(time.Duration(int(next.Unix()-now.Unix())) * time.Second)
}

func (t *TimeTools)SleepDay(day int){
	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(time.Hour*24*time.Duration(day))
	time.Sleep(time.Duration(int(next.Unix()-now.Unix())) * time.Second)
}

func (t *TimeTools)SleepParse(timeStr string)error{
	now := time.Now()
	next, err := t.ParseInLocation(timeStr)
	if err!=nil{
		return err
	}
	time.Sleep(time.Duration(int(next.Unix()-now.Unix())) * time.Second)
	return nil
}

func (t *TimeTools)ParseInLocation(timeStr string)(time.Time,error){
	next, err := time.ParseInLocation(layout, timeStr, time.Local)
	return next,err
}

func  (t *TimeTools)TimeParseToString(in time.Time)string{
	return in.Format(layout)
}

func  (t *TimeTools)TimeParseToTime(in time.Time)(time.Time,error){
	return t.ParseInLocation(in.Format(layout))
}

package carbon

import (
	"errors"
	"reflect"
	"time"

	"github.com/golang-module/carbon"
)

type CarbonPeriod struct {
	startDatetime *carbon.Carbon
	endDatetime   *carbon.Carbon

	isDefaultInterval bool

	recurrences int
	options     int
}

func CreateCarbonPeriod() (p *CarbonPeriod) {
	startDatetime := GetCarbonNow()
	endDatetime := startDatetime.AddDay()
	p = &CarbonPeriod{
		&startDatetime,
		&endDatetime,
		true,
		0,
		0,
	}
	// xin-da-fmt.Printf("%+v \r\n", p)
	return p
}

func CreateCarbonPeriodWithCarbon(start *carbon.Carbon, end *carbon.Carbon) (p *CarbonPeriod) {
	p = CreateCarbonPeriod()
	p.startDatetime = start
	p.endDatetime = end

	return p
}

func CreateCarbonPeriodWithTime(start time.Time, end time.Time) (p *CarbonPeriod) {
	startDate := carbon.Time2Carbon(start)
	endDate := carbon.Time2Carbon(end)

	p = CreateCarbonPeriod()
	p.startDatetime = &startDate
	p.endDatetime = &endDate

	return p
}

func CreateCarbonPeriodWithString(start string, end string, format string) (p *CarbonPeriod) {
	if format == "" {
		format = DATETIME_FORMAT
	}

	startDate := carbon.ParseByFormat(start, format)
	endDate := carbon.ParseByFormat(end, format)

	p = CreateCarbonPeriod()
	p.startDatetime = &startDate
	p.endDatetime = &endDate

	return p
}

func (period *CarbonPeriod) SetStartDate(date interface{}, inclusive interface{}) *CarbonPeriod {
	// xin-da-fmt.Println("set start datetime")
	setDate(&period.startDatetime, date)
	return period
}

func (period *CarbonPeriod) SetEndDate(date interface{}, inclusive interface{}) *CarbonPeriod {
	// xin-da-fmt.Println("set end datetime")
	setDate(&period.endDatetime, date)

	return period
}

func setDate(toSetDate **carbon.Carbon, date interface{}) (err error) {
	dType := reflect.TypeOf(date).String()
	// xin-da-fmt.Printf("%v\r\n", dType)
	// 解析字符串
	if dType == "string" {
		parsedDate := carbon.Parse(date.(string))
		if parsedDate.Error == nil {
			*toSetDate = &parsedDate
		} else {
			err = errors.New("Invalid date string xin-da-fmt.")
			return err
		}

	} else if dType == "carbon.Carbon" {
		// 直接赋值carbon指针
		ptr := date.(carbon.Carbon)
		*toSetDate = &ptr
	} else if dType == "*carbon.Carbon" {
		// 直接赋值carbon指针
		*toSetDate = date.(*carbon.Carbon)
	} else {
		// 如果不是string或者*carbon.Carbon， 抛出panic
		err = errors.New("Invalid date.")
	}

	return nil
}

func (period *CarbonPeriod) Overlaps(insideRange *CarbonPeriod) bool {
	// xin-da-fmt.Printf("start is : %#v", period.startDatetime.ToDateTimeString())
	// xin-da-fmt.Printf("current start :%s %d\r\n", period.startDatetime.ToString(), period.calculateStart())
	// xin-da-fmt.Printf("current end   :%s %d\r\n", period.endDatetime.ToString(), period.calculateEnd())
	// xin-da-fmt.Printf("range start   :%s %d\r\n", insideRange.startDatetime.ToString(), insideRange.calculateStart())
	// xin-da-fmt.Printf("range end     :%s %d\r\n\n", insideRange.endDatetime.ToString(), insideRange.calculateEnd())

	return period.calculateEnd() > insideRange.calculateStart() && insideRange.calculateEnd() > period.calculateStart()
}

func (period *CarbonPeriod) calculateStart() int64 {
	return period.startDatetime.Timestamp()
}

func (period *CarbonPeriod) calculateEnd() int64 {
	return period.endDatetime.Timestamp()
}

func (period *CarbonPeriod) DiffInDays() int64 {
	diffDays := period.startDatetime.DiffInDays(*period.endDatetime)

	return diffDays
}

func (period *CarbonPeriod) IsDiffInDays(inDays int64) bool {
	diffDays := period.startDatetime.DiffInDaysWithAbs(*period.endDatetime)

	return diffDays <= inDays
}

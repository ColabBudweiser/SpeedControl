package models

import (
	"fmt"
	v "github.com/go-ozzo/ozzo-validation"
	"regexp"
	"strconv"
	"strings"
)

const (
	datetimeRegExp      = "^[0-9]{2}\\.[0-9]{2}\\.[0-9]{4}\\s[0-9]{2}:[0-9]{2}:[0-9]{2}$"
	dateRegExp          = "^[0-9]{2}\\.[0-9]{2}\\.[0-9]{4}$"
	vehicleNumberRegExp = "^[0-9]{4}\\s[A-Z]{2}-[0-9]$"
	speedValueRegExp    = "^\\d+(,\\d{1,2})?$"
)

type ValidationType interface {
	Validate() error
}

type SpeedControlMsg struct {
	Day               string
	Month             string
	Year              string
	ParsedSpeed       int
	ParsedSpeedString string
	StringFormat      string
}

func NewSpeedControlMsg(time string, number string, speed string) *SpeedControlMsg {
	parsedSpeed, _ := strconv.ParseFloat(strings.Replace(speed, ",", ".", 1), 64)
	parsedSpeed *= 100
	parsedSpeedString := strconv.Itoa(int(parsedSpeed))

	return &SpeedControlMsg{
		Day:               strings.Split(strings.Split(time, " ")[0], ".")[0], 	// format 20
		Month:             strings.Split(strings.Split(time, " ")[0], ".")[1], 	// format 12
		Year:              strings.Split(strings.Split(time, " ")[0], ".")[2], 	// format 2019
		ParsedSpeed:       int(parsedSpeed),                                   				// format 6500
		ParsedSpeedString: parsedSpeedString,                                  				// format "6500"
		StringFormat:      fmt.Sprintf("%s %s %s", time, number, speed),       				// format 20.12.2019 14:31:25 1234 PP-7 65,5
	}
}

type ByDateMsg struct {
	Day               string
	Month             string
	Year              string
	ParsedSpeedString string
	ParsedSpeed       int
}

func NewByDateQuery(date string, speed string) *ByDateMsg {
	parsedSpeed, _ := strconv.ParseFloat(strings.Replace(speed, ",", ".", 1), 64)
	parsedSpeed *= 100
	parsedSpeedString := strconv.Itoa(int(parsedSpeed))

	return &ByDateMsg{
		Day:               strings.Split(date, ".")[0], 	// format 20
		Month:             strings.Split(date, ".")[1], 	// format 12
		Year:              strings.Split(date, ".")[2], 	// format 2019
		ParsedSpeed:       int(parsedSpeed),            		// format 6500
		ParsedSpeedString: parsedSpeedString,           		// format "6500"
	}
}

type ExtremesMsg struct {
	Day   string
	Month string
	Year  string
}

func NewExtremesMessage(date string) *ExtremesMsg {
	return &ExtremesMsg{
		Day:   strings.Split(date, ".")[0], // format 20
		Month: strings.Split(date, ".")[1], // format 12
		Year:  strings.Split(date, ".")[2], // format 2019
	}
}

type SpeedInfoInputQuery struct {
	DateTime string
	Number   string
	Speed    string
}

func (i SpeedInfoInputQuery) Validate() error {
	return v.ValidateStruct(&i,
		v.Field(&i.DateTime, v.Required, v.Match(regexp.MustCompile(datetimeRegExp))),
		v.Field(&i.Number, v.Required, v.Match(regexp.MustCompile(vehicleNumberRegExp))),
		v.Field(&i.Speed, v.Required, v.Match(regexp.MustCompile(speedValueRegExp))),
	)
}

type InfoByDateQuery struct {
	Date  string
	Speed string
}

func NewInfoByDateQuery(d string, s string) *InfoByDateQuery {
	return &InfoByDateQuery{Date: d, Speed: s}
}

func (i InfoByDateQuery) Validate() error {
	return v.ValidateStruct(&i,
		v.Field(&i.Date, v.Required, v.Match(regexp.MustCompile(dateRegExp))),
		v.Field(&i.Speed, v.Required, v.Match(regexp.MustCompile(speedValueRegExp))),
	)
}

type ExtremesByDateQuery struct {
	Date string
}

func NewExtremesByDateQuery(d string) *ExtremesByDateQuery {
	return &ExtremesByDateQuery{Date: d}
}

func (i ExtremesByDateQuery) Validate() error {
	return v.ValidateStruct(&i,
		v.Field(&i.Date, v.Required, v.Match(regexp.MustCompile(dateRegExp))),
	)
}

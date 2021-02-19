package config

import (
	"encoding/json"
	"errors"
	v "github.com/go-ozzo/ozzo-validation"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const confPath = "conf.json"
const timeRegExp = "^(0[0-9]|1[0-9]|2[0-3]|[0-9]):[0-5][0-9]$"

var Conf *Config

func init() {
	if err := initConfiguration(confPath); err != nil {
		log.Fatal("Couldn't load config: ", err)
	}
}

type Config struct {
	Host            string `json:"host"`
	ListenPort      string `json:"listen_port"`
	QueryStartTime  string `json:"select_queries_start_time"`
	QueryFinishTime string `json:"select_queries_finish_time"`
	DirsCacheSize   int    `json:"dirs_cache_size"`
	Start           Time
	Finish          Time
}

type Time struct {
	Hour   int
	Minute int
}

func (c *Config) extractTimeFromStrings() error {
	c.Start.extractTime(c.QueryStartTime)
	c.Finish.extractTime(c.QueryFinishTime)

	return CheckStartLessThanFinish(c.Start, c.Finish)
}

func CheckStartLessThanFinish(start Time, finish Time) error {
	const errorText = "invalid format of start and finish time. Start time more or equal that finish time"
	if start.Hour > finish.Hour {
		return errors.New(errorText)
	} else if start.Hour == finish.Hour {
		if start.Minute >= finish.Minute {
			return errors.New(errorText)
		}
	}
	return nil
}

func (time *Time) extractTime(timeString string) {
	timeValuesArray := strings.Split(timeString, ":")
	hours, err := strconv.Atoi(timeValuesArray[0])
	if err != nil {
		log.Fatal("Couldn't load config: ", err)
	}
	minutes, err := strconv.Atoi(timeValuesArray[1])
	if err != nil {
		log.Fatal("Couldn't load config: ", err)
	}
	time.Hour = hours
	time.Minute = minutes
}

func (c Config) Validate() error {
	return v.ValidateStruct(&c,
		v.Field(&c.QueryStartTime, v.Required, v.Match(regexp.MustCompile(timeRegExp))),
		v.Field(&c.QueryFinishTime, v.Required, v.Match(regexp.MustCompile(timeRegExp))),
	)
}

func initConfiguration(file string) error {
	Conf = &Config{}
	configFile, err := os.Open(file)
	defer configFile.Close()

	if err != nil {
		log.Fatal("Couldn't load config: ", err)
		return err
	}

	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(Conf); err != nil {
		return err
	}

	if err = Conf.Validate(); err != nil {
		return err
	}

	if err = Conf.extractTimeFromStrings(); err != nil {
		return err
	}

	log.Println("load config:", *Conf)
	return nil
}

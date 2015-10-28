package dictutil

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
)

var (
	token_file = flag.String("token", "", "File contains token")
)

func GetToken() string {
	flag.Parse()
	if *token_file == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	token, err := ioutil.ReadFile(*token_file)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return string(token)
}

//func PrepareIndex(base string) {}
//func Check(word string) {}

func LogInit() {
	log.SetFormatter(&log.JSONFormatter{})
}

func LogInfo(e string, t int, info interface{}) {
	log.WithFields(log.Fields{
		"event": e,
		"user":  t,
	}).Info(info)
}

func LogWarn(e string, t int, info interface{}) {
	log.WithFields(log.Fields{
		"event": e,
		"user":  t,
	}).Warn(info)
}

func LogError(e string, t int, info interface{}) {
	log.WithFields(log.Fields{
		"event": e,
		"user":  t,
	}).Error(info)
}

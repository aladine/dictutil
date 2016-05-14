package dictutil

import (
	"bytes"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
)

var (
	tokenFile = flag.String("token", "", "File contains token")
)

type dictionary struct {
	base    string
	wordIdx map[string]string
}

// NewDictionary returns new dictionary
func NewDictionary(b string) *dictionary {
	return &dictionary{
		base:    b,
		wordIdx: make(map[string]string),
	}
}

// PrepareIndex defines a function to index all the links
func (d *dictionary) PrepareIndex() {
	// read  dictionary
	idxData, err := ioutil.ReadFile(d.base + ".idx")
	if err != nil {
		fmt.Printf("Failed to open the dictionary: %s\n", err.Error())
	}

	dictData, err := os.Open(d.base + ".dict")
	if err != nil {
		fmt.Printf("Failed to open the dictionary: %s\n", err.Error())
	}

	// close fi on exit and check for its returned error
	defer func() {
		if err := dictData.Close(); err != nil {
			panic(err)
		}
	}()

	reader := bytes.NewBuffer(idxData)

	for {
		word, err := reader.ReadString('\x00')
		word = strings.TrimRight(word, "\x00")

		if err != nil {
			break
		}

		// offset
		offset, err := GetNumber(reader)
		if err != nil {
			break
		}

		// length
		length, err := GetNumber(reader)
		if err != nil {
			break
		}

		// desc
		desc := make([]byte, length)
		dictData.ReadAt(desc, int64(offset))
		d.wordIdx[word] = string(desc)
	}
}

// Check defines function to check whether word existed
func (d *dictionary) Check(word string) (string, error) {
	desc, exists := d.wordIdx[word]
	if !exists {
		return "", errors.New("notFound")
	}
	return desc, nil
}

// Change base of folder
func (d *dictionary) ChangeBase(base string) {
	d.base = base
}

// GetNumber is function to get Number
func GetNumber(b *bytes.Buffer) (int32, error) {
	var length int32
	bLength := make([]byte, 4)
	if n, err := b.Read(bLength); err != nil || n != 4 {
		fmt.Println("length err", n, err)
		return 0, errors.New("length err")
	}
	binary.Read(bytes.NewBuffer(bLength), binary.BigEndian,
		&length)
	return length, nil
}

// GetToken is function to get token
func GetToken() string {
	flag.Parse()
	if *tokenFile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	token, err := ioutil.ReadFile(*tokenFile)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return string(token)
}

// LogInit init logging
func LogInit() {
	log.SetFormatter(&log.JSONFormatter{})
}

// LogInfo defines logging info
func LogInfo(e string, t int64, info interface{}) {
	log.WithFields(log.Fields{
		"event": e,
		"user":  t,
	}).Info(info)
}

// LogWarn defines logging warning
func LogWarn(e string, t int64, info interface{}) {
	log.WithFields(log.Fields{
		"event": e,
		"user":  t,
	}).Warn(info)
}

// LogError defines logging error
func LogError(e string, t int64, info interface{}) {
	log.WithFields(log.Fields{
		"event": e,
		"user":  t,
	}).Error(info)
}

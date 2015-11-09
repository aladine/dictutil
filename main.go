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
	token_file = flag.String("token", "", "File contains token")
)

type dictionary struct {
	base     string
	word_idx map[string]string
}

func NewDictionary(b string) *dictionary {
	return &dictionary{
		base:     b,
		word_idx: make(map[string]string),
	}
}

func (d *dictionary) PrepareIndex() {
	// read  dictionary
	idx_data, err := ioutil.ReadFile(d.base + ".idx")
	if err != nil {
		fmt.Printf("Failed to open the dictionary: %s\n", err.Error())
	}

	dict_data, err := os.Open(d.base + ".dict")
	if err != nil {
		fmt.Printf("Failed to open the dictionary: %s\n", err.Error())
	}

	// close fi on exit and check for its returned error
	defer func() {
		if err := dict_data.Close(); err != nil {
			panic(err)
		}
	}()

	reader := bytes.NewBuffer(idx_data)

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
		dict_data.ReadAt(desc, int64(offset))
		d.word_idx[word] = string(desc)
	}
}

func (d *dictionary) Check(word string) (string, error) {
	desc, exists := d.word_idx[word]
	if !exists {
		return "", errors.New("notFound")
	}
	return desc, nil
}

func (d *dictionary) ChangeBase(base string) {
	d.base = base
}

func GetNumber(b *bytes.Buffer) (int32, error) {
	var length int32
	b_length := make([]byte, 4)
	if n, err := b.Read(b_length); err != nil || n != 4 {
		fmt.Println("length err", n, err)
		return 0, errors.New("length err")
	}
	binary.Read(bytes.NewBuffer(b_length), binary.BigEndian,
		&length)
	return length, nil
}

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

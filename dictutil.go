package dictutil

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

var (
	tokenFile = flag.String("token", "", "File contains token")
	words     = `{ "words": [
        {
            "name": "bye" ,
            "request": ["goodbye","bye bye","bye now","bye","later","laters","see you later","see you","see ya later","see ya","cya","au revoir","good night","good day"],
            "response": ["Have a good time!", "You are so nice!"]
        },
        {
            "name": "hi" ,
            "request": ["hello","hi","hola","howdy","hiya","hey","heya","yello","aloha","hi there","hey there"],
            "response": ["Hello to you too", "You are so nice!", "Hello"]
        },
        {
            "name": "thank" ,
            "request": ["okay","ok","thanks","thank you","ok thanks","ok thank you","okay thanks","okay thank you","all right","alright","great","cool","aww","awww","oh","oh okay","oh ok","ah","ah okay","ah ok","got it","gotcha"],
            "response": ["You are welcome!", "You are so nice!"]
        },
        {
            "name": "yes",
            "request": ["yes", "yep", "yup", "yeah", "yeh", "y", "that's right", "sure", "yes thanks", "yes thank you", "yes please", "yah", "ya", "for sure", "fo shizzle", "fo shiz", "yeppers", "you betcha", "you bet", "you bet ya", "certainly"],
            "response": ["You are right!"]
        },
        {
            "name": "no",
            "request": ["no","n","nope","no way","not really","nah","no thanks","no thank you"],
            "response": ["Hm...", "Okie, fine"]
        },
        {
            "name": "haha",
            "request": ["haha", "ha", "hehe", "lol", "rofl", "lmao"],
            "response": ["You are so funny", "You make me laugh"]
        }
    ] 
}`
	wordsData map[string][]wordObject
)

type wordObject struct {
	Name     string   `json:"name"`
	Request  []string `json:"request"`
	Response []string `json:"response"`
}

type wordList struct {
	Words []wordObject `json:"words"`
}

type dictionary struct {
	base    string
	wordIdx map[string]string
}

func init() {
	err := json.Unmarshal([]byte(words), &wordsData)
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixNano())
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

// IsGreetingMsg defines whether a string is greeting type msg
func IsGreetingMsg(str string) (result string, isGreeting bool) {
	for _, w := range wordsData["words"] {
		for _, k := range w.Request {
			if k == str {
				i := rand.Intn(len(w.Response))
				return w.Response[i], true
			}
		}
	}
	return "", false
}

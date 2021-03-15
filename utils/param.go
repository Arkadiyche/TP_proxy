package utils

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GetParams() []string {
	var Params = make([]string, 0, 0)
	inputFile, err := os.Open("params")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		if scanner.Err() != nil {
			log.Fatal(scanner.Err().Error())
		}
		Params = append(Params, scanner.Text())
	}
	return Params
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes() string {
	b := make([]rune, rand.Intn(20))
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

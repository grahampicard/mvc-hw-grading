package main

import (
	"errors"
	"log"
	"net/url"
	"os"

	"github.com/yale-mgt-656-fall-2018/mvc-hw-grading/questions"
)

func isValidURL(testURL string) bool {
	_, err := url.ParseRequestURI(testURL)
	return err == nil
}

func getArgs(args []string) (string, string, error) {
	if len(args) != 3 {
		return "", "", errors.New("Error: you must provide exactly two arguments: class nickname and app URL")
	}
	theNickname := args[1]
	theURL := args[2]
	if isValidURL(theURL) == false {
		return "", "", errors.New("Error: you provided an invalid URL")
	}

	return theNickname, theURL, nil
}

func main() {
	userNickname, userURL, err := getArgs(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	questions.TestAll(userNickname, userURL, true)
}

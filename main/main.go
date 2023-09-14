package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"regexp"
	"time"
)

func isEmailValid(email string) bool {
	if len(email) <= 0 {
		return false
	}
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}
func isValidCategory(category string) bool {
	var categories [5]string = [5]string{"new", "top", "hot", "rising", "controversial"}
	for i := 1; i < 5; i++ {
		if categories[i] == category {
			return true
		}
	}
	return false
}
func validateCMDFlags(
	subredditName string,
	exportPath string,
	email string,
	category string,
	exportToFile bool,
	limit int,
	sendEmail bool,
) error {
	if exportToFile && sendEmail {
		return errors.New("you can either export the results to local filesystem or receive them as email")
	}
	if res := isValidCategory(category); !res {
		return errors.New("category should be one of the following : new, top, hot, controversial or rising")
	}
	if sendEmail && (!isEmailValid(email)) {
		return errors.New("please enter a valid email address")
	}
	if limit > 100 || limit <= 0 {
		return errors.New("please set the value of limit, ranging from 1 to 100")
	}
	return nil
}
func main() {
	var subredditName string
	var exportPath string
	var email string
	var category string
	var exportToFile bool
	var limit int
	var sendEmail bool

	flag.StringVar(&subredditName, "subreddit", "AskReddit", "Name of the subreddit.")
	flag.StringVar(&exportPath, "export_path", fmt.Sprintf("../exports/%v", time.Now().Format("2006-01-02 - 15:04:05")), "Path of file to write data to.")
	flag.StringVar(&email, "email", "", "Receive zipped data to your email.")
	flag.StringVar(&category, "category", "hot", "Sort content by category.")
	flag.BoolVar(&exportToFile, "export_to_file", true, "Name of the subreddit.")
	flag.IntVar(&limit, "limit", 1, "No of posts to retrieve.")
	flag.BoolVar(&sendEmail, "send_email", false, "Should data be sent through email")

	flag.Parse()
	if err := validateCMDFlags(
		subredditName,
		exportPath,
		email,
		category,
		exportToFile,
		limit,
		sendEmail,
	); err != nil {
		log.Fatalln(err.Error())
		return
	}
	fmt.Println(email)
}

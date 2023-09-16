package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	// "io"
	"log"
	"net/http"
	"regexp"
	"time"
)

type ExportData struct {
	*ResponseData
}

type ResponseData struct {
	Data struct {
		Name        string `json:"display_name"`
		Description string `json:"public_description"`
		HeaderTitle string `json:"header_title"`
		Subscribers int    `json:"subscribers"`
		UsersCount  int64  `json:"active_user_count"`
		ImageUrl    string `json:"icon_img"`
	} `json:"data"`
}

type SubredditValidatorStruct struct {
	Data struct {
		NSubs      int    `json:"subscribers"`
		TypeOfSubR string `json:"subreddit_type"`
	} `json:"data"`
}

func HandleError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}
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
	filename string,
	email string,
	category string,
	exportToFile bool,
	limit int,
	sendEmail bool,
) error {
	if exportToFile && sendEmail {
		sendEmail = true
		exportToFile = false
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

func validateSubreddit(subreddit string) bool {
	subredditUrl := fmt.Sprintf("http://www.reddit.com/r/%s/about.json", subreddit)
	client := &http.Client{}
	var responseData *SubredditValidatorStruct
	req, err := http.NewRequest("GET", subredditUrl, nil)
	HandleError(err)
	res, err := client.Do(req)
	HandleError(err)
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&responseData)
	HandleError(err)
	fmt.Printf("Data %v\n", responseData)
	if responseData != nil && responseData.Data.NSubs != 0 && responseData.Data.TypeOfSubR == "public" {
		return true
	}

	return false
}

func getSubRedditInfo(subreddit string) ResponseData {
	subredditUrl := fmt.Sprintf("http://www.reddit.com/r/%s/about.json", subreddit)
	client := &http.Client{}
	var responseData *ResponseData
	req, err := http.NewRequest("GET", subredditUrl, nil)
	HandleError(err)
	res, err := client.Do(req)
	HandleError(err)
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&responseData)
	HandleError(err)
	return *responseData
}

func saveDataToFile(path, filename string, permission fs.FileMode) *os.File {

	// fmt.Printf("path : %v\n", path)
	// fmt.Printf("filename : %v\n", filename)

	// fmt.Printf("File path : %v", filepath.Join(".", path, filename))
	HandleError(os.Mkdir(path, permission))
	fmt.Println("here")
	file, err := os.Create(filepath.Join(path, filename))
	HandleError(err)
	return file
}
func main() {

	var subredditName string
	var filename string
	var email string
	var category string
	var exportToFile bool
	var limit int
	var sendEmail bool

	flag.StringVar(&subredditName, "subreddit", "", "Name of the subreddit.")
	flag.StringVar(&filename, "filename", "", "Path of file to write data to.")
	flag.StringVar(&email, "email", "", "Receive zipped data to your email.")
	flag.StringVar(&category, "category", "hot", "Sort content by category.")
	flag.BoolVar(&exportToFile, "export_to_file", true, "Name of the subreddit.")
	flag.IntVar(&limit, "limit", 1, "No of posts to retrieve.")
	flag.BoolVar(&sendEmail, "send_email", false, "Should data be sent through email")

	flag.Parse()
	filename = fmt.Sprintf("%v.json", filename)
	if err := validateCMDFlags(
		subredditName,
		filename,
		email,
		category,
		exportToFile,
		limit,
		sendEmail,
	); err != nil {
		log.Fatalln(err.Error())
		return
	}
	dir, err := os.Getwd()
	HandleError(err)

	var exportPath string = fmt.Sprintf("%v\\exports\\%v - %v", dir, subredditName, time.Now().Format("2006-01-02 15-04-05 PM"))

	validSubR := validateSubreddit(subredditName)
	if !validSubR {
		HandleError(errors.New("the provided subreddit is either invalid or private"))
	}
	subredditInfo := getSubRedditInfo(subredditName)
	fmt.Println(subredditInfo)

	// get subreddit posts
	//get subreddit post comments
	exportData := &ExportData{
		&subredditInfo,
	}
	jsonData, err := json.MarshalIndent(exportData, "", "  ")
	HandleError(err)
	if exportToFile && filename != "" {
		file := saveDataToFile(exportPath, filename, 0777)
		defer file.Close()

		n, err := io.WriteString(file, string(jsonData))
		HandleError(err)
		fmt.Printf("Saved data for subreddit %s to path %s\nBytes: %d\n", subredditName, exportPath, n)
	}
	if sendEmail && email != "" {
		// archive data
		// send data to the email
	}
}

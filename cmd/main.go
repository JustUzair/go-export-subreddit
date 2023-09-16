package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"
)

var wg sync.WaitGroup

type ExportData struct {
	*SubredditAboutData `json:"about"`
	*SubredditPostsData `json:"posts"`
}

type SubredditAboutData struct {
	Data struct {
		Name         string `json:"display_name"`
		Description  string `json:"public_description"`
		HeaderTitle  string `json:"header_title"`
		Subscribers  int    `json:"subscribers"`
		UsersCount   int64  `json:"active_user_count"`
		ImageUrl     string `json:"icon_img"`
		HeaderImgUrl string `json:"header_img"`
	} `json:"data"`
}

type SubredditPostsData struct {
	Data struct {
		Children []struct {
			Data struct {
				PostId        string        `json:"id"`
				SelfText      string        `json:"selftext"`
				Author        string        `json:"author"`
				Title         string        `json:"title"`
				SubredditName string        `json:"subreddit_name_prefixed"`
				UpVotes       int           `json:"ups"`
				Created       json.Number   `json:"created"`
				NComments     int           `json:"num_comments"`
				PostUrl       string        `json:"url"`
				Permalink     string        `json:"permalink"`
				Comments      []interface{} `json:"comments"`
			} `json:"data"`
		} `json:"children"`
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
	var categories [5]string = [5]string{"new", "top", "hot", "rising"}
	for i := 0; i < 5; i++ {
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
	exportPosts bool,
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

func getSubRedditInfo(subreddit string) *SubredditAboutData {
	subredditUrl := fmt.Sprintf("http://www.reddit.com/r/%s/about.json", subreddit)
	client := &http.Client{}
	var responseData *SubredditAboutData
	req, err := http.NewRequest("GET", subredditUrl, nil)
	HandleError(err)

	req.Header.Set("authority", "www.reddit.com")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="89", "Chromium";v="89", ";Not A Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("dnt", "1")
	req.Header.Set("sec-fetch-site", "none")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("accept-language", "en-GB,en;q=0.9")
	res, err := client.Do(req)
	HandleError(err)
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&responseData)
	HandleError(err)
	return responseData
}

func getSubRedditPosts(subreddit, category string, limit int, exportComments bool, exportPath string) *SubredditPostsData {
	subredditUrl := fmt.Sprintf("http://www.reddit.com/r/%s/%s.json?limit=%d", subreddit, category, limit)
	client := &http.Client{}
	var responseData *SubredditPostsData

	req, err := http.NewRequest("GET", subredditUrl, nil)
	req.Header.Set("authority", "www.reddit.com")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="89", "Chromium";v="89", ";Not A Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("dnt", "1")
	req.Header.Set("sec-fetch-site", "none")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("accept-language", "en-GB,en;q=0.9")
	HandleError(err)
	res, err := client.Do(req)
	HandleError(err)
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&responseData)
	HandleError(err)
	if exportComments {
		for i := range responseData.Data.Children {
			wg.Add(1)
			i := i
			go func() {
				defer wg.Done()
				permalink := responseData.Data.Children[i].Data.Permalink
				commentsURL := fmt.Sprintf("http://reddit.com%v.json", permalink)
				log.Println(commentsURL)
				req, err := http.NewRequest("GET", commentsURL, nil)
				HandleError(err)
				res, err := client.Do(req)
				HandleError(err)
				defer res.Body.Close()
				// var comments []interface{}
				var comments []struct {
					Data struct {
						Children []struct {
							Data struct {
								Permalink  string `json:"permalink"`
								Body       string `json:"body"`
								PostId     string `json:"id"`
								CommentUrl string `json:"url"`
							}
						} `json:"children"`
					} `json:"data"`
				}
				err = json.NewDecoder(res.Body).Decode(&comments)
				responseData.Data.Children[i].Data.Comments = append(responseData.Data.Children[i].Data.Comments, comments)
				HandleError(err)

				// ------------- Logic to export comments to individual files -------------
				// file := saveDataToFile(exportPath, fmt.Sprintf("post-author-%s comments - %d.json", responseData.Data.Children[i].Data.Author, i+1), 0777)
				// commentJSONData, err := json.MarshalIndent(comments, "", "  ")
				// HandleError(err)

				// n, err := io.WriteString(file, string(commentJSONData))
				// HandleError(err)
				// fmt.Printf("Saved data for subreddit %s's comment-%d to path %s\nBytes: %d\n", subreddit, i+1, exportPath, n)
				// ------------- --------------------------------------------- -------------

			}()

			wg.Wait()
		}

	}
	return responseData
}

func saveDataToFile(path, filename string, permission fs.FileMode) *os.File {
	os.MkdirAll(path, permission)
	file, err := os.Create(filepath.Join(path, filename))
	HandleError(err)
	return file
}
func main() {
	var subredditPostsInfo *SubredditPostsData

	var subredditName string
	var filename string
	var email string
	var category string
	var exportToFile bool
	var limit int
	var sendEmail bool
	var exportPosts bool
	var exportComments bool
	flag.StringVar(&subredditName, "subreddit", "", "Name of the subreddit.")
	flag.StringVar(&filename, "filename", "", "Path of file to write data to.")
	flag.StringVar(&email, "email", "", "Receive zipped data to your email.")
	flag.StringVar(&category, "category", "hot", "Sort content by category.")
	flag.BoolVar(&exportToFile, "export_to_file", true, "Name of the subreddit.")
	flag.IntVar(&limit, "posts_limit", 1, "No of posts to retrieve.")
	flag.BoolVar(&sendEmail, "send_email", false, "Should data be sent through email")
	flag.BoolVar(&exportPosts, "export_posts", false, "Should posts data be fetched")
	flag.BoolVar(&exportComments, "export_comments", false, "Should comments data be fetched")

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
		exportPosts,
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
	if exportPosts {
		subredditPostsInfo = getSubRedditPosts(subredditName, category, limit, exportComments, exportPath)
	}

	exportData := &ExportData{
		subredditInfo,
		subredditPostsInfo,
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


# Go-Export Subreddit

The project is a Command-Line-Application aimed at selectively fetching the reddit data about a subreddit (incl. posts, comments, data about the subreddit itself, etc).

## Usage/Examples
### Fetch the repository and install the required dependencies
```bash
   git clone  https://github.com/JustUzair/go-export-subreddit.git
   cd go-export-subreddit
   go mod tidy
```

### Setting Up the ```config.env``` file
- Please note that these flags are necessary only when you want to use the email functionality to send the fetched data to your email address
    - ```SENDER_NAME``` - Name of the sender of email
    - ```SENDER_EMAIL``` - Email address of the sender
    - ```SENDER_PASSWORD``` - Application password of the email sender
        - Generate password for your app [here](https://myaccount.google.com/apppasswords).

## Usaing the CLI Application

#### CLI Flags and their Usage


##### Syntax of the command 
```bash
    go run -v ./cmd --subreddit=NAME_OF_SUBREDDIT --filename=NAME_OF_FILE_WITHOUT_EXTENSION --export_comments=[true | false] --export_comments=[true | false] --export_to_file=[true | false]  --posts_limit=[1 - 100]  --email=YOUR_EMAIL_ADDRESS --send_email=[true | false]
```


##### Example command 

```bash
go run -v ./cmd --filename=Meh --export_to_file=true --subreddit=dadjokes --category=rising  --export_posts=true --export_comments=true --posts_limit=10  --send_email=true --email=uzairhajra76330@gmail.com
```




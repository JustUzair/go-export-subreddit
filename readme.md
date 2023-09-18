```bash
    go run -v ./cmd --subreddit=NAME_OF_SUBREDDIT --filename=NAME_OF_FILE_WITHOUT_EXTENSION --export_comments=[true | false] --export_comments=[true | false] --export_to_file=[true | false]  --posts_limit=[1 - 100]  --email=YOUR_EMAIL_ADDRESS --send_email=[true | false]
```


## Example command 

```bash
go run -v ./cmd --filename=Meh --export_to_file=true --subreddit=dadjokes --category=rising  --export_posts=true --export_comments=true --posts_limit=10  --send_email=true --email=uzairhajra76330@gmail.com
```
# Web Crawler
You can use this crawler to get all the links contained in a web page.

You can either use the crawler.go, or crawler which I explain each below:

* To use crawler.go, you need to have "go" already installed, then you can open up a terminal and run the file using this command:
```
go run crawler.go https://google.com 
    (or any other website)
```
* To use crawler file, you have the advantage of using the crawler without "go" being installed on your computer, so you can use this command:
```
./crawler https://google.com
```
## How does it work?
It basically (as you can see in the source code), looks for the url given by the user, and downloads the content recursively,
then it looks for the "href" in the content of the page (which you can change it yourself) and returns the result in output.

If there is no error in the process, You will get "No failures" and "Done" at the end.

The timeout is also set to one minute, meaning that if one minute is reached and crawler is still looking for content in that website,
it just finishes the program, which also, you can modify in the source code and build it yourself using:
```
go build crawler.go
```
package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	Id           string `json:"id"`
	OringinalUrl string `json:"originalUrl"`
	ShortUrl     string `json:"shorturl"`
	OriginalDate time.Time `json:"originaltime"`
}
/*
      id---->{
	          id
			  originalurl
			  shorturl
			  originadldate
	         }
	 // in the above format the data strores according to the id it will sort in database it will redirect using map
*/

var urlDB = make(map[string]URL)

func generateShortUrl(OriginalUrl string) string{
	hasher:=md5.New() //we need to convert the original url to hash values
	hasher.Write([]byte(OriginalUrl)) //then hash to bytes
	// fmt.Print(hasher)
	data:=hasher.Sum(nil)//sum the bytes
	// fmt.Println(data)
	hash:=hex.EncodeToString(data)//then covert bytes to the string
	// fmt.Println(hash)
	fmt.Println(hash[:8])//then retuen the only 8 character we can adjuct how many char we need
	return hash[:8]
}

func createUrl(originalUrl string) string{
	shortUrl:=generateShortUrl(originalUrl)
	id:=shortUrl
	urlDB[id]=URL{
		Id: id,
		OringinalUrl: originalUrl,
		ShortUrl: shortUrl,
		OriginalDate: time.Now(),
	}
	return shortUrl
}

func getUrl(id string)(URL ,error){
	url,ok:=urlDB[id]
	if !ok{
		return URL{},errors.New("url not found")
	}
	return url,nil
}

func handler(w http.ResponseWriter,r *http.Request){
	// fmt.Println("GET Method")//it will be printed in terminal that the server is working
	fmt.Fprintln(w,"Hello world")//it will be printed that response writer where the server is running screen in writer
}

func shortUrlHandler(w http.ResponseWriter, r *http.Request){
	var data struct{
		URL string `json:"url"`
	}
	er:=json.NewDecoder(r.Body).Decode(&data)
	if er!=nil{
		http.Error(w,"invalid request body",http.StatusBadRequest)
		return
	}
	shortUrl:=createUrl(data.URL)

	// fmt.Fprintln(w,shortUrl)
	response:=struct{
		ShortUrl string `json:"short_url"`
	}{ShortUrl: shortUrl}
	w.Header().Set("Content-Type","application/json")//these 2 lines for sending response
	json.NewEncoder(w).Encode(response)
}

func redirectUrlHandler(w http.ResponseWriter, r *http.Request){
	id:=r.URL.Path[len("/redirect/"):]
	url,er:=getUrl(id)
	if er!=nil{
		http.Error(w,"inavlid request",http.StatusNotFound)
		return
	}
	http.Redirect(w,r,url.OringinalUrl,http.StatusFound)
}


func main() {
	// generateShortUrl("https://github.com/GNANESH9704/")

	//register the handler function to handle the all requests in the server
	http.HandleFunc("/",handler)
	http.HandleFunc("/shorten",shortUrlHandler)
	http.HandleFunc("/redirect/",redirectUrlHandler)

	//start HTTP server on port 3000...
	fmt.Println("Starting server on port 3000....")
	er:=http.ListenAndServe(":3000",nil)
	if er!=nil{
		fmt.Println("error",er)
		return
	}
}
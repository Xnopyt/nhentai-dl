package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Xnopyt/nhentai-go"
)

type job struct {
	path string
	url  string
}

type doujinshiJSON struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	TotalPages int      `json:"pages"`
	Tags       []string `json:"tags"`
	URL        string   `json:"url"`
}

var verbose *bool
var threads *int
var activeJobs int

func init() {
	verbose = flag.Bool("verbose", false, "Enable verbose mode")
	threads = flag.Int("threads", 15, "The maximum goroutines that will be ran at once. (Excluding main.main()")
	flag.Parse()
}

func main() {
	sendVerbose("Verbose mode is enabled.")
	fmt.Println("nHentai Downloader - A tool for bulk downloads from nhentai written in Go.")
	fmt.Print("Made by Xnopyt (https://github.com/Xnopyt)\n\n")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please enter the ID's to download separated by commas: ")
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\r", "", -1)
	idstring := strings.Split(strings.Split(text, "\n")[0], ",")
	var ids []int
	for _, v := range idstring {
		idint, err := strconv.Atoi(v)
		if err != nil {
			panic(errors.New("all ids must be integers"))
		}
		ids = append(ids, idint)
	}
	var queue []*nhentai.Doujinshi
	for _, v := range ids {
		sendVerbose("Getting info for " + strconv.Itoa(v))
		info, err := nhentai.Get(v)
		if err != nil {
			fmt.Println("Couldn't find anything for " + strconv.Itoa(v))
			continue
		}
		fmt.Println("Adding \"" + info.Title + "\" to the queue.")
		queue = append(queue, info)
	}
	if len(queue) <= 0 {
		panic("There is nothing to do.")
	}
	fmt.Println("Creating jobs, please wait...")
	var jobs []job
	var jobCount int
	for _, v := range queue {
		sendVerbose(v.Title + " has " + strconv.Itoa(v.TotalPages) + " pages.")
		for i := 1; i <= v.TotalPages; i++ {
			var x job
			x.url = "https://i.nhentai.net/galleries/" + v.MediaID + "/" + strconv.Itoa(i) + "." + v.Pages[i-1].Ext
			x.path = v.ID + "/" + strconv.Itoa(i) + "." + v.Pages[i-1].Ext
			jobs = append(jobs, x)
			jobCount++
			sendVerbose("Created job for " + x.url)
		}
		sendVerbose("Creating output folder for " + v.ID)
		os.Mkdir(v.ID, 0777)
		sendVerbose("Creating metadata file for " + v.ID)
		var meta doujinshiJSON
		meta.Title = v.Title
		meta.ID = v.ID
		meta.Tags = v.Tags
		meta.TotalPages = v.TotalPages
		meta.URL = v.URL
		metajson, _ := json.MarshalIndent(meta, "", " ")
		ioutil.WriteFile(v.ID+"/metadata.json", metajson, 0644)
	}
	fmt.Println("Created " + strconv.Itoa(jobCount) + " jobs.")
	fmt.Println(strconv.Itoa(*threads) + " goroutines will be used to download.")
	for jobCount > 0 {
		for activeJobs < *threads {
			if jobCount < 1 {
				break
			}
			go download(jobs[jobCount-1])
			jobCount--
			activeJobs++
		}
		time.Sleep(time.Second)
	}
	sendVerbose("Waiting for all routines to cleanup...")
	time.Sleep(5 * time.Second)
	fmt.Println("Download complete!")
}

func download(j job) {
	sendVerbose("Downloading: " + j.url)
	resp, err := http.Get(j.url)
	if err != nil {
		fmt.Println("Error downloading " + j.url + ": " + err.Error())
		activeJobs--
		return
	}
	defer resp.Body.Close()
	out, err := os.Create(j.path)
	if err != nil {
		fmt.Println("Error downloading " + j.url + ": " + err.Error())
		activeJobs--
		return
	}
	defer out.Close()
	io.Copy(out, resp.Body)
	sendVerbose("Finished job for " + j.url)
	activeJobs--
	time.Sleep(10 * time.Second)
}

func sendVerbose(txt string) {
	if *verbose == true {
		fmt.Println("Verbose: " + txt)
	}
}

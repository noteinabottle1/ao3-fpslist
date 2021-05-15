package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/pflag"
)

var (
	rePodsOk     bool // WARNING: setting this to false will greatly increase program execution time
	ao3FandomId  int
	maxWordCount int
	minWordCount int
	bpAuthors    []string
)

func main() {
	pflag.BoolVar(&rePodsOk, "rePodsOk", true, "check if work already has been podficced")
	pflag.IntVar(&ao3FandomId, "id", 0, "AO3 Fandom ID")
	pflag.StringSliceVar(&bpAuthors, "authors", []string{}, "Authors with BP statements (e.g. a,b,c)")
	pflag.IntVar(&maxWordCount, "maxWords", 5000, "Max words in a fic")
	pflag.IntVar(&minWordCount, "minWords", 50, "Min words in a fic")

	pflag.Parse()

	// TODO make this work
	//fpsListFandomId := "13862"
	//bpAuthors := getBPAuthorsList(fpsListFandomId)

	// Output to the following CSV file
	outputCSVFile := strconv.Itoa(ao3FandomId) + "_bplist.csv"
	os.Remove(outputCSVFile)
	f, err := os.Create(outputCSVFile)
	defer f.Close()
	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	w := csv.NewWriter(f)
	defer w.Flush()

	// CSV Header
	w.Write([]string{"Title", "Author", "Rating", "Pairing", "Archive Warnings", "Word Count", "Link"})

	for _, author := range bpAuthors {
		authorsPageLink := "https://archiveofourown.org/users/" + author + "/works?fandom_id=" + strconv.Itoa(ao3FandomId)
		authorsWorks := scrapeAuthorsPage(authorsPageLink)
		for _, work := range authorsWorks {
			if err := w.Write(work); err != nil {
				log.Fatalln("error writing record to file", err)
			}
		}
		// Don't get rate limited by Ao3
		time.Sleep(time.Second * 10)
	}
}

func scrapeAuthorsPage(page string) [][]string {
	var records [][]string
	doc := loadPage(page)
	if doc == nil {
		return records
	}

	numPages := (doc.Find(".pagination.actions").Find("li").Length() - 4) / 2
	if numPages < 0 { // Only 1 page of works
		works := parseWorks(doc)
		records = append(records, works...)
	} else {
		for pageN := 1; pageN < numPages; pageN++ {
			doc = loadPage(fmt.Sprintf("%v&page=%d", page, pageN))
			if doc != nil {
				works := parseWorks(doc)
				records = append(records, works...)
				// Don't get rate limited by Ao3
				time.Sleep(time.Second * 10)
			}
		}
	}
	return records
}

func parseWorks(doc *goquery.Document) [][]string {
	var works [][]string
	doc.Find(".work.blurb").Each(func(_ int, work *goquery.Selection) {
		name := work.Find("a").First().Text()
		if strings.Contains(strings.ToLower(name), "podfic") {
			// This is already a podfic
		} else {
			author := work.Find("a").Eq(1).Text()
			wordsString := work.Find(".words").Eq(1).Text()
			numWords, _ := strconv.Atoi(strings.ReplaceAll(wordsString, ",", ""))
			if numWords > maxWordCount || numWords < minWordCount {
				fmt.Printf("%v by %s is too many words: %s\n", name, author, wordsString)
			} else {
				fullRating := work.Find("span.rating").Find(".text").Text()
				warnings := work.Find(".warnings").Find(".tag").Text()
				relationships := work.Find(".relationships").Find(".tag").Text()
				if strings.Contains(relationships, "&") {
					relationships = "Gen"
				}
				relativeLink, _ := work.Find("a").Attr("href")
				authorLink := strings.ReplaceAll(author, " ", "%20")
				link := fmt.Sprintf("https://archiveofourown.org/users/%v%v", authorLink, relativeLink)

				// Check that it has no 'inspired works', which we'll assume to be podfics
				if rePodsOk || noPodficsYet(link) {
					summary := []string{name, author, string(fullRating[0]), relationships, warnings, wordsString, link}
					works = append(works, summary)
					fmt.Printf("%s by %s, %s words\n", name, author, wordsString)
				}
			}
		}
	})
	return works
}

func loadPage(page string) *goquery.Document {
	// Request the HTML page.
	res, err := http.Get(page)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Printf("status code error: %d %s", res.StatusCode, res.Status)
		return nil
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil
	}
	return doc
}

func noPodficsYet(ficPage string) bool {
	// Don't get rate limited by Ao3
	time.Sleep(time.Second * 10)

	res, _ := http.Get(ficPage + "#children")
	defer res.Body.Close()
	return res.StatusCode == 404
}

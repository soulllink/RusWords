package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/mattn/go-sqlite3"
)

var RusWords []string
var Stypes = []string{
	"существительное",
	"прилагательное",
	"числительное",
	"местоимение",
	"глагол",
	"наречие",
	"предлог",
	"союз",
	"частица",
	"междометие",
}

func querryStype(s string) int {
	res, _ := http.Get("https://ru.wiktionary.org/wiki/" + s)
	defer res.Body.Close()
	if strings.Contains(res.Status, "200") {
		doc, _ := goquery.NewDocumentFromReader(res.Body)
		preprocess := doc.Find("a").Text()
		preprocess = strings.ToLower(preprocess)
		resultdec := decodeStype(preprocess)
		if resultdec == -1 {
			resultdec = querryAlt(preprocess)
		}
		addtodb(s, resultdec)
		return resultdec

	} else {
		return 404
	}

}

func querryAlt(s string) int {
	res, _ := http.Get("https://pishugramotno.ru/morfologiya/" + s)
	defer res.Body.Close()
	if strings.Contains(res.Status, "200") {
		doc, _ := goquery.NewDocumentFromReader(res.Body)
		preprocess := doc.Find("h2").Text()
		preprocess = strings.ToLower(preprocess)
		resultdec := decodeStype(preprocess)
		if resultdec == -1 {
			resultdec = querryAlttwo(preprocess)
		}
		addtodb(s, resultdec)
		return resultdec

	} else {
		return 404
	}

}

func querryAlttwo(s string) int {
	res, _ := http.Get("https://rustxt.ru/morfologicheskij-razbor-slova/" + s)
	defer res.Body.Close()
	if strings.Contains(res.Status, "200") {
		doc, _ := goquery.NewDocumentFromReader(res.Body)
		preprocess := doc.Find("span").Text()
		preprocess = strings.ToLower(preprocess)
		return decodeStype(preprocess)

	} else {
		return 404
	}

}

func decodeStype(s string) int {
	for i, v := range Stypes {
		if strings.Contains(s, v) {
			return i
		}
	}
	return -1
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func sqldb() *sql.DB {
	//init
	database, _ :=
		sql.Open("sqlite3", "words.db")

	statement, _ :=
		database.Prepare("CREATE TABLE IF NOT EXISTS ruswords (id INTEGER PRIMARY KEY, words TEXT, stype INTEGER)")

	statement.Exec()
	return database
}

func addtodb(word string, stype int) {
	statement, _ :=
		sqldb().Prepare("INSERT INTO ruswords (words, stype) VALUES (?, ?)")

	statement.Exec(word, stype)
}

func main() {
	//init
	sqldb()

	file, _ := os.Open("Gugo.txt")
	scanner := bufio.NewScanner(file)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	file.Close()

	for _, eachline := range txtlines {
		s := strings.ToLower(eachline)
		words := strings.Replace(s, "\n", " ", -1)
		words = strings.Replace(words, "!", " ", -1)
		words = strings.Replace(words, ",", " ", -1)
		words = strings.Replace(words, "?", " ", -1)
		words = strings.Replace(words, ".", " ", -1)
		words = strings.Replace(words, "(", " ", -1)
		words = strings.Replace(words, ")", " ", -1)
		words = strings.Replace(words, "«", " ", -1)
		words = strings.Replace(words, "»", " ", -1)
		words = strings.Replace(words, "…", " ", -1)
		out := strings.Split(words, " ")
		for _, v := range out {
			if contains(RusWords, v) == false {
				RusWords = append(RusWords, v)
			}
		}

	}

	fmt.Println("Done s1: ", len(RusWords))
	for _, word := range RusWords {
		fmt.Println("Doing: " + word)
		querryStype(word)
		fmt.Println("OK")

	}
}

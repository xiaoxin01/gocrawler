package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"time"

	"github.com/gocolly/colly"
	"github.com/spf13/viper"
)

// Web web to crawl
type Web struct {
	Enabled      bool
	URL          string
	ListSelector string
	MinFields    int
	PageCursor   *PageCursor
	Fields       map[string]Field
	Headers      map[string]string
	Visited      map[string]bool
}

// Field field to add to each item
type Field struct {
	Operator  string
	Parameter string
	Selector  string
	Regexp    *RegexOperation
}

// PageCursor visit page by identity
type PageCursor struct {
	URLFormat string
	Start     int
	End       int
}

// RegexOperation regexp to change field value
type RegexOperation struct {
	Expression string
	Group      int
}

func main() {
	rand.Seed(time.Now().UnixNano())
	viper.AddConfigPath(".")
	viper.SetConfigName("webs")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	webs := viper.GetStringMap("webs")

	for path := range webs {

		var web Web
		if err := viper.UnmarshalKey("webs."+path, &web); err != nil {
			panic(err)
		}
		if !web.Enabled {
			continue
		}

		bytes, err := ioutil.ReadFile(fmt.Sprintf("data/%s.state", path))
		web.Visited = make(map[string]bool)
		if err == nil {
			json.Unmarshal(bytes, &web.Visited)
		}
		os.Mkdir("data", 0744)
		f, err := os.OpenFile(fmt.Sprintf("data/%s.json", path), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
		defer f.Close()

		if err != nil {
			panic(err)
		}

		// Instantiate default collector
		c := colly.NewCollector()

		// Before making a request put the URL with
		// the key of "url" into the context of the request
		c.OnRequest(func(r *colly.Request) {
			for key, value := range web.Headers {
				r.Headers.Add(key, value)
			}
		})

		c.OnHTML(web.ListSelector, func(e *colly.HTMLElement) {
			obj := make(map[string]interface{}, 0)

			for key, field := range web.Fields {
				if value, ok := getValue(e, field); ok {
					obj[key] = value
				}
			}

			if len(obj) < web.MinFields {
				return
			}

			objString, _ := json.Marshal(obj)

			f.WriteString(string(objString))
			f.WriteString("\n")

			//fmt.Println(string(objString))
		})

		c.OnResponse(func(r *colly.Response) {
			//r.Body
		})

		// Start scraping on https://en.wikipedia.org
		fmt.Println("visit url: ", web.URL)
		web.Visited[web.URL] = true
		c.Visit(web.URL)

		if web.PageCursor != nil {
			for start := web.PageCursor.Start; start <= web.PageCursor.End; start++ {
				url := fmt.Sprintf(web.PageCursor.URLFormat, start)
				if _, ok := web.Visited[url]; !ok {
					time.Sleep(time.Second * time.Duration(rand.Intn(1)+1))
					web.Visited[url] = true
					fmt.Println("visit url: ", url)
					c.Visit(url)
				} else {
					fmt.Println("skip url: ", url)
				}
			}
		}
		// viper.Set("webs."+path, &web)
		// viper.WriteConfig()

		file, _ := json.MarshalIndent(web.Visited, "", " ")
		_ = ioutil.WriteFile(fmt.Sprintf("data/%s.state", path), file, 0644)
	}
}

func getValue(e *colly.HTMLElement, field Field) (v interface{}, ok bool) {
	ok = true
	switch field.Operator {
	case "Attr":
		v = e.ChildAttr(field.Selector, field.Parameter)
	case "Attrs":
		if attr := e.ChildAttrs(field.Selector, field.Parameter); attr != nil {
			v = attr
		} else {
			ok = false
		}
	case "Text":
		v = e.ChildText(field.Selector)
	case "Const":
		v = field.Parameter
	default:
		ok = false
	}

	if ok && field.Regexp != nil {
		values := regexp.MustCompile(field.Regexp.Expression).FindStringSubmatch(v.(string))
		if len(values) > field.Regexp.Group {
			v = values[field.Regexp.Group]
		}
	}

	ok = ok && v != ""

	return
}

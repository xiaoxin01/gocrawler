package main

import (
	"encoding/json"
	"fmt"

	"github.com/gocolly/colly"
	"github.com/spf13/viper"
)

// Web web to crawl
type Web struct {
	URL          string
	ListSelector string
	minFields    int
	Fields       map[string]Field
	Headers      map[string]string
}

// Field field to add to each item
type Field struct {
	Operator  string
	Parameter string
	Selector  string
}

func main() {
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

			if len(obj) < web.minFields {
				return
			}

			objString, _ := json.Marshal(obj)

			fmt.Println(string(objString))
		})

		c.OnResponse(func(r *colly.Response) {
			//r.Body
		})

		// Start scraping on https://en.wikipedia.org
		c.Visit(web.URL)
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
	default:
		ok = false
	}

	ok = ok && v != ""

	return
}

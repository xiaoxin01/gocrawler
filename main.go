package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"gocrawler/model"
	"gocrawler/service"

	"github.com/gocolly/colly"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	defaultSchedule = "CRON_TZ=Asia/Shanghai * 18 * * *"
)

func init() {
	model.InitConfig(".")
	model.InitDb()
	model.InitField("")
}

func main() {
	defer func() {
		err := recover()
		if err != nil {
			service.Alert("crawl failed!", fmt.Sprintf("%v", err))
		}
	}()

	rand.Seed(time.Now().UnixNano())

	db := model.Db

	fmt.Println("Connected to MongoDB!")

	cronJobs := cron.New(cron.WithChain(
		cron.SkipIfStillRunning(cron.DefaultLogger),
	))

	for _, web := range model.Webs {
		if !web.Enabled {
			continue
		}

		schedule := defaultSchedule
		if web.Schedule != nil {
			schedule = *web.Schedule
		}
		realWeb := web
		collectionName := db.Collection
		if realWeb.Collection != nil {
			collectionName = *realWeb.Collection
		}
		collection := model.Database.Collection(collectionName)
		_, err := cronJobs.AddFunc(schedule, func() { crawlWeb(realWeb, collection) })
		if err != nil {
			panic(err)
		}
	}

	cronJobs.Start()
	for {
		for _, entity := range cronJobs.Entries() {
			fmt.Printf("%d, next run: %v, last run: %v\n", entity.ID, entity.Next, entity.Prev)
		}
		time.Sleep(time.Hour)
	}
}

func crawlWeb(web model.Web, collection *mongo.Collection) {
	fmt.Printf("start: %s\n", web.Name)
	initialState(&web)
	os.Mkdir("data", 0744)
	f, err := os.OpenFile(fmt.Sprintf("data/%s.json", web.Name), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	defer f.Close()

	if err != nil {
		panic(err)
	}

	// Instantiate default collector
	c := colly.NewCollector()
	c.SetRequestTimeout(time.Duration(time.Minute * 2))
	c.RedirectHandler = func(req *http.Request, via []*http.Request) error {
		fmt.Println("redirect to: ", req.Response.Header["Location"])
		service.Alert("redirect", strings.Join(req.Response.Header["Location"], " "))
		return nil
	}

	// Before making a request put the URL with
	// the key of "url" into the context of the request
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Del("User-Agent")
		for key, value := range web.Headers {
			r.Headers.Add(key, value)
		}
	})

	c.OnHTML(web.ListSelector, func(e *colly.HTMLElement) {
		obj := make(map[string]interface{}, 0)

		for key, field := range web.Fields {
			if value, ok := field.GetValue(e); ok {
				obj[key] = value
			}
		}

		if len(obj) < web.MinFields {
			return
		}

		// check cache
		if web.ItemKey != nil {
			if key, ok := obj[*web.ItemKey]; ok {
				if _, ok := web.VisitedItems[fmt.Sprintf("%v", key)]; ok {
					// item visited, skip
					return
				}
				web.VisitedItems[fmt.Sprintf("%v", key)] = true
			}
		}

		objString, _ := json.Marshal(obj)

		f.WriteString(string(objString))
		f.WriteString("\n")

		_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": obj["_id"]}, bson.M{"$set": obj}, options.Update().SetUpsert(true))

		if err != nil {
			panic(fmt.Errorf("%v", err))
		}

		if web.Subscribe {
			checkAndSendAlert(string(objString))
		}
		//fmt.Println(string(objString))
	})

	c.OnResponse(func(r *colly.Response) {
		// r.Body
		// fmt.Println("code: ", r.StatusCode)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("visit url: ", r.Request.URL, "failed.", string(r.Body), " error:", err)
		web.Visited[r.Request.URL.String()] = false
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

	saveState(&web)
}

func initialState(web *model.Web) {
	// initial page visit
	bytes, err := ioutil.ReadFile(fmt.Sprintf("data/%s.state", web.Name))
	web.Visited = make(map[string]bool)
	if err == nil {
		json.Unmarshal(bytes, &web.Visited)
	}

	// initial item visit
	bytes, err = ioutil.ReadFile(fmt.Sprintf("data/%s.items", web.Name))
	web.VisitedItems = make(map[string]bool)
	if err == nil {
		json.Unmarshal(bytes, &web.VisitedItems)
	}
}

func saveState(web *model.Web) {
	// save page visit
	file, _ := json.MarshalIndent(web.Visited, "", " ")
	_ = ioutil.WriteFile(fmt.Sprintf("data/%s.state", web.Name), file, 0644)

	// save item visit
	file, _ = json.MarshalIndent(web.VisitedItems, "", " ")
	_ = ioutil.WriteFile(fmt.Sprintf("data/%s.items", web.Name), file, 0644)
}

func checkAndSendAlert(item string) bool {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	var subscribes []string = getSubscribes()

	for _, subscribe := range subscribes {
		if strings.Contains(item, subscribe) {
			service.Alert("subscribe new!", item)
			return true
		}
	}

	return false
}

func getSubscribes() []string {
	var subscribes []string
	collection := model.Database.Collection("subscribes")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		subscribes = append(subscribes, result["_id"].(string))
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return subscribes
}

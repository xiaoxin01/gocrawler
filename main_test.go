package main

import (
	"context"
	"fmt"
	"gocrawler/model"
	"math/rand"
	"regexp"
	"testing"
	"time"

	"github.com/liuzl/gocc"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestRegex(t *testing.T) {
	t.Run("test regex match", func(t *testing.T) {
		url := "details.php?id=123456&ref="
		finds := regexp.MustCompile("id=(\\d+)").FindStringSubmatch(url)
		assert.Equal(t, len(finds), 2)
		assert.Equal(t, finds[1], "123456")
	})
}

func TestMongoDb(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	viper.AddConfigPath(".")
	viper.SetConfigName("webs")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig() // Find and read the config file
	assert.Equal(t, nil, err)

	// Set client options
	clientOptions := options.Client().ApplyURI(model.Db.Connection)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	assert.Equal(t, nil, err)

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	assert.Equal(t, nil, err)

	fmt.Println("Connected to MongoDB!")
	collection := client.Database("test_database").Collection("tests")
	obj, objFind := make(map[string]interface{}, 0), make(map[string]interface{}, 0)
	id := rand.Int63n(100000)
	obj["id"] = id
	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": "asbsdf"}, bson.M{"$set": obj}, options.Update().SetUpsert(true))
	assert.Equal(t, nil, err)

	err = collection.FindOne(context.TODO(), bson.M{"_id": "asbsdf"}).Decode(&objFind)
	assert.Equal(t, nil, err)
	assert.Equal(t, id, objFind["id"].(int64))
}

func TestCronjob(t *testing.T) {
	c := cron.New()
	id, err := c.AddFunc("CRON_TZ=Asia/Shanghai 1 * * * *", func() { fmt.Println("Every hour on the half hour") })
	assert.Equal(t, nil, err)

	c.Start()
	entry := c.Entry(id)
	now := time.Now()
	//shanghaiLocation, _ := time.LoadLocation("Asia/Shanghai")
	targetTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 1, 0, 0, time.Local)
	if targetTime.Before(now) {
		targetTime = targetTime.Add(time.Hour)
	}
	assert.Equal(t, targetTime, entry.Next)
}

func TestT2s(t *testing.T) {
	*gocc.Dir = `./module/gocc`
	t2s, err := gocc.New("t2s")
	assert.Equal(t, err, nil)

	got, err := t2s.Convert("自然語言處理是人工智能領域中的一個重要方向")
	assert.Equal(t, err, nil)
	assert.Equal(t, got, "自然语言处理是人工智能领域中的一个重要方向")
}

func TestGetSubscribes(t *testing.T) {
	subscribes := getSubscribes()

	assert.Equal(t, 1, len(subscribes))
	assert.Equal(t, "测试", subscribes[0])
}

func TestCheckAndSendAlert(t *testing.T) {
	alerted1 := checkAndSendAlert("测试")
	alerted2 := checkAndSendAlert("test")

	assert.Equal(t, true, alerted1)
	assert.Equal(t, false, alerted2)
}

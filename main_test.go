package main

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
	"github.com/spf13/viper"
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

	var db Db
	err = viper.UnmarshalKey("db", &db)
	assert.Equal(t, nil, err)

	// Set client options
	clientOptions := options.Client().ApplyURI(db.Connection)

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

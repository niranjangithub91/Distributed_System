package db

import (
	"context"
	"fmt"
	"log"
	"time"
	"user_entry/model"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var collections1 *mongo.Collection
var collections2 *mongo.Collection
var rdb *redis.Client

func init() {
	clientoptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientoptions)
	if err != nil {
		log.Fatal(err)
	}
	collections1 = client.Database("Distributed_systems").Collection("Users")
	collections2 = client.Database("Distributed_systems").Collection("MetaData")

	fmt.Println("Database connected successfully")

	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password by default
		DB:       0,                // Default DB
	})

	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	} else {
		fmt.Println("Connected to Redis successfully!")
	}

}

func Add_User(t model.User) bool {
	insert, err := collections1.InsertOne(context.Background(), t)
	if err != nil {
		log.Fatal(err)
		return false
	}
	fmt.Println(insert)
	return true
}
func Find_User(t model.User) bool {
	curr, err := collections1.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	for curr.Next(context.Background()) {
		var a model.User
		err := curr.Decode(&a)
		if err != nil {
			log.Fatal(err)
		}
		if t.Name == a.Name && t.Password == a.Password {
			return true
		}
	}
	return false
}

func Collection_update(t model.Collection) {
	insert, err := collections2.InsertOne(context.Background(), t)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(insert)
}

func Collection_retreive_details(filename string, user string) map[int][]int {
	curr, err := collections2.Find(context.Background(), bson.M{})
	j := make(map[int][]int)
	if err != nil {
		log.Fatal(err)
		return j
	}
	for curr.Next(context.Background()) {
		var a model.Collection
		err := curr.Decode(&a)
		if err != nil {
			log.Fatal(err)
		}
		if a.File_name == filename && a.Name == user {
			j[a.Port] = append(j[a.Port], a.Chunk)
		}
	}
	return j
}
func Set_cache_mem(key string, value []byte) {
	err := rdb.Set(context.Background(), key, value, 1*time.Minute).Err() // 0 = no expiration
	if err != nil {
		log.Fatal(err)
	}
	return
}

func Get_cache_mem(key string) []byte {
	var q []byte
	val, err := rdb.Get(context.Background(), key).Bytes()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("key not found in cache")
			return q
		}
		fmt.Println("error in getting from cache:", err)
		return q
	}
	return val
}

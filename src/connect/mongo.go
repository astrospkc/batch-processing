package connect

import (
	"batch-processing/env"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	dbName              = "CMS_portfolio"
	colPosts      = "posts"
	colNameUsers        = "users"
	colLikes      = "likes"
	
)
var UsersCollection *mongo.Collection
var PostsCollection *mongo.Collection
var PostLikesCollection *mongo.Collection
// var MongoClient
func MongoConnect(){
	envs:=env.NewEnv()
	uri:=envs.MONGODB_URI
	if uri==" "{
		log.Fatal("get the uri first , uri must be invalid or empty string")
	}
	ctx, cancel:=context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts:=options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client,err:=mongo.Connect(opts)
	if err != nil {
		log.Fatal("the error while connecting mongo client", err)
		return
	}

	err = client.Ping(ctx, nil)

	if err != nil {
		log.Fatal("ping :", err)
		return
	}

	UsersCollection  = client.Database(dbName).Collection(colNameUsers)

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: -1}},
		Options: options.Index().SetUnique(true),
	}
	_, err = UsersCollection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		log.Fatal("error occured while connecting to users collection : ", err)
	}


	PostsCollection = client.Database(dbName).Collection(colPosts)
	PostLikesCollection = client.Database(dbName).Collection(colLikes)
	indexModel = mongo.IndexModel{
		Keys: bson.D{{Key:"user_id", Value:-1}},
		Options: options.Index().SetUnique(true),
	}
	_,err = PostsCollection.Indexes().CreateOne(context.TODO(), indexModel)
	if err!=nil{
		log.Fatal("error occurred while connecting to posts collection: ",err)
	}
	_,err = PostsCollection.Indexes().CreateOne(context.TODO(), indexModel)
	if err!=nil{
		log.Fatal("error occurred while connecting to post_likes collection: ",err)
	}
	fmt.Println("Set up is done")

}
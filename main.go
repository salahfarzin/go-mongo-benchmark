package main

import (
	"context"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
)

type Post struct {
	ID    primitive.ObjectID `bson:"_id, omitempty"`
	Title string
	Data  string
}

func main() {

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	r := router.New()
	r.GET("/", func(ctx *fasthttp.RequestCtx) {
		collection := client.Database("test").Collection("posts")

		post := Post{}
		if err = collection.FindOne(ctx, bson.M{"title": "insert a value \n"}).Decode(&post); err != nil {
			log.Fatal(err)
		}

		// update
		collection.UpdateOne(
			ctx,
			bson.M{"_id": post.ID},
			bson.D{
				{"$set", bson.D{{"data", post.Data + " test "}}},
			},
		)

		ctx.WriteString(post.Data)
	})

	r.GET("/fakeData", func(ctx *fasthttp.RequestCtx) {
		collection := client.Database("test").Collection("posts")

		byteValues, err := ioutil.ReadFile("data.json")
		if err != nil {
			log.Fatal(err)
		}

		//insert 300,000 documents about 6 Kb per record
		for i := 1; i <= 300000; i++ {
			collection.InsertOne(context.TODO(), bson.M{
				"title": "insert a value " + string(i), "data": string(byteValues),
			})
		}

		ctx.WriteString("Done!")
	})

	log.Fatal(fasthttp.ListenAndServe(":8080", r.Handler))
}

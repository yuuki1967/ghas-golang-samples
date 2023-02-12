package controllers

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"example.com/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Connection URI
const uri = "mongodb://localhost:27017/?maxPoolSize=20&w=majority"

type Address struct {
	Street string
	City   string
	State  string
}

type Student struct {
	FirstName string  `bson:"first_name,omitempty"`
	LastName  string  `bson:"last_name,omitempty"`
	Address   Address `bson:"inline"`
	Age       int
	MyNumber  string
}

func dbInit() (client *mongo.Client) {
	// return client
	// Create a new client and connect to the server
	cl, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}
	defer func() {
		if err = cl.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Ping the primary
	if err = cl.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected and pinged.")
	client = cl
	return client
}

func dbInsertOne(client *mongo.Client) (collection *mongo.Collection) {
	coll := client.Database("school").Collection("students")
	address1 := Address{"1 Lakewood Way", "Elwood City", "PA"}
	student1 := Student{FirstName: "Arthur", Address: address1, Age: 8, MyNumber: "111111"}
	_, err := coll.InsertOne(context.TODO(), student1)
	fmt.Printf("err = %s", err)
	collection = coll
	return collection
}

func dbFindOne(coll *mongo.Collection, mynum string) {
	filter := bson.D{{"mynumber", "111111"}}

	var result bson.D
	err := coll.FindOne(context.TODO(), filter).Decode(&result)
	fmt.Printf("err = %s", err)
	fmt.Println(result)
}

func hello(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	b, _ := io.ReadAll(req.Body)
	mynum := string(b)
	logger.Log().Errorf("user %s logged in.\n", mynum)
	fmt.Fprintf(w, "hello %s\n", mynum)

	/*
		client := dbInit()
		coll := dbInsertOne(client)
		dbFindOne(coll, mynum)

	*/
	cl, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}
	defer func() {
		if err = cl.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Ping the primary
	if err = cl.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected and pinged.")

	coll := cl.Database("school").Collection("students")
	address1 := Address{"1 Lakewood Way", "Elwood City", "PA"}
	student1 := Student{FirstName: "Arthur", Address: address1, Age: 8, MyNumber: "111111"}
	_, err = coll.InsertOne(context.TODO(), student1)

	filter := bson.D{{"first_name", mynum}}

	var result bson.D
	err = coll.FindOne(context.TODO(), filter).Decode(&result)
	fmt.Println(result)
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func Controllersbody() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	http.ListenAndServe(":8090", nil)
}

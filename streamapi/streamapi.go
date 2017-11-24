package main

import (
	"gopkg.in/mgo.v2"
	"flag"
	"log"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/Machinestream/types"
	"gopkg.in/mgo.v2/bson"
	"fmt"
)

var mongoAddress = flag.String("mongo", "mongodb:27017", "mongoDB address")
var database = flag.String("database", "streamengine", "mongoDB database name")

func main() {
	flag.Parse()
	log.SetFlags(0)

	// Init MongoDB connection
	log.Printf("connecting to %s", *mongoAddress)
	session, err := mgo.Dial(*mongoAddress)
	if err != nil {
		log.Fatal("connectionto MongoDB failed:", err)
	}
	defer session.Close()
	collection := session.DB(*database).C("events")

	engine := gin.Default()
	engine.GET("machines", func(c *gin.Context) {
		result := []types.Event{}
		collection.Find(nil).All(&result)
		c.JSON(http.StatusOK, result)
	})

	engine.GET("machines/:machine_id", func(c *gin.Context) {
		machineId := c.Param("machine_id")
		query := collection.Find(bson.M{"payload.machineid": machineId})
		count, err := query.Count()
		if err != nil {
			log.Fatal("failed to query mongodb")
			c.String(http.StatusInternalServerError, "failed to query mongodb")
		}

		if count == 0 {
			c.String(http.StatusNotFound, fmt.Sprintf("machine with machineId %v not found", machineId))
		} else if count == 1{
			result := types.Event{}
			query.One(&result)
			c.JSON(http.StatusOK, result)
		} else {
			c.String(http.StatusBadRequest, fmt.Sprintf("machine_id %v not identifies one machine uniquely", machineId))
		}
	})

	engine.Run()
}

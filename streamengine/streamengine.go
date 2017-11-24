package main

import (
	"github.com/gorilla/websocket"
	"log"
	"flag"
	"os"
	"net/url"
	"os/signal"
	"time"
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/Machinestream/types"
)

var websocketAddress = flag.String("websocket", "machinestream.herokuapp.com", "websocket address")
var mongoAddress = flag.String("mongo", "mongodb:27017", "mongoDB address")
var database = flag.String("database", "streamengine", "mongoDB database name")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *websocketAddress, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	// Init websocket connection
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	// Init MongoDB connection
	log.Printf("connecting to %s", *mongoAddress)
	session, err := mgo.Dial(*mongoAddress)
	if err != nil {
		log.Fatal("connectionto MongoDB failed:", err)
	}
	defer session.Close()
	collection := session.DB(*database).C("events")

	done := make(chan struct{})
	events := readEvents(done, conn)
	processEvents(done, events, collection)

	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Printf("Sending ping")
			err := conn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")
			close(done)
			session.Close()
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			conn.Close()
			return
		}
	}
}

func readEvents(done <-chan struct{}, conn *websocket.Conn) <-chan types.Event {
	events := make(chan types.Event)

	go func() {
		defer conn.Close()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			event := types.Event{}
			err = json.Unmarshal(message, &event)
			if err != nil {
				log.Println("json.unmarshal:", err)
			}
			events <- event
		}
	}()

	return events
}

func processEvents(done <-chan struct{}, events <-chan types.Event, collection *mgo.Collection) {
	go func(){
		for {
			event := <-events
			query := collection.Find(bson.M{"payload.machineid": event.Payload.MachineId})
			count, err := query.Count()
			if err != nil {
				log.Fatal("unable to query mongodb", err)
			}

			if count == 0 {
				err = collection.Insert(event)
				if err != nil {
					log.Fatal("failed to insert event into database", err)
				} else {
					log.Printf("database record inserted: %v", event.Payload.MachineId)
				}
			} else if count == 1 {
				result := types.Event{}
				err = query.One(&result)
				if err != nil {
					log.Fatal("unable to retrieve mongoDB record")
				}
				err = collection.Update(result, event)
				if err != nil {
					log.Fatal("failed to update the monboDB record")
				} else {
					log.Printf("database record updated: %v", event.Payload.MachineId)
				}
			} else {
				log.Fatal("found multiple records in database for machineID", event.Payload.MachineId)
			}
		}
	}()
}

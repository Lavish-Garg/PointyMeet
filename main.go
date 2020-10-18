package main

import (
	"context"
	"encoding/json"
	"fmt"
  "log"
	"net/http"
	"time"
	"regexp"
  "strconv"
   "strings"
        "go.mongodb.org/mongo-driver/bson"
        "go.mongodb.org/mongo-driver/bson/primitive"
        "go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

type Meeting struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title     string             `json:"title,omitempty" bson:"title,omitempty"`
	Participants  string        `json:"participants,omitempty" bson:"participants,omitempty"`
        Start_Time    string         `json:"start_time,omitempty" bson:"start_time,omitempty"`
        End_Time      string         `json:"end_time,omitempty" bson:"end_time,omitempty"`
        Creation_TimeStamp  time.Time   `json:"creation_timestamp,omitempty" bson:"creation_timestamp,omitempty"`
        
}

type Participant struct {
        Name string                  `json:"Name,omitempty" bson:"Name,omitempty"`
        Email string                 `json:"email,omitempty" bson:"email,omitempty"`
        RSVP string                  `json:"rsvp,omitempty" bson:"rsvp,omitempty"`
}

func CreateMeetingEndpoint(response http.ResponseWriter, request *http.Request) {
response.Header().Set("content-type", "application/json")
	var meeting Meeting
	json.NewDecoder(request.Body).Decode(&meeting)
	collection := client.Database("Appointy").Collection("participant")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, meeting)
	json.NewEncoder(response).Encode(result) 
}

func GetMeetingsEndpoint(response http.ResponseWriter, request *http.Request) { 
response.Header().Set("content-type", "application/json")
         var meetings []Meeting
         collection := client.Database("Appointy").Collection("participant")
         
        json.NewEncoder(reponse).Encode(meetings)

}

func GetMeetingEndpoint(response http.ResponseWriter, request *http.Request) { 
response.Header().Set("content-type", "application/json")
	slug := getField(request,0)
	id, _ := primitive.ObjectIDFromHex(slug["id"])
	var meeting Meeting
	collection := client.Database("Appointy").Collection("participant")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Meeting{ID: id}).Decode(&meeting)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(meeting)
}
func GetParticipantEndpoint(response http.ResponseWriter, request *http.Request) {

var participant []Participant
         collection := client.Database("Appointy").Collection("participant")
         
        json.NewEncoder(reponse).Encode(participant)

}

var routes = []route{
         newRoute("POST", "/meetings", CreateMeetingEndpoint),
         newRoute("GET", "/meeting/([0-9]+)", GetMeetingEndpoint),
         newRoute("GET", "/([\s\S])", GetMeetingsEndpoint),
         newRoute("GET", "/([\s\S])", GetParticipantEndpoint),
}
func newRoute(method, pattern string, handler http.HandlerFunc) route {
	return route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}



func main() {
	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
        clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	http.ListenAndServe(":12345", routes)
}

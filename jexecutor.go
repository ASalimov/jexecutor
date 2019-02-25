package main

import (
	"encoding/json"
	"flag"
	"github.com/ASalimov/jexecutor/logic"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
)

func main() {
	flag.Parse()
	logrus.SetLevel(logrus.DebugLevel)
	filePath := flag.Arg(0)
	if filePath == "" {
		filePath = "test.json"
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("failed to open file %s, %s", filePath, err)
	}
	var jInfo map[string]interface{}
	err = json.Unmarshal(data, &jInfo)
	if err != nil {
		log.Fatalf("failed to unmarshal file %s", err)
	}

	URL := jInfo["url"].(string)
	username := jInfo["username"].(string)
	token := jInfo["token"].(string)
	threads := int(jInfo["threads"].(float64))
	jobs := jInfo["jobs"].([]interface{})
	for i := 0; i < len(jobs); i++ {
		log.Printf("Goal №%d, starting...", i+1)
		e := logic.NewExecutor(threads, URL, username, token)
		jobs1 := jobs[i].([]interface{})
		orders := make([]logic.Order, len(jobs1))
		for j := 0; j < len(jobs1); j++ {
			job1 := jobs1[j].(map[string]interface{})
			order := logic.Order{Job: job1["id"].(string), Query: job1["q"].(map[string]interface{})}
			orders[j] = order
		}
		e.Execute(orders...)
		log.Printf("Goal №%d, finished", i+1)
	}

}

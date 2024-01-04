package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	c8s "github.com/furon-kuina/cuternetes/pkg"
)

func notifyEvent(event c8s.ContainerEvent) (err error) {
	c8s.Wrap(&err, "notifyEvent(%q)", event)
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	resp, err := http.Post(c8sConfig.ApiServer.Url+"/event", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	log.Printf("response of notifyEvent(%q): %+v", event, resp)
	return nil
}

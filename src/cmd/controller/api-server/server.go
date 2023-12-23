package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/docker/docker/api/types/strslice"
	restful "github.com/emicklei/go-restful/v3"
)

func main() {
	ws := new(restful.WebService)
	ws.Route(ws.POST("/apply").To(postApply))
	restful.Add(ws)

	http.ListenAndServe(":7070", nil)
}

func postApply(req *restful.Request, resp *restful.Response) {
	container := new(Container)
	if err := req.ReadEntity(&container); err != nil {
		resp.WriteHeaderAndJson(http.StatusInternalServerError, &resError{"Invalid Response"}, restful.MIME_JSON)
	}
	if container.Image == "" {
		resp.WriteHeaderAndJson(http.StatusInternalServerError, &resError{"Image name required"}, restful.MIME_JSON)
	}

	payload, err := json.Marshal(container)
	if err != nil {
		resp.WriteHeaderAndJson(http.StatusInternalServerError, &resError{"json.Marshal failed"}, restful.MIME_JSON)
	}

	workerReq, err := http.NewRequest("POST", "http://127.0.0.1:8080", bytes.NewBuffer(payload))
	if err != nil {
		resp.WriteHeaderAndJson(http.StatusInternalServerError, &resError{"http.NewRequest failed"}, restful.MIME_JSON)
	}
	workerReq.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	workerResp, err := client.Do(workerReq)
	if err != nil {
		resp.WriteHeaderAndJson(http.StatusInternalServerError, &resError{"request to worker failed"}, restful.MIME_JSON)
	}
	defer workerResp.Body.Close()
	body, err := io.ReadAll(workerResp.Body)

	if err != nil {
		resp.WriteHeaderAndJson(http.StatusInternalServerError, &resError{"invalid response from worker"}, restful.MIME_JSON)
	}
	fmt.Println("Response from worker: ", string(body))

	resp.WriteAsJson(container)

}

type Container struct {
	Image string
	Cmd   strslice.StrSlice `json:",omitempty"`
}

type resError struct {
	Error string
}

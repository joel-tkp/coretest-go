package base

import (
	"net/http" // base http handler
	"encoding/json" // JSON Response Support
	"time" // time lib
)

const Success = "200"
const Error = "400"

type Base struct {
    Status            string      `json:"status"`
    Config            interface{} `json:"config"`
    ServerProcessTime string      `json:"server_process_time"`
    ErrorMessage      []string    `json:"message_error,omitempty"`
    StatusMessage     []string    `json:"message_status,omitempty"`
}

type Response struct {
    Base
    Data interface{} `json:"data"`
}

func WriteSuccessResponse(w http.ResponseWriter, statusCode string, statusMessage string, start time.Time, data interface{}) {
    var response Response
    response.Status = statusCode
    response.Base.ServerProcessTime = time.Since(start).String()
    response.Data = data
    if statusCode == Error {
        response.Base.ErrorMessage = []string { statusMessage, "" }
    } else {
        response.Base.StatusMessage = []string { statusMessage, "" }
    }

    b, err := json.Marshal(response)
    if err != nil {
        panic(err.Error())
    }

    w.Header().Set("content-type", "application/json")
    w.WriteHeader(200)
    w.Write(b)
}

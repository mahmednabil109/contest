package loggerhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/linuxboot/contest/cmds/admin_server/storage"
	"github.com/sirupsen/logrus"
)

type HttpHook struct {
	Addr string
}

func NewHttpHook(addr string) *HttpHook {
	fmt.Println("creating new http logger hook")
	// add the endpoint to the server addr
	if addr[len(addr)-1] != '/' {
		addr += "/"
	}
	addr += "log"

	return &HttpHook{
		Addr: addr,
	}
}

// this implements logrus Hook interface
func (hh *HttpHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// this implements logrus Hook interface
func (hh *HttpHook) Fire(entry *logrus.Entry) error {
	msg := strings.TrimRight(entry.Message, "\n")
	if msg == "" {
		return nil
	}
	jobId, ok := entry.Data["job_id"]
	jobIdInt, noInt := jobId.(int)
	if !ok || !noInt {
		// to indicate an invalid job id
		jobIdInt = -1
	}

	logJson, err := json.Marshal(storage.Log{
		LogData: msg,
		JobID:   jobIdInt,
	})
	if err != nil {
		return err
	}

	requestBody := bytes.NewBuffer(logJson)
	http.Post(hh.Addr, "application/json", requestBody)
	return nil
}

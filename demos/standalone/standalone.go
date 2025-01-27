package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Xunzhuo/async"
	log "github.com/sirupsen/logrus"
)

func main() {

	async.DefaultAsyncQueue().Start()

	stop := make(chan bool)
	stopData := make(chan bool)
	jobID := make(chan async.Job, 1000)

	go func() {
		for {
			select {
			case _, ok := <-stop:
				if !ok {
					return
				}
				return
			default:
				id := fmt.Sprintf("%d", rand.Intn(1000000))
				job, err := async.NewJob(longTimeJob, "xunzhuo")
				if err != nil {
					log.Error("New Job Failed:%v", err)
					continue
				}
				if ok := async.DefaultAsyncQueue().AddJobAndRun(job); ok {
					jobID <- *job
					log.Warning("Send Job ID: ", id)
				} else {
					log.Warning("Reject Job ID: ", id)
				}
			}
		}
	}()

	time.Sleep(5 * time.Second)
	stop <- true
	close(stop)

	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			select {
			case _, ok := <-stopData:
				if !ok {
					return
				}
				return
			case job := <-jobID:
				log.Warning("Received Job ID: ", job.JobID)
				if data, ok := async.DefaultAsyncQueue().GetJobData(job); ok {
					log.Warningf(fmt.Sprintf("Get data from WorkQueue %s with ID: %s", data[0].(string), job.GetJobID()))
				}
			}
		}
	}()

	time.Sleep(10 * time.Second)
	stopData <- true
	close(stopData)
}

func longTimeJob(value string) string {
	time.Sleep(1000 * time.Millisecond)
	return "Hello World from " + value
}

/*
Copyright © 2020 Vlad Krava <vkrava4@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"curlson/logutil"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
	"net/http"
	"os"
	"sync"
	"time"
)

var count int
var sleepMs int
var threads int
var duration int
var persistLogs = false
var loggingSupported = false

var getCmd = &cobra.Command{
	Use:   "get <URL>",
	Short: "Performs HTTP GET request with options",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ValidateInput()
		Run(args[0])
	},
}

var yellowColor = color.New(color.FgYellow)
var greenColor = color.New(color.FgGreen)
var redColor = color.New(color.FgRed)

var log = logrus.New()

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().IntVarP(&threads, "threads", "t", 1, "A number of concurrent GET requests")
	getCmd.Flags().IntVarP(&count, "count", "c", 1, "A number of GET requests per single thread")
	getCmd.Flags().IntVarP(&sleepMs, "sleep", "s", 0, "A delay in millis after each GET requests. Doesn't impact performance report results if set (default 0)")
	getCmd.Flags().IntVarP(&duration, "duration", "d", 0, "A maximum duration in seconds by reaching which requests execution will be terminated regardless of a 'count' flag value. When the value set to '0' this flag is ignored")
	getCmd.Flags().BoolVarP(&persistLogs, "persist-logs", "p", false, "A property which defines whether execution log files will be persisted or automatically cleaned up")
}

func ValidateInput() {
	var validationResults = ""
	if threads < 1 {
		validationResults = fmt.Sprintf(validationResults+"\n - The number of threads has to be grater or equal to 1. Currently it's: '%d'", threads)
	}

	if count < 1 {
		validationResults = fmt.Sprintf(validationResults+"\n - The amount of requests per thread has to be grater or equal to 1. Currently it's: '%d'", count)
	}

	if sleepMs < 0 {
		validationResults = fmt.Sprintf(validationResults+"\n - The delay in millis per request has to be grater or equal to 0. Currently it's: '%d' millis", sleepMs)
	}

	if duration < 0 {
		validationResults = fmt.Sprintf(validationResults+"\n - The maximum execution duration in seconds has to be grater or equal to 0. Currently it's: '%d' seconds", duration)
	}

	if len(validationResults) > 0 {
		fmt.Println(redColor.Sprintf(validationResults))
		os.Exit(1)
	}
}

func Run(url string) {
	var supported, logFile = logutil.SetupLogs(log, &persistLogs)
	loggingSupported = supported
	defer logutil.ShutdownLogs(logFile, &persistLogs)

	logutil.InfoLog(fmt.Sprintf("Setting up GET execution to URL address %s with threads = %d, amount of requests = %d, sleep millis timeout = %d", url, threads, count, sleepMs), log, &supported)

	var executionResults = make([]int, threads)
	var waitGroup sync.WaitGroup
	var multiProgress = mpb.New(mpb.WithWaitGroup(&waitGroup))
	waitGroup.Add(threads)

	for i := 0; i < len(executionResults); i++ {
		logutil.InfoLog(fmt.Sprintf("Setting up new thread with id: %d", i), log, &loggingSupported)
		go ThreadStart(i, &url, executionResults, multiProgress, &waitGroup)
		logutil.InfoLog(fmt.Sprintf("Thread with id: %d successfully started", i), log, &loggingSupported)
	}

	MonitorActivity(executionResults)
	multiProgress.Wait()

}

func ThreadStart(threadId int, url *string, executionResults []int, multiProgress *mpb.Progress, waitGroup *sync.WaitGroup) {
	var threadDescription = yellowColor.Sprintf("Thread #%-4d", threadId)

	var progress = multiProgress.AddBar(int64(count),
		mpb.BarStyle("╢▌▌ ╟"),
		mpb.PrependDecorators(
			decor.Name(threadDescription),
			decor.Percentage(decor.WCSyncSpace),
		),
		mpb.AppendDecorators(
			decor.OnComplete(
				decor.EwmaETA(decor.ET_STYLE_GO, 10),
				greenColor.Sprintf("COMPLETE"),
			),
		),
	)
	defer waitGroup.Done()

	var maxExecutionEndTime = time.Now().Add(time.Second * time.Duration(duration))
	logutil.InfoLog(fmt.Sprintf("Determined maximum execution duration time: %#v for thread with id: %d", maxExecutionEndTime, threadId), log, &loggingSupported)

	for i := 0; i < count; i++ {
		var requestStartTime = time.Now()
		var getResponse, getResponseErr = http.Get(*url)

		if getResponseErr == nil {
			if getResponse.StatusCode >= 200 && getResponse.StatusCode <= 299 {
				logutil.InfoLog(fmt.Sprintf("Successfully received HTTP GET responce from address '%s' with body %#v", *url, getResponse), log, &loggingSupported)
			} else {
				logutil.WarnLog(fmt.Sprintf("Received HTTP GET responce with status code: %d from address '%s' with body %#v", getResponse.StatusCode, *url, getResponse), log, &loggingSupported)
			}
		} else {
			logutil.ErrorLog(fmt.Sprintf("Received an error on HTTP GET request from address: '%s' with message: %s", *url, getResponseErr.Error()), log, &loggingSupported)
		}

		progress.Increment(time.Since(requestStartTime))

		if sleepMs > 0 {
			logutil.InfoLog(fmt.Sprintf("Sleeping thread with id: %d for %d millis before the next itteration", threadId, sleepMs), log, &loggingSupported)
			time.Sleep(time.Millisecond * time.Duration(sleepMs))
			logutil.InfoLog(fmt.Sprintf("Resumed thread with id: %d after sleeping for %d millis", threadId, sleepMs), log, &loggingSupported)
		}

		if duration != 0 && maxExecutionEndTime.Before(time.Now()) {
			logutil.WarnLog(fmt.Sprintf("Exceeded maximum execution duration time: %#v. Terminating execution of thread with id: %d", maxExecutionEndTime, threadId), log, &loggingSupported)
			progress.SetTotal(progress.Current(), true)
			break
		}
	}

	executionResults[threadId] = 1
}

func MonitorActivity(executionResults []int) {
	for true {
		var RunningThreads = 0
		for _, value := range executionResults {
			if value == 0 {
				RunningThreads++
			}
		}

		if RunningThreads == 0 {
			break
		}

		time.Sleep(time.Millisecond * time.Duration(200))
	}
}

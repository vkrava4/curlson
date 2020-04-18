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
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
	"github.com/vkrava4/curlson/cmd/util"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var count int
var sleepMs int
var threads int
var maxDuration int
var template string
var persistLogs = false
var verbose = false
var loggingSupported = false
var templateEnabled = false

var getCmd = &cobra.Command{
	Use:   "get <URL> [flags]",
	Short: "Performs HTTP GET request(s) based on specified options",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ValidateInput()
		runGet(args[0])
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
	getCmd.Flags().IntVarP(&maxDuration, "duration-max", "D", 0, "A maximum duration in seconds by reaching which requests execution will be terminated regardless of a 'count' flag value. When the value set to '0' this flag is ignored (default 0)")
	getCmd.Flags().StringVarP(&template, "template-file", "T", "", "")
	getCmd.Flags().BoolVarP(&persistLogs, "persist-logs", "p", false, "A flag which defines whether execution log files will be persisted or automatically cleaned up")
	getCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "A flag which defines whether additional execution information such as log creations or other actions will be logged in console output")
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

	if maxDuration < 0 {
		validationResults = fmt.Sprintf(validationResults+"\n - The maximum execution duration in seconds has to be grater or equal to 0. Currently it's: '%d' seconds", maxDuration)
	}

	if template != "" {
		var absTemplatePath, absError = filepath.Abs(template)
		if absError != nil {
			validationResults = fmt.Sprintf(validationResults+"\n - An error occured while identifying template file path: %s", absError.Error())
		}

		if util.FileExists(absTemplatePath) {
			var templateFile, errOpenTemplate = os.OpenFile(template, os.O_RDONLY, 0666)

			if errOpenTemplate != nil {
				validationResults = fmt.Sprintf(validationResults+"\n - An error occured while openning template file: %s", errOpenTemplate.Error())
			} else {
				template = absTemplatePath
				var errCloseTemplate = templateFile.Close()
				if errCloseTemplate != nil {
					validationResults = fmt.Sprintf(validationResults+"\n - An error occured while working with template file: %s", errCloseTemplate.Error())
				}
				templateEnabled = true
			}
		} else {
			validationResults = fmt.Sprintf(validationResults+"\n - The specyfied template file: '%s' can not be found.", absTemplatePath)
		}
	}

	if len(validationResults) > 0 {
		fmt.Println(redColor.Sprintf(validationResults))
		os.Exit(1)
	}
}

func runGet(url string) {
	var supported, logFile = util.SetupLogs(log, &persistLogs, &verbose)
	loggingSupported = supported
	defer util.ShutdownLogs(logFile, &persistLogs, &verbose)

	util.InfoLog(fmt.Sprintf("Setting up GET execution to URL address %s with threads = %d, amount of requests = %d, sleep millis timeout = %d", url, threads, count, sleepMs), log, &supported)

	var executionResults = make([]int, threads)
	var waitGroup sync.WaitGroup
	var multiProgress = mpb.New(mpb.WithWaitGroup(&waitGroup))
	waitGroup.Add(threads)

	var linesCount = util.CountLines(&template, &templateEnabled, log, &loggingSupported)
	for i := 0; i < len(executionResults); i++ {
		util.InfoLog(fmt.Sprintf("Setting up new thread with id: %d", i), log, &loggingSupported)
		go ThreadStart(i, url, executionResults, multiProgress, &waitGroup, linesCount)
		util.InfoLog(fmt.Sprintf("Thread with id: %d successfully started", i), log, &loggingSupported)
	}

	MonitorActivity(executionResults)
	multiProgress.Wait()

}

func ThreadStart(threadId int, url string, executionResults []int, multiProgress *mpb.Progress, waitGroup *sync.WaitGroup, linesCount int) {
	var threadDescription = yellowColor.Sprintf("Thread #%-4d", threadId)

	var onCompleteDecorator = decor.OnComplete(
		decor.EwmaETA(decor.ET_STYLE_GO, 10),
		greenColor.Sprintf("DONE"),
	)

	var progress = multiProgress.AddBar(int64(count),
		mpb.BarStyle("╢▌▌ ╟"),
		mpb.PrependDecorators(
			decor.Name(threadDescription),
			decor.Percentage(decor.WCSyncSpace),
		),
		mpb.AppendDecorators(onCompleteDecorator),
	)
	defer waitGroup.Done()

	var maxExecutionEndTime = time.Now().Add(time.Second * time.Duration(maxDuration))
	util.InfoLog(fmt.Sprintf("Determined maximum execution duration time: %#v for thread with id: %d", maxExecutionEndTime, threadId), log, &loggingSupported)

	for i := 0; i < count; i++ {
		var requestStartTime = time.Now()
		var getUrl string
		if templateEnabled && linesCount > 0 {
			var lineNum, templateLine = util.ReadRandomLine(template, linesCount)
			util.InfoLog(fmt.Sprintf("Received template line %s from the line %d", templateLine, lineNum), log, &loggingSupported)
			var updatedUrl, errPrepareUrl = util.PrepareUrl(url, templateLine)
			if errPrepareUrl != nil {
				util.ErrorLog("Can not make GET request with broken URL. Skipping this iteration", log, &loggingSupported)
				progress.Increment()
				progress.DecoratorEwmaUpdate(time.Since(requestStartTime))
				continue
			}

			util.InfoLog(fmt.Sprintf("Updated URL address ccording to template file %s with line number %d. WAS: %s BECOME: %s", template, lineNum, url, updatedUrl), log, &loggingSupported)
			getUrl = updatedUrl
		} else {
			getUrl = url
		}

		var getResponse, getResponseErr = http.Get(getUrl)

		if getResponseErr == nil {
			util.WarnLog(fmt.Sprintf("Received HTTP GET responce with status code: %d from address '%s' with ContentLength: %d", getResponse.StatusCode, getUrl, getResponse.ContentLength), log, &loggingSupported)
			_ = getResponse.Body.Close()
		} else {
			util.ErrorLog(fmt.Sprintf("Received an error on HTTP GET request from address: '%s' with message: %s", getUrl, getResponseErr.Error()), log, &loggingSupported)
		}

		progress.Increment()
		progress.DecoratorEwmaUpdate(time.Since(requestStartTime))

		if sleepMs > 0 {
			util.InfoLog(fmt.Sprintf("Sleeping thread with id: %d for %d millis before the next itteration", threadId, sleepMs), log, &loggingSupported)
			time.Sleep(time.Millisecond * time.Duration(sleepMs))
			util.InfoLog(fmt.Sprintf("Resumed thread with id: %d after sleeping for %d millis", threadId, sleepMs), log, &loggingSupported)
		}

		if maxDuration != 0 && maxExecutionEndTime.Before(time.Now()) {
			util.WarnLog(fmt.Sprintf("Exceeded maximum execution duration of %d second(s). Terminating execution of thread with id: %d as it did not complete before time: %s", maxDuration, threadId, maxExecutionEndTime.Format(time.RFC3339)), log, &loggingSupported)
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

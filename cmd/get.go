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
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vkrava4/curlson/app"
	"github.com/vkrava4/curlson/ui"
	"github.com/vkrava4/curlson/util"
	"net/http"
	"time"
)

var count int
var sleepMs int
var threads int
var maxDuration int
var template string
var persistLogs = false
var verbose = false

var appConf = &app.Configuration{}

var getCmd = &cobra.Command{
	Use:   "get <URL> [flags]",
	Short: "Performs HTTP GET request(s) based on specified options",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appConf.Logs = &app.LogConfiguration{
			Enabled: false,
			Persist: persistLogs,
			Verbose: verbose,
			Log:     logrus.New(),
		}

		var getValidator = &util.GetValidator{}
		var validatorEntity = getValidator.
			AddRequestCount(count).
			AddThreads(threads).
			AddUrl(args[0]).
			AddTemplate(template).
			AddSleep(sleepMs).
			AddMaxDuration(maxDuration).
			WithAppConfiguration(appConf).
			Entity()

		validatorEntity.Validate().ProcessErrors()

		runGet(args[0])
	},
}

var yellowColor = color.New(color.FgYellow)

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

func runGet(url string) {
	var errSetupLogs = util.SetupLogs(appConf.Logs)
	defer util.ShutdownLogs(appConf.Logs)
	if errSetupLogs != nil && verbose {
		_, _ = yellowColor.Println(fmt.Sprintf("Unable to setup log file. Reason: %s", errSetupLogs.Error()))
	}

	var progressWrapper = ui.InitMultiProgress(threads, count)
	for i := 0; i < threads; i++ {
		go ThreadStart(i, url, progressWrapper, appConf.Template.Size)
	}

	progressWrapper.WaitForCompletion()
}

func ThreadStart(threadID int, url string, progressWrapper *ui.ProgressWrapper, linesCount int) {
	progressWrapper.AddBar(threadID)
	defer progressWrapper.DoneExecution()

	var maxExecutionEndTime = time.Now().Add(time.Second * time.Duration(maxDuration))
	util.InfoLog(fmt.Sprintf("Determined maximum execution duration time: %#v for thread with id: %d", maxExecutionEndTime, threadID), appConf.Logs)

	for i := 0; i < count; i++ {
		var requestStartTime = time.Now()
		var getUrl string
		if appConf.Template.Enabled && linesCount > 0 {
			var lineNum, templateLine = util.ReadRandomLine(template, linesCount)
			util.InfoLog(fmt.Sprintf("Received template line %s from the line %d", templateLine, lineNum), appConf.Logs)
			var updatedUrl, errPrepareUrl = util.PrepareUrl(url, templateLine)
			if errPrepareUrl != nil {
				util.ErrorLog("Can not make GET request with broken URL. Skipping this iteration", appConf.Logs)
				progressWrapper.Increment(threadID, time.Since(requestStartTime))
				continue
			}

			util.InfoLog(fmt.Sprintf("Updated URL address ccording to template file %s with line number %d. WAS: %s BECOME: %s", template, lineNum, url, updatedUrl), appConf.Logs)
			getUrl = updatedUrl
		} else {
			getUrl = url
		}

		var getResponse, getResponseErr = http.Get(getUrl)
		if getResponseErr == nil {
			util.WarnLog(fmt.Sprintf("Received HTTP GET response with status code: %d from address '%s' with ContentLength: %d", getResponse.StatusCode, getUrl, getResponse.ContentLength), appConf.Logs)
			_ = getResponse.Body.Close()
		} else {
			util.ErrorLog(fmt.Sprintf("Received an error on HTTP GET request from address: '%s' with message: %s", getUrl, getResponseErr.Error()), appConf.Logs)
		}

		progressWrapper.Increment(threadID, time.Since(requestStartTime))

		if sleepMs > 0 {
			util.InfoLog(fmt.Sprintf("Sleeping thread with id: %d for %d millis before the next itteration", threadID, sleepMs), appConf.Logs)
			time.Sleep(time.Millisecond * time.Duration(sleepMs))
			util.InfoLog(fmt.Sprintf("Resumed thread with id: %d after sleeping for %d millis", threadID, sleepMs), appConf.Logs)
		}

		if maxDuration != 0 && maxExecutionEndTime.Before(time.Now()) {
			util.WarnLog(fmt.Sprintf("Exceeded maximum execution duration of %d second(s). Terminating execution of thread with id: %d as it did not complete before time: %s", maxDuration, threadID, maxExecutionEndTime.Format(time.RFC3339)), appConf.Logs)
			progressWrapper.CompleteProgress(threadID)
			break
		}
	}
}

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

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().IntVarP(&threads, "threads", "t", 1, "A number of concurrent GET requests")
	getCmd.Flags().IntVarP(&count, "count", "c", 1, "A number of GET requests per single thread")
	getCmd.Flags().IntVarP(&sleepMs, "sleep", "s", 0, "A delay in millis after each GET requests. Doesn't impact performance report results if set (default 0)")
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

	if len(validationResults) > 0 {
		fmt.Println(redColor.Sprintf(validationResults))
		os.Exit(1)
	}
}

func Run(url string) {
	var waitGroup sync.WaitGroup
	var executionResults = make([]int, threads)

	var multiProgress = mpb.New(mpb.WithWaitGroup(&waitGroup))
	waitGroup.Add(threads)

	for i := 0; i < len(executionResults); i++ {
		go ThreadStart(i, &url, executionResults, multiProgress, &waitGroup)
	}
	//AwaitForCompletion(executionResults)
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
				decor.EwmaETA(decor.ET_STYLE_GO, 10), greenColor.Sprintf("COMPLETE"),
			),
		),
	)

	defer waitGroup.Done()
	for i := 0; i < count; i++ {
		var requestStartTime = time.Now()
		var _, _ = http.Get(*url)
		progress.Increment(time.Since(requestStartTime))

		if sleepMs > 0 {
			time.Sleep(time.Millisecond * time.Duration(sleepMs))
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

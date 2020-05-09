package ui

import (
	"github.com/fatih/color"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
	"sync"
	"time"
)

var yellowColor = color.New(color.FgYellow)
var greenColor = color.New(color.FgGreen)

// InitMultiProgress creates WaitGroup with optional progress components and fills
// ProgressWrapper struct with created data
func InitMultiProgress(threads int, count int) *ProgressWrapper {
	var multiProgress *mpb.Progress
	var waitGroup = &sync.WaitGroup{}
	var silent = true
	waitGroup.Add(threads)

	if count > 9 {
		multiProgress = mpb.New(mpb.WithWaitGroup(waitGroup))
		silent = false
	}

	return &ProgressWrapper{
		multiProgress: multiProgress,
		progressBars:  make([]*mpb.Bar, threads),
		waitGroup:     waitGroup,
		threads:       threads,
		count:         count,
		silent:        silent,
	}
}

// AddBar creates *mpb.Bar and adds newly created progress component to progressBars slice.
// A threadID determines an ID of a Thread and also slice's index.
// This scenario applicable if pw.multiProgress initialized and progress mode not silent
func (pw *ProgressWrapper) AddBar(threadID int) {
	if pw.multiProgress != nil && !pw.silent {
		var threadDescription = yellowColor.Sprintf("Thread #%-4d", threadID)
		var onCompleteDecorator = decor.OnComplete(
			decor.EwmaETA(decor.ET_STYLE_GO, 10),
			greenColor.Sprintf("DONE"),
		)

		var progress = pw.multiProgress.AddBar(int64(pw.count),
			mpb.BarStyle("╢▌▌ ╟"),
			mpb.PrependDecorators(
				decor.Name(threadDescription),
				decor.Percentage(decor.WCSyncSpace),
			),
			mpb.AppendDecorators(onCompleteDecorator),
		)

		pw.progressBars[threadID] = progress
	}
}

// Increment increments a value of progress bar with ID on 1 item.
// This scenario applicable if pw.multiProgress initialized, pw.progressBars[threadID] is present
// and progress mode not silent
func (pw *ProgressWrapper) Increment(threadID int, timeSince time.Duration) {
	if pw.multiProgress != nil && !pw.silent && pw.progressBars[threadID] != nil {
		pw.progressBars[threadID].IncrBy(1, timeSince)
	}
}

// DoneExecution decrements the WaitGroup counter by one
func (pw *ProgressWrapper) DoneExecution() {
	pw.waitGroup.Done()
}

// CompleteProgress sets complete flag of a progress component with index at threadID to true, to trigger bar complete event now
// This scenario applicable if pw.multiProgress initialized, pw.progressBars[threadID] is present
// and progress mode not silent
func (pw *ProgressWrapper) CompleteProgress(threadID int) {
	if pw.multiProgress != nil && !pw.silent && pw.progressBars[threadID] != nil {
		pw.progressBars[threadID].SetTotal(pw.progressBars[threadID].Current(), true)
	}
}

// WaitForCompletion:
// 1. Blocks until the WaitGroup counter is zero in silent mode.
// OR
// 2. Wait first waits for user provided *sync.WaitGroup, if any, then
// waits far all bars to complete and finally shutdowns master goroutine.
func (pw *ProgressWrapper) WaitForCompletion() {
	if pw.multiProgress == nil && pw.silent {
		pw.waitGroup.Wait()
	} else {
		pw.multiProgress.Wait()
	}
}

// ProgressWrapper structure contains information about execution mode (silent: true/false),
// Progress Bars and WaitGroups and other related info
type ProgressWrapper struct {
	multiProgress *mpb.Progress
	progressBars  []*mpb.Bar
	waitGroup     *sync.WaitGroup
	threads       int
	count         int
	silent        bool
}

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

func (pw *ProgressWrapper) Increment(threadID int, timeSince time.Duration) {
	if pw.multiProgress != nil && !pw.silent && pw.progressBars[threadID] != nil {
		pw.progressBars[threadID].IncrBy(1, timeSince)
	}
}
func (pw *ProgressWrapper) DoneExecution() {
	pw.waitGroup.Done()
}

func (pw *ProgressWrapper) CompleteProgress(threadID int) {
	if pw.multiProgress != nil && !pw.silent && pw.progressBars[threadID] != nil {
		pw.progressBars[threadID].SetTotal(pw.progressBars[threadID].Current(), true)
	}
}

func (pw *ProgressWrapper) WaitForCompletion() {
	if pw.multiProgress == nil && pw.silent {
		pw.waitGroup.Wait()
	} else {
		pw.multiProgress.Wait()
	}
}

type ProgressWrapper struct {
	multiProgress *mpb.Progress
	progressBars  []*mpb.Bar
	waitGroup     *sync.WaitGroup
	threads       int
	count         int
	silent        bool
}

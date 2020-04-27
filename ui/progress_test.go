package ui

import (
	"testing"
	"time"
)

func TestInitMultiProgressInSilentMode(t *testing.T) {
	var pw = InitMultiProgress(10, 5)

	if !pw.silent {
		t.Error("ProgressWrapper should be in silent mode")
	}
}

func TestInitMultiProgressInSilentModeWithoutMultiProgress(t *testing.T) {
	var pw = InitMultiProgress(10, 5)

	if pw.multiProgress != nil {
		t.Error("ProgressWrapper in silent mode should not have multiProgress")
	}
}

func TestInitMultiProgressInSilentModeWithoutMultiProgressBars(t *testing.T) {
	var pw = InitMultiProgress(10, 5)

	if len(pw.progressBars) == 0 {
		t.Error("ProgressWrapper in silent mode should not have progressBars")
	}
}

func TestInitMultiProgressInSilentModeWithWaitGroup(t *testing.T) {
	var pw = InitMultiProgress(10, 5)

	if pw.waitGroup == nil {
		t.Error("ProgressWrapper in silent mode should not have waitGroup")
	}
}

func TestInitMultiProgress(t *testing.T) {
	var pw = InitMultiProgress(10, 10)

	if pw.silent {
		t.Error("ProgressWrapper should be not silent if count > 10")
	}
}

func TestInitMultiProgressWithMultiProgress(t *testing.T) {
	var pw = InitMultiProgress(10, 10)

	if pw.multiProgress == nil {
		t.Error("ProgressWrapper should have multiProgress")
	}
}

func TestInitMultiProgressWithMultiProgressBars(t *testing.T) {
	var pw = InitMultiProgress(10, 10)

	if pw.progressBars == nil || len(pw.progressBars) != 10 {
		t.Error("ProgressWrapper should have progressBars")
	}

	for _, bar := range pw.progressBars {
		if bar != nil {
			t.Error("ProgressWrapper should have empty progressBars")
		}
	}
}

func TestInitMultiProgressWithWaitGroup(t *testing.T) {
	var pw = InitMultiProgress(10, 10)

	if pw.waitGroup == nil {
		t.Error("ProgressWrapper should have waitGroup")
	}
}

func TestProgressWrapper_AddBar(t *testing.T) {
	var pw = InitMultiProgress(10, 10)
	var givenThreadID = 5

	if pw.progressBars[givenThreadID] != nil {
		t.Error("ProgressWrapper should have empty progressBars")
	}

	pw.AddBar(givenThreadID)

	for i, bar := range pw.progressBars {
		if i == givenThreadID && bar == nil {
			t.Error("ProgressWrapper should not have empty progressBar after calling AddBar")
		} else if i != givenThreadID && bar != nil {
			t.Error("ProgressWrapper should have empty progressBars")
		}
	}
}

func TestProgressWrapper_ShouldNotAddBarInSilentMode(t *testing.T) {
	var pw = InitMultiProgress(10, 5)
	var givenThreadID = 2

	pw.AddBar(givenThreadID)

	for _, bar := range pw.progressBars {
		if bar != nil {
			t.Error("ProgressWrapper should not have empty progressBars after calling AddBar in silent mode")
		}
	}
}

func TestProgressWrapper_Increment(t *testing.T) {
	var pw = InitMultiProgress(10, 10)
	var givenThreadID = 5
	pw.AddBar(givenThreadID)

	if pw.progressBars[givenThreadID].Current() != 0 {
		t.Error("progressBar should not be incremented from the beginning")
	}

	pw.Increment(givenThreadID, time.Since(time.Now()))

	if pw.progressBars[givenThreadID].Current() < 1 {
		t.Error("progressBar should be incremented")
	}
}

func TestProgressWrapper_ShouldNotIncrementInSilentMode(t *testing.T) {
	var pw = InitMultiProgress(10, 5)
	var givenThreadID = 5
	pw.AddBar(givenThreadID)
	pw.Increment(givenThreadID, time.Since(time.Now()))
}

func TestProgressWrapper_DoneExecution(t *testing.T) {
	var pw = InitMultiProgress(10, 10)
	var givenThreadID = 5
	pw.AddBar(givenThreadID)

	pw.DoneExecution()
}

func TestProgressWrapper_WaitForCompletionWithProgressBars(t *testing.T) {
	var pw = InitMultiProgress(10, 10)
	for i := 0; i < 10; i++ {
		i := i
		go func() {
			pw.AddBar(i)
			defer pw.DoneExecution()
			time.Sleep(time.Millisecond * time.Duration(10))
			pw.CompleteProgress(i)
		}()
	}

	pw.WaitForCompletion()
}

func TestProgressWrapper_WaitForCompletion(t *testing.T) {
	var pw = InitMultiProgress(3, 3)
	for i := 0; i < 3; i++ {
		i := i
		go func() {
			pw.AddBar(i)
			defer pw.DoneExecution()
			time.Sleep(time.Millisecond * time.Duration(10))
		}()
	}

	pw.WaitForCompletion()
}

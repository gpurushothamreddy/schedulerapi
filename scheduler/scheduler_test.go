package scheduler

import (
	"testing"
	"time"

	"github.com/schedulerapi/scheduler/storage"
	"github.com/schedulerapi/scheduler/task"
)

const TestTaskName = "github.com/schedulerapi/scheduler/task.(*CallbackMock).CallNoArgs-fm"

func TestRunAt(t *testing.T) {
	mock := task.CallbackMock{}

	timeNow := time.Now()
	scheduler := New(storage.NewMemoryStorage())
	taskID, err := scheduler.RunAt(timeNow, mock.CallNoArgs)
	if err != nil {
		t.Error("Creating a task should succeed")
	}

	_, err = scheduler.RunAt(timeNow, "InvalidFunction")
	if err == nil {
		t.Error("InvalidFunction should have failed RunAt")
	}

	if len(scheduler.tasks) > 1 {
		t.Error("There should only be one task")
	}

	if scheduler.tasks[taskID].NextRun != timeNow {
		t.Error("The task's NextRun should be equal to passed parameter")
	}
}

func TestRunEvery(t *testing.T) {
	mock := task.CallbackMock{}
	scheduler := New(storage.NewMemoryStorage())
	taskID, err := scheduler.RunEvery(5, mock.CallNoArgs)
	if err != nil {
		t.Error("Creating a task should succeed")
	}

	_, err = scheduler.RunEvery(5, "InvalidFunction")
	if err == nil {
		t.Error("InvalidFunction should have failed RunAt")
	}

	if !scheduler.tasks[taskID].IsRecurring {
		t.Error()
	}
}

func TestRunPending(t *testing.T) {
	mock := task.CallbackMock{}
	scheduler := New(storage.NewMemoryStorage())
	_, err := scheduler.RunAt(time.Now(), mock.CallNoArgs)
	if err != nil {
		t.Error("Creating a task should succeed")
	}

	mock.On("CallNoArgs").Return()

	scheduler.runPending()

	time.Sleep(100 * time.Millisecond)
	mock.AssertExpectations(t)

	if len(scheduler.tasks) > 0 {
		t.Error("Non-recurring task should be removed once executed")
	}

	// Test again with a recurring task
	_, _ = scheduler.RunEvery(5, mock.CallNoArgs)

	mock.On("CallNoArgs").Return()

	// Task should be executed and then rescheduled
	scheduler.runPending()
	time.Sleep(100 * time.Millisecond)
	mock.AssertExpectations(t)
	if len(scheduler.tasks) == 0 {
		t.Error("The recurring task should still exist")
	}
}

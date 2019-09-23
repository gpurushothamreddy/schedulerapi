package scheduler

import (
	"testing"

	"github.com/schedulerapi/scheduler/storage"
	"github.com/schedulerapi/scheduler/task"
)

func TestStore(t *testing.T) {
	mock := task.CallbackMock{}
	funcRegistry := task.NewFuncRegistry()
	store := getStoreBridge(funcRegistry, nil)
	task := newTask(funcRegistry, mock.CallNoArgs)
	task.IsRecurring = true
	err := store.Add(task)
	if err != nil {
		t.Error("Failed to store task")
	}
}

func TestStoreTaskWithMultipleParams(t *testing.T) {
	mock := task.CallbackMock{}
	funcRegistry := task.NewFuncRegistry()
	store := getStoreBridge(funcRegistry, nil)
	task := newTask(funcRegistry, mock.CallWithArgs, "Hello", "World")
	err := store.Add(task)
	if err != nil {
		t.Error("Failed to store task with multiple params")
	}
}

func TestStoreThatFails(t *testing.T) {
	mock := task.CallbackMock{}
	funcRegistry := task.NewFuncRegistry()
	store := getStoreBridge(funcRegistry, nil)
	task := newTask(funcRegistry, mock.CallWithChan, make(chan bool))
	err := store.Add(task)
	if err == nil {
		t.Error("Wrong storage of a task with a channel arg took place")
	}
}

func TestRemoveTask(t *testing.T) {
	mock := task.CallbackMock{}
	funcRegistry := task.NewFuncRegistry()
	store := getStoreBridge(funcRegistry, nil)
	task := newTask(funcRegistry, mock.CallWithArgs, "Hello", "World")
	_ = store.Add(task)
	err := store.Remove(task)
	if err != nil {
		t.Error("Failed to remove task")
	}
}

func TestRemoveThatFails(t *testing.T) {
	mock := task.CallbackMock{}
	funcRegistry := task.NewFuncRegistry()
	store := getStoreBridge(funcRegistry, nil)
	task := newTask(funcRegistry, mock.CallWithChan, make(chan bool))
	err := store.Remove(task)
	if err == nil {
		t.Error("Wrong call to remove a task with a channel arg took place")
	}
}

func TestFetch(t *testing.T) {
	mock := task.CallbackMock{}
	funcRegistry := task.NewFuncRegistry()
	store := getStoreBridge(funcRegistry, nil)
	task := newTask(funcRegistry, mock.CallNoArgs)
	err := store.Add(task)
	if err != nil {
		t.Error("Failed to store task")
	}
	tasks, err := store.Fetch()
	if err != nil {
		t.Error("Could not read tasks from store")
	}

	if len(tasks) != 1 {
		t.Error("Found wrong task count")
	}
}

func TestFetchWithParams(t *testing.T) {
	mock := task.CallbackMock{}
	funcRegistry := task.NewFuncRegistry()
	store := getStoreBridge(funcRegistry, nil)
	task := newTask(funcRegistry, mock.CallWithArgs, "Test", true)
	err := store.Add(task)
	if err != nil {
		t.Error("Failed to store task")
	}
	tasks, err := store.Fetch()
	if err != nil {
		t.Error("Could not read tasks from store")
	}

	if len(tasks) != 1 {
		t.Error("Found wrong task count")
	}
}

func newTask(funcRegistry *task.FuncRegistry, function task.Function, params ...task.Param) *task.Task {
	funcMeta, err := funcRegistry.Add(function)
	if err != nil {
		return nil
	}

	task := task.New(funcMeta, params)
	_, _ = funcRegistry.Add(task.Func)

	return task
}

func getStoreBridge(funcRegistry *task.FuncRegistry, store storage.TaskStore) storeBridge {
	if store == nil {
		store = storage.NewMemoryStorage()
	}
	storeBridge := storeBridge{
		store:        store,
		funcRegistry: funcRegistry,
	}
	return storeBridge
}

package httpserver

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/schedulerapi/scheduler"
)

type getSchedulerHandler struct {
	scheduler *scheduler.Scheduler
}

// NewScheduleHandler to configure Schedule endpoint
func NewScheduleHandler(scheduler *scheduler.Scheduler) http.Handler {
	return &getSchedulerHandler{scheduler}
}

func (h *getSchedulerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var schedule ScheduleInput
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&schedule)
	if err != nil || schedule.Command == "" || schedule.ScheduleDateTime == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ErrorMessage{ErrorCode: "INVALID_BODY", Message: "Invalid body"})
		return
	}

	pst, _ := time.LoadLocation("America/Los_Angeles")
	var longForm = "2006-01-02 15:04:05 PST"
	s, err := time.ParseInLocation(longForm, schedule.ScheduleDateTime, pst)

	if err != nil || s.Before(time.Now()) {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ErrorMessage{ErrorCode: "INVALID_BODY", Message: "Invalid body"})
		return
	}
	h.scheduler.RunAt(s, TaskWithArgs, schedule.Command)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
}

// TaskWithArgs to execute given bash command
func TaskWithArgs(script string) {
	out, _ := exec.Command("sh", "-c", script).Output()
	output := string(out[:])
	log.Println(output)
	log.Println("Task executed successfully")
}

// ScheduleInput struct used to deserialize request body
type ScheduleInput struct {
	Command          string `json:"command"`
	ScheduleDateTime string `json:"scheduleDateTime"`
}

// ErrorMessage to write error message in response body
type ErrorMessage struct {
	ErrorCode string `json:"errorCode"`
	Message   string `json:"message"`
}

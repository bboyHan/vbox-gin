package timer

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var job = mockJob{}

type mockJob struct{}

func (job mockJob) Run() {
	mockFunc()
}

func mockFunc() {
	time.Sleep(time.Second)
	fmt.Println("1s...")
	log.Printf("1s...")
}

func TestNewTimerTask(t *testing.T) {
	tm := NewTimerTask()
	_tm := tm.(*Task)

	{
		_, err := tm.AddTaskByFunc("func", "@every 5s", mockFunc)
		assert.Nil(t, err)
		_, ok := _tm.TaskList["func"]
		if !ok {
			t.Error("no find func")
		}
	}

	{
		_, err := tm.AddTaskByJob("job", "@every 5s", job)
		assert.Nil(t, err)
		_, ok := _tm.TaskList["job"]
		if !ok {
			t.Error("no find job")
		}
	}

	{
		_, ok := tm.FindCron("func")
		if !ok {
			t.Error("no find func")
		}
		_, ok = tm.FindCron("job")
		if !ok {
			t.Error("no find job")
		}
		_, ok = tm.FindCron("none")
		if ok {
			t.Error("find none")
		}
	}
	{
		tm.Clear("func")
		_, ok := tm.FindCron("func")
		if ok {
			t.Error("find func")
		}
	}
}

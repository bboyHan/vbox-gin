package task

import (
	"github.com/flipped-aurora/gin-vue-admin/server/utils/timer"
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
}

func TestNewTimerTask1(t *testing.T) {
	tm := timer.NewTimerTask()
	_tm := tm.(*timer.Task)

	{
		_, err := tm.AddTaskByFunc("func", "@every 1s", mockFunc)
		assert.Nil(t, err)
		_t, ok := _tm.TaskList["func"]
		if !ok {
			t.Error("no find func")
		} else {
			t.Logf("222222, %v", _t)
		}
	}

	{
		_, err := tm.AddTaskByJob("job", "@every 1s", job)
		assert.Nil(t, err)
		_t, ok := _tm.TaskList["job"]
		if !ok {
			t.Error("no find job")
		} else {
			t.Logf("3333, %v", _t)
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

package timer

import (
	"sync"

	"github.com/robfig/cron/v3"
)

type Timer interface {
	AddTaskByFunc(taskName string, spec string, task func(), option ...cron.Option) (cron.EntryID, error)
	AddTaskByJob(taskName string, spec string, job interface{ Run() }, option ...cron.Option) (cron.EntryID, error)
	FindCron(taskName string) (*cron.Cron, bool)
	StartTask(taskName string)
	StopTask(taskName string)
	Remove(taskName string, id int)
	Clear(taskName string)
	Close()
}

// Task 定时任务管理
type Task struct {
	TaskList map[string]*cron.Cron
	sync.Mutex
}

// AddTaskByFunc 通过函数的方法添加任务
func (t *Task) AddTaskByFunc(taskName string, spec string, task func(), option ...cron.Option) (cron.EntryID, error) {
	t.Lock()
	defer t.Unlock()
	if _, ok := t.TaskList[taskName]; !ok {
		t.TaskList[taskName] = cron.New(option...)
	}
	id, err := t.TaskList[taskName].AddFunc(spec, task)
	t.TaskList[taskName].Start()
	return id, err
}

// AddTaskByFuncWithSeconds 通过函数的方法使用WithSeconds添加任务
func (t *Task) AddTaskByFuncWithSeconds(taskName string, spec string, task func(), option ...cron.Option) (cron.EntryID, error) {
	t.Lock()
	defer t.Unlock()
	option = append(option, cron.WithSeconds())
	if _, ok := t.TaskList[taskName]; !ok {
		t.TaskList[taskName] = cron.New(option...)
	}
	id, err := t.TaskList[taskName].AddFunc(spec, task)
	t.TaskList[taskName].Start()
	return id, err
}

// AddTaskByJob 通过接口的方法添加任务
func (t *Task) AddTaskByJob(taskName string, spec string, job interface{ Run() }, option ...cron.Option) (cron.EntryID, error) {
	t.Lock()
	defer t.Unlock()
	if _, ok := t.TaskList[taskName]; !ok {
		t.TaskList[taskName] = cron.New(option...)
	}
	id, err := t.TaskList[taskName].AddJob(spec, job)
	t.TaskList[taskName].Start()
	return id, err
}

// AddTaskByJobWithSeconds 通过接口的方法添加任务
func (t *Task) AddTaskByJobWithSeconds(taskName string, spec string, job interface{ Run() }, option ...cron.Option) (cron.EntryID, error) {
	t.Lock()
	defer t.Unlock()
	option = append(option, cron.WithSeconds())
	if _, ok := t.TaskList[taskName]; !ok {
		t.TaskList[taskName] = cron.New(option...)
	}
	id, err := t.TaskList[taskName].AddJob(spec, job)
	t.TaskList[taskName].Start()
	return id, err
}

// FindCron 获取对应taskName的cron 可能会为空
func (t *Task) FindCron(taskName string) (*cron.Cron, bool) {
	t.Lock()
	defer t.Unlock()
	v, ok := t.TaskList[taskName]
	return v, ok
}

// StartTask 开始任务
func (t *Task) StartTask(taskName string) {
	t.Lock()
	defer t.Unlock()
	if v, ok := t.TaskList[taskName]; ok {
		v.Start()
	}
}

// StopTask 停止任务
func (t *Task) StopTask(taskName string) {
	t.Lock()
	defer t.Unlock()
	if v, ok := t.TaskList[taskName]; ok {
		v.Stop()
	}
}

// Remove 从taskName 删除指定任务
func (t *Task) Remove(taskName string, id int) {
	t.Lock()
	defer t.Unlock()
	if v, ok := t.TaskList[taskName]; ok {
		v.Remove(cron.EntryID(id))
	}
}

// Clear 清除任务
func (t *Task) Clear(taskName string) {
	t.Lock()
	defer t.Unlock()
	if v, ok := t.TaskList[taskName]; ok {
		v.Stop()
		delete(t.TaskList, taskName)
	}
}

// Close 释放资源
func (t *Task) Close() {
	t.Lock()
	defer t.Unlock()
	for _, v := range t.TaskList {
		v.Stop()
	}
}

func NewTimerTask() Timer {
	return &Task{TaskList: make(map[string]*cron.Cron)}
}

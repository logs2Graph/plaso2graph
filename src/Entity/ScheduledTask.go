package Entity

import (
	"regexp"
	"strings"
)

type ScheduledTask struct {
	Application string
	Comment     string
	Trigger     string
	User        string
}

// containsScheduledTask check if a ScheduledTask is already in a slice of ScheduledTask
func containsScheduledTask(tasks []ScheduledTask, task ScheduledTask) bool {
	for _, v := range tasks {
		if v.Application == task.Application && v.Comment == task.Comment && v.Trigger == task.Trigger && v.User == task.User {
			return true
		}
	}
	return false
}

// AddScheduledTask add a ScheduledTask to a slice of ScheduledTask
func AddScheduledTask(tasks []ScheduledTask, task ScheduledTask) []ScheduledTask {
	if containsScheduledTask(tasks, task) == false {
		tasks = append(tasks, task)
	}
	return tasks
}

// GetScheduledTask return a slice of ScheduledTask from a slice of PlasoLog
func GetScheduledTasks(data []PlasoLog) []ScheduledTask {
	var res []ScheduledTask

	// Iterate over all logs
	for _, d := range data {
		switch d.DataType {
		case "windows:tasks:job":
			res = AddScheduledTask(res, NewScheduledTaskFromTask(d))
		}
	}

	return res
}

func NewScheduledTaskFromTask(task PlasoLog) ScheduledTask {
	var res ScheduledTask

	res.Application = task.Application
	res.Comment = task.Comment

	//Parse Trigger and User
	r, _ := regexp.Compile(`by: (?P<User>.*) Working directory.*Trigger type: (?P<Trigger>.*)`)
	matches := r.FindStringSubmatch(task.Message)

	res.Trigger = matches[0]

	// When User like "NT AUTHORITY\SYSTEM"
	splittedUser := strings.Split(matches[1], "\\")
	if len(splittedUser) > 1 {
		res.User = splittedUser[1]
	} else {
		res.User = matches[1]
	}

	return res
}

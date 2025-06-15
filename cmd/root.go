package cmd

import (
	"encoding/binary"
	"os"
	"time"

	"github.com/spf13/cobra"
)

type TaskStatus int

const (
	StatusPending TaskStatus = iota
	StatusDone
	StatusDeleted
)

var statusName = map[TaskStatus]string{
	StatusPending: "pending",
	StatusDone:    "done",
	StatusDeleted: "deleted",
}

type Task struct {
	Name      string     `json:"name"`
	Status    TaskStatus `json:"status"`
	CreatedAt string     `json:"createdAt"`
}

func createTask(taskName string) *Task {
	currentDate := time.Now().Format("2006-01-02")
	newTask := Task{
		Name:      taskName,
		Status:    StatusPending,
		CreatedAt: currentDate,
	}

	return &newTask
}

var rootCmd = &cobra.Command{
	Use:   "task",
	Short: "simple task manager",
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

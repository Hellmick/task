package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

func init() {
	rootCmd.AddCommand(doCmd)
}

func doTask(taskId int, dbPath string) error {
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}

	db.Update(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("Tasks"))
		v := b.Get([]byte(itob(taskId)))

		if v != nil {
			var task Task
			if err := json.Unmarshal(v, &task); err != nil {
				return err
			}

			task.Status = StatusDone
			taskJSON, err := json.Marshal(task)
			if err != nil {
				return err
			}

			err = b.Put([]byte(itob(taskId)), taskJSON)
			if err != nil {
				return err
			}
		} else {
			fmt.Printf("Task %s doesn't exist", strconv.Itoa(taskId))
		}

		return nil
	})
	return nil
}

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "mark task as done",
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Improper task id, please use numbers to refer to tasks", err)
			return
		}
		err = doTask(id, "./db")
		if err != nil {
			panic(err)
		}
	},
}

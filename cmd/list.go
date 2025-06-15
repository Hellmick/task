package cmd

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

func listTasks(dbPath string, listAll bool) error {
	fmt.Printf("Current tasks:\n")
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}

	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("Tasks"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var task Task
			if err := json.Unmarshal(v, &task); err != nil {
				return err
			}

			if listAll {
				id := strconv.Itoa(int(binary.BigEndian.Uint64(k)))
				fmt.Printf("%s. Task: %s, Status: %s, Created at: %s\n", id, task.Name, statusName[task.Status], task.CreatedAt)
			}
			if !listAll && task.Status == StatusPending {
				fmt.Printf("Task: %s, Status: %s, Created at: %s\n", task.Name, statusName[task.Status], task.CreatedAt)
			}

		}

		return nil
	})

	defer db.Close()
	return nil
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List current tasks",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		listTasks("./db", false)
	},
}

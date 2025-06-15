package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

func addTask(task Task, dbPath string) error {
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("Tasks"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		id, _ := b.NextSequence()
		buf, err := json.Marshal(task)
		if err != nil {
			return err
		}

		err = b.Put(itob(int(id)), []byte(buf))
		return err
	})
	if err != nil {
		return err
	}

	defer db.Close()
	return nil
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add task to the list",
	Run: func(cmd *cobra.Command, args []string) {
		task := *createTask(args[0])
		err := addTask(task, "./db")
		if err != nil {
			log.Fatal(err)
		}
	},
}

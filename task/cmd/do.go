package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/angusgmorrison/gophercises/task/db"
	"github.com/spf13/cobra"
)

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Marks a task as complete.",
	Run: func(cmd *cobra.Command, args []string) {
		var ids []int
		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Println("Failed to parse the argument:", arg)
			}
			ids = append(ids, id)
		}
		tasks, err := db.AllTasks()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Something went wrong: %v", err)
			return
		}
		for _, id := range ids {
			if id <= 0 || id > len(tasks) {
				fmt.Fprintf(os.Stderr, "Invalid task number: %d.\n", id)
				continue
			}
			task := tasks[id-1]
			if err := db.DeleteTask(task.Key); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to mark task %d as completed. Error: %v\n", id, err)
			} else {
				fmt.Printf("Marked task %d as completed.\n", id)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(doCmd)
}

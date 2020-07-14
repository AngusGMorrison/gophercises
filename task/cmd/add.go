package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/angusgmorrison/gophercises/task/db"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a task to your task list.",
	Run: func(cmd *cobra.Command, args []string) {
		task := strings.Join(args, " ")
		_, err := db.CreateTask(task)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Something went wrong: %v", err)
			return
		}
		fmt.Printf("Added %q to your task list.\n", task)
	},
}

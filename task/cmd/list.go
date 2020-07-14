package cmd

import (
	"fmt"
	"os"

	"github.com/angusgmorrison/gophercises/task/db"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all of your tasks.",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := db.AllTasks()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Something went wrong: %v", err)
			os.Exit(1)
		}
		if len(tasks) == 0 {
			fmt.Println("You have no tasks to complete. Why not take a vacation? 🏖")
			return
		}
		fmt.Println("You have the following tasks:")
		for i, t := range tasks {
			fmt.Printf("%d. %s\n", i+1, t.Value)
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}

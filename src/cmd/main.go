package main

import (
	"fmt"
	"os"

	"github.com/razvvan/qi/src/internal/aws"
	"github.com/spf13/cobra"
)

const profilePrefix = "qi-"

func main() {
	rootCmd := &cobra.Command{
		Use: "qi",
	}

	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Authentication commands",
		Run: func(cmd *cobra.Command, args []string) {
			env, err := cmd.Flags().GetString("env")
			if err != nil {
				handleErr(err)
			}

			err = aws.GenerateNewSessionCredentials(env, profilePrefix)
			if err != nil {
				handleErr(err)
			}
		},
	}

	rootCmd.AddCommand(authCmd)

	rootCmd.PersistentFlags().StringP("env", "e", "dev", "(default: dev)")

	err := rootCmd.Execute()
	if err != nil {
		handleErr(err)
	}
}

func handleErr(err error) {
	fmt.Println(err)
	os.Exit(1)
}

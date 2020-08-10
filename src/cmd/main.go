package main

import (
	"fmt"
	"os"

	"github.com/razvvan/qi/src/internal/aws"
	"github.com/spf13/cobra"
)

const (
	profilePrefix = "qi-"
	colorReset    = "\033[0m"
	colorGreen    = "\033[32m"
	colorYellow   = "\033[33m"
)

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

			fmt.Println(colorGreen, " Success!", colorReset)
			fmt.Println(colorYellow, " Start using your new profile:")
			fmt.Println()
			fmt.Println(colorReset, "   $ export AWS_PROFILE="+profilePrefix+env+"-mfa\n", colorReset)
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

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/henrywallace/homelab/go/notify/gmail"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "netwatch",
	Short: "Watch for activity on a LAN",
	RunE:  main,
}

func init() {
	rootCmd.Flags().StringP("type", "t", "gmail", "kind of notifation to send")
	rootCmd.Flags().StringP("subject", "s", "", "subject of the notification")
	rootCmd.Flags().StringP("body", "b", "", "body of the notification")
}

func main(cmd *cobra.Command, args []string) error {
	subject := mustString(cmd, "subject")
	body := mustString(cmd, "body")
	gmail.Run(subject, body)
	return nil
}

// Execute adds all child commands to the root command and sets flags
// appropriately. This is called by main.main(). It only needs to happen once
// to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func mustString(cmd *cobra.Command, name string) string {
	val, err := cmd.Flags().GetString(name)
	if err != nil {
		log.Fatal(err)
	}
	return val
}

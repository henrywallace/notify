package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/henrywallace/homelab/go/notify/gmail"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "notify",
	Short: "Send notifications, e.g. gmail emails",
	RunE:  main,
}

func init() {
	rootCmd.Flags().Bool("setup", false, "whether to just do setup")
	rootCmd.Flags().Bool("qr", false, "whether to print qr code for setup url")
	rootCmd.Flags().StringP("type", "t", "gmail", "kind of notifation to send")
	rootCmd.Flags().StringP("subject", "s", "", "subject of the notification")
	rootCmd.Flags().StringP("body", "b", "", "body of the notification")
}

func main(cmd *cobra.Command, args []string) error {
	subject := mustString(cmd, "subject")
	body := mustString(cmd, "body")
	setup := mustBool(cmd, "setup")
	qr := mustBool(cmd, "qr")
	gmail.Run(setup, subject, body, qr)
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

func mustBool(cmd *cobra.Command, name string) bool {
	val, err := cmd.Flags().GetBool(name)
	if err != nil {
		log.Fatal(err)
	}
	return val
}

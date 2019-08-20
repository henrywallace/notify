package cmd

import (
	"fmt"
	"os"

	"github.com/henrywallace/homelab/go/notify/gmail"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "netwatch",
	Short: "Watch for activity on a LAN",
	RunE:  main,
}

func init() {
	rootCmd.Flags().StringP("type", "t", "gmail", "kind of notifation to send")
	rootCmd.Flags().StringP("subject", "s", "config.toml", "toml file to trigger config")
	rootCmd.Flags().StringSliceP("body", "b", nil, "config trigger names to only run")
}

func main(cmd *cobra.Command, args []string) error {
	gmail.Run()
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

func mustString(
	log *logrus.Logger,
	cmd *cobra.Command,
	name string,
) string {
	val, err := cmd.Flags().GetString(name)
	if err != nil {
		log.Fatal(err)
	}
	return val
}

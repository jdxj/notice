/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/jdxj/notice/config"

	"github.com/spf13/cobra"
)

// emailCmd represents the email command
var emailCmd = &cobra.Command{
	Use:   "email",
	Short: "set email config",
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.SetEmail(emailFlags); err != nil {
			fmt.Fprintf(os.Stderr, "set email config failed: %s", err)
		} else {
			fmt.Printf("set email success")
		}
	},
}

var emailFlags = &config.Email{}

func init() {
	rootCmd.AddCommand(emailCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// emailCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// emailCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	emailCmd.Flags().StringVarP(&emailFlags.Addr, "addr", "a",
		"985759262@qq.com", "set email service account")
	emailCmd.Flags().StringVarP(&emailFlags.Token, "token", "t",
		"", "set email service password")
}

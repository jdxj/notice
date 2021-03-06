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

// neoCmd represents the neo command
var neoCmd = &cobra.Command{
	Use:   "neo",
	Short: "set neo config",
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.SetNeo(neoFlags); err != nil {
			fmt.Fprintf(os.Stderr, "set neo config failed: %s\n", err)
		} else {
			fmt.Printf("set neo config success\n")
		}
	},
}

var neoFlags = &config.Neo{}

func init() {
	rootCmd.AddCommand(neoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// neoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// neoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	neoCmd.Flags().StringVar(&neoFlags.Host, "host", "", "neo site")
	neoCmd.Flags().StringVarP(&neoFlags.Domain, "domain", "d", "", "neo domain")
	neoCmd.Flags().StringVarP(&neoFlags.Cookies, "cookies", "c", "", "neo cookies")
	neoCmd.Flags().StringVarP(&neoFlags.Services, "services", "s", "", "neo services")
	neoCmd.Flags().StringVarP(&neoFlags.User, "user", "u", "", "neo email")
}

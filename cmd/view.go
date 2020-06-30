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

	"github.com/dgraph-io/badger/v2"

	"github.com/jdxj/notice/config"

	"github.com/spf13/cobra"
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "view task config",
	Run: func(cmd *cobra.Command, args []string) {
		printConfig(vf)
	},
}

var (
	vf = &viewFlags{}
)

type viewFlags struct {
	bilibiliFlag    bool
	neoFlag         bool
	sourceforgeFlag bool
	emailFlag       bool
	ryfFlag         bool
}

func init() {
	rootCmd.AddCommand(viewCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// viewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// viewCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	viewCmd.Flags().BoolVarP(&vf.bilibiliFlag, "bilibili", "b", true, "view BiliBili's config")
	viewCmd.Flags().BoolVarP(&vf.neoFlag, "neo", "n", true, "view neo_proxy's config")
	viewCmd.Flags().BoolVarP(&vf.sourceforgeFlag, "sourceforge", "s", true, "view sourceforge config")
	viewCmd.Flags().BoolVarP(&vf.emailFlag, "email", "e", true, "view email config")
	viewCmd.Flags().BoolVarP(&vf.ryfFlag, "ryf", "r", true, "view ruanyifeng config")
}

func printConfig(vf *viewFlags) {
	var result string
	defer func() {
		fmt.Printf("%s", result)
	}()

	if vf.bilibiliFlag {
		b, err := config.GetBiliBili()
		if err != nil {
			if err != badger.ErrKeyNotFound {
				fmt.Fprintf(os.Stderr, "%s\n", err)
			}
		} else {
			result = fmt.Sprintf("%s\nbilibili config:\n%s\n", result, b)
		}
	}

	if vf.neoFlag {
		n, err := config.GetNeo()
		if err != nil {
			if err != badger.ErrKeyNotFound {
				fmt.Fprintf(os.Stderr, "%s\n", err)
			}
		} else {
			result = fmt.Sprintf("%s\nneo config:\n%s\n", result, n)
		}
	}

	if vf.sourceforgeFlag {
		s, err := config.GetSourceforge()
		if err != nil {
			if err != badger.ErrKeyNotFound {
				fmt.Fprintf(os.Stderr, "%s\n", err)
			}
		} else {
			result = fmt.Sprintf("%s\nsourceforge config:\n%s\n", result, s)
		}
	}

	if vf.emailFlag {
		e, err := config.GetEmail()
		if err != nil {
			if err != badger.ErrKeyNotFound {
				fmt.Fprintf(os.Stderr, "%s\n", err)
			}
		} else {
			result = fmt.Sprintf("%s\nemail config:\n%s\n", result, e)
		}
	}

	if vf.ryfFlag {
		e, err := config.GetRYF()
		if err != nil {
			if err != badger.ErrKeyNotFound {
				fmt.Fprintf(os.Stderr, "%s\n", err)
			}
		} else {
			result = fmt.Sprintf("%s\nruanyifeng config:\n%s\n", result, e)
		}
	}
}

// Copyright © 2019 James Routley <jroutley@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/jamesroutley/mmv/mmv"
	"github.com/spf13/cobra"
)

var (
	dryRun   bool
	includes string
	excludes string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mmv [directory]",
	Short: "multi move files",
	Long:  `mmv lets you rename multiple files in a directory at once, by editing their names in your favourite text editor`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var options []func(*mmv.MultiMover)
		if dryRun {
			options = append(options, mmv.OptionDryRun())
		}
		if includes != "" {
			re, err := regexp.Compile(includes)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			options = append(options, mmv.OptionInclude(re))
		}
		if excludes != "" {
			re, err := regexp.Compile(excludes)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			options = append(options, mmv.OptionExclude(re))
		}
		mover := mmv.NewMultiMover(options...)
		if err := mover.MultiMoveDir(args[0]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print out the the changes without actually making them")
	rootCmd.Flags().StringVar(&includes, "includes", "", "Only include files which match this regular expression")
	rootCmd.Flags().StringVar(&excludes, "excludes", "", "Exclude files which match this regular expression")
}

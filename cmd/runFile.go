// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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
	"github.com/dpaolella/shp2csvApp/shp2csvApp"
	"github.com/spf13/cobra"
)

// runFileCmd represents the runFile command
var runFileCmd = &cobra.Command{
	Use:   "runFile",
	Short: "Run on a single shapefile",
	Long: `runFile exports data from a shapefile to a
	a csv of the same name`,
	Run: func(cmd *cobra.Command, args []string) {
		shp2csvApp.Run(false, args[0])
	},
}

func init() {
	RootCmd.AddCommand(runFileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runFileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runFileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

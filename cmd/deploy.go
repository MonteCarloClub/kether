/*
Copyright (c) 2022 Zhang Zhanpeng <zhangregister@outlook.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"

	"github.com/MonteCarloClub/kether/flag"
	"github.com/MonteCarloClub/kether/log"
	"github.com/MonteCarloClub/kether/object"
	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var (
	dryRun   bool
	yamlPath string

	deployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.WithValue(context.Background(), flag.ContextKey, flag.ContextValType{
				DryRun: dryRun,
			})
			ketherObject, ketherObjectState, err := object.Register(ctx, yamlPath)
			if err != nil {
				log.Error("fail to register kether object", "err", err)
				return
			}
			log.Info("kether object registered")

			err = object.Deploy(ctx, ketherObject, ketherObjectState)
			if err != nil {
				log.Error("fail to deploy ketherObject", "err", err)
				return
			}
			log.Info("kether object deployed")
		},
	}
)

func init() {
	rootCmd.AddCommand(deployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	deployCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Output actions to be performed without changing any state")
	deployCmd.Flags().StringVarP(&yamlPath, "file", "f", "", "Construct Kether object and its state with this YAML file path (required)")
	deployCmd.MarkFlagRequired("file")
}

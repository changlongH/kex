/*
Copyright © 2023 hcl <hcl2685@gmail.com>
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/changlongH/kex/client"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// podsCmd represents the pods command
var podsCmd = &cobra.Command{
	Use:   "pods",
	Short: "show pods in current use cluster",
	Long: `show pods in current use cluster.
You can search and select pod to
	describe: show pod info
or
	enter:
		select container
		select cwd
or
	cmd: run cmd and gets result`,
	Run: func(cmd *cobra.Command, args []string) {
		kubeCli := client.NewClient()
		kubeCli.InitClientSet()

		ns := ""
		if len(args) > 0 {
			ns = string(args[0])
		}
		podList := kubeCli.GetPods(ns)
		cnt := len(podList.Items)

		prompt := promptui.Select{
			Label:     "Fond Pods Count: " + strconv.FormatInt(int64(cnt), 10),
			Items:     podList.Items,
			Templates: PodsTemplate,
			Size:      60,
			/*
				Searcher: func(input string, idx int) bool {
					line := podList[idx]
					return strings.Contains(line, input)
				},
			*/
		}

		_, result, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}
		fmt.Println(promptui.IconGood + " " + result)
	},
}

func init() {
	rootCmd.AddCommand(podsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// podsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//podsCmd.Flags().String(ns, "ns", "specify namespace to get pods in particular namespace")
}

/*
Copyright © 2023 hcl <hcl2685@gmail.com>
*/
package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// describeCmd represents the describe command
var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Show details of a specific resource or group of resources",
	Long: `Print a detailed description of the selected resources, including related resources such as events or controllers.
	You may select a single object by name, all objects of that type to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := exec.LookPath("kubectl")
		if err != nil {
			fmt.Println("'kubectl' not found. install pls")
			return
		}

		podcmd := exec.Command("kubectl", "get", "pods", "-A")
		out, err := podcmd.CombinedOutput()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if len(out) <= 0 {
			fmt.Println("not pods. check config pls")
			return
		}

		lines := strings.Split(string(out), "\n")
		if len(lines) <= 1 {
			fmt.Println("not pods. check config pls")
			return
		}

		template := &promptui.SelectTemplates{
			Label:    "{{ . | yellow }} " + promptui.IconInitial,
			Active:   promptui.IconSelect + " {{ . | red }}",
			Inactive: "{{ . | cyan }}",
			//Selected: "{{ . | yellow }}",
		}

		prompt := promptui.Select{
			Label:     "Search Pod Name:",
			Items:     lines,
			Templates: template,
			Size:      60,
			/*
				Searcher: func(input string, idx int) bool {
					line := lines[idx]
					if strings.Contains(line, input) {
						return true
					}
					return false
				},
			*/
		}

		_, result, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		values := strings.Fields(result)
		namespace := values[0]
		podname := values[1]

		//fmt.Printf("Describe Namespace: %s | Pod: %s %s\n", namespace, podname, promptui.IconGood)

		//kubectl describe pod ${pod} -n ${namespace}
		syscmd := exec.Command("kubectl", "describe", "pod", podname, "-n", namespace)
		out, err = syscmd.CombinedOutput()
		if err != nil {
			fmt.Printf("%v\n err:%s\n", string(out), err.Error())
			return
		}
		fmt.Println(string(out))
	},
}

func init() {
	rootCmd.AddCommand(describeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// describeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// describeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

/*
Copyright © 2023 hcl <hcl2685@gmail.com>
*/
package cmd

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// clusterCmd represents the cluster command
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Select kubernetes clusters config list",
	Long: `Select kubernetes clusters config.
	Copy your select config to $HOME/.kube/config or .kex.yaml config path`,

	Run: func(cmd *cobra.Command, args []string) {
		kubeConfigs := []string{}

		srcPath := GetKubeConfigsPath()
		if len(srcPath) <= 0 {
			log.Printf("Not set 'kubeConfigsPath' value in %v", Viper.ConfigFileUsed())
			return
		}

		files, err := os.ReadDir(srcPath)
		if err != nil {
			log.Printf("Clusters Config Dir Error:%v\n", err.Error())
			return
		}

		for _, file := range files {
			kubeConfigs = append(kubeConfigs, file.Name())
		}

		if len(kubeConfigs) <= 0 {
			log.Printf("Clusters Config Directory Is Empty\nDir: %v\n", srcPath)
			return
		}

		prompt := promptui.Select{
			Label:     "Select Cluster Config:",
			Items:     kubeConfigs,
			Templates: ClusterTemplate,
		}

		_, filename, err := prompt.Run()

		if err != nil {
			log.Printf("Prompt failed %v\n", err)
			return
		}

		src := filepath.Join(srcPath, filename)
		dst := GetKubeConfigFile()

		syscmd := exec.Command("cp", src, dst)
		out, err := syscmd.CombinedOutput()
		if err != nil {
			log.Printf("%v\ncp config failed %s\n %s\n", string(out), promptui.IconBad, err.Error())
			return
		}
		log.Printf("cp %v %v", src, dst)
		log.Println("Select SUCCESS ", promptui.IconGood)
	},
}

func init() {
	rootCmd.AddCommand(clusterCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clusterCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}

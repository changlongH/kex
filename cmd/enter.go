/*
Copyright © 2023 hcl <hcl2685@gmail.com>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"golang.org/x/term"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/util/homedir"
)

const (
	containerRunPath = "/home/game/running/server/"
	containerLogPath = "/home/game/log/server/"

	kubeconfigPath = ".kube/config"
)

var template = &promptui.SelectTemplates{
	Label:    "{{ . | green }}",
	Active:   promptui.IconSelect + " {{ . | red }}",
	Inactive: "{{ . | cyan }}",
	//Selected:  promptui.IconGood + " {{ . | yellow }}",
}

type sizeQueue chan remotecommand.TerminalSize

func (s sizeQueue) Next() *remotecommand.TerminalSize {
	size, ok := <-s
	if !ok {
		return nil
	}
	return &size
}

func kubectlExec(namespace string, podname string, bizname string, podcmd []string) {
	home := homedir.HomeDir()
	if home == "" {
		fmt.Printf("not found .kube/config in homeDir: %s\n Please select config begin\n", home)
		return
	}

	kubeconfig := filepath.Join(home, kubeconfigPath)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Printf("build config: %s err:%s\n", kubeconfig, err.Error())
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("new config err:%s\n", err.Error())
		return
	}

	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podname).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: bizname,
			Command:   podcmd,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())

	if !term.IsTerminal(0) || !term.IsTerminal(1) {
		fmt.Println("stdin/stdout should be terminal")
		return
	}

	fd := int(os.Stdin.Fd())
	// Put the terminal into raw mode to prevent it echoing characters twice.
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Println(err)
		return
	}

	termWidth, termHeight, _ := term.GetSize(fd)
	termSize := remotecommand.TerminalSize{Width: uint16(termWidth), Height: uint16(termHeight)}
	s := make(sizeQueue, 1)
	s <- termSize

	defer func() {
		err := term.Restore(fd, oldState)
		if err != nil {
			fmt.Println(err)
		}
	}()

	err = exec.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdin:             os.Stdin,
		Stdout:            os.Stdout,
		Stderr:            os.Stderr,
		Tty:               true,
		TerminalSizeQueue: s,
	})

	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Println()
}

// enterCmd represents the enterpod command
var enterCmd = &cobra.Command{
	Use:   "enterpod",
	Short: "enter select pod",
	Long:  `enter select kuber pod`,
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

		// remove status line
		lines = lines[1:]

		prompt := promptui.Select{
			Label:     "Search Pod Name:",
			Items:     lines,
			Templates: template,
			Size:      70,
			Searcher: func(input string, idx int) bool {
				line := lines[idx]
				if strings.Contains(line, input) {
					return true
				}
				return false
			},
		}

		_, result, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		values := strings.Fields(result)
		namespace := values[0]
		pod := values[1]

		// p.s.m
		psm := strings.Split(pod, "-")
		bizname := strings.Join(psm[:3], "-")
		biztype := psm[2]

		pathList := []string{}
		pathList = append(pathList, "")

		logPath := filepath.Join(containerLogPath, biztype, pod)
		pathList = append(pathList, logPath)
		runPath := filepath.Join(containerRunPath, biztype)
		pathList = append(pathList, runPath)

		prompt = promptui.Select{
			Label:     "Select enter PATH:",
			Items:     pathList,
			Templates: template,
			Size:      10,
		}

		_, path, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		// change work path
		var podCmd []string
		if len(path) == 0 {
			podCmd = []string{"bash"}
		} else {
			cwd := fmt.Sprintf("cd %s ; /bin/bash", path)
			podCmd = []string{"bash", "-c", cwd}
		}
		kubectlExec(namespace, pod, bizname, podCmd)
	},
}

func init() {
	rootCmd.AddCommand(enterCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// enterCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// enterCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

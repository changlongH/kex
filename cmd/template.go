package cmd

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
)

func GetStringDefault(key string, s string) string {
	v := viper.GetString(key)
	if len(v) > 0 {
		return v
	}
	return s
}

func GetColorDefault(key string, color string) string {
	return GetStringDefault(key, color)
}

var ClusterTemplate = &promptui.SelectTemplates{
	Label:    fmt.Sprintf("{{ . | %s }} ", GetColorDefault("colors.cluster.label", "yellow")),
	Active:   fmt.Sprintf("%s {{ . | %s }}", promptui.IconSelect, GetColorDefault("colors.cluster.active", "red")),
	Inactive: fmt.Sprintf("{{ . | %s }}", GetColorDefault("colors.cluster.inactive", "cyan")),
}

var PodsTemplate = &promptui.SelectTemplates{
	Label:    fmt.Sprintf("{{ . | %s }} ", GetColorDefault("colors.pods.label", "yellow")),
	Active:   fmt.Sprintf("%s {{ . | %s }}", promptui.IconSelect, GetColorDefault("colors.pods.active", "red")),
	Inactive: fmt.Sprintf("{{ . | %s }}", GetColorDefault("colors.pods.inactive", "cyan")),
}

package cmd

import (
	"code-reviewer/internal/config"
	"fmt"
	"os"
	"reflect"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage application configuration",
	Long:  `View and modify the application configuration (API Key, AI Model, Output Language).`,
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration values",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		v := reflect.ValueOf(*cfg)
		typeOfCfg := v.Type()

		for i := 0; i < v.NumField(); i++ {
			field := typeOfCfg.Field(i)
			value := v.Field(i).Interface()
			jsonTag := field.Tag.Get("json")
			fmt.Printf("%s: %v\n", jsonTag, value)
		}
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a specific configuration value",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		v := reflect.ValueOf(*cfg)
		typeOfCfg := v.Type()

		found := false
		for i := 0; i < v.NumField(); i++ {
			field := typeOfCfg.Field(i)
			jsonTag := field.Tag.Get("json")
			if jsonTag == key {
				fmt.Println(v.Field(i).Interface())
				found = true
				break
			}
		}

		if !found {
			fmt.Fprintf(os.Stderr, "Error: Config key '%s' not found.\n", key)
			os.Exit(1)
		}
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a specific configuration value",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		v := reflect.ValueOf(cfg).Elem()
		typeOfCfg := v.Type()

		found := false
		for i := 0; i < v.NumField(); i++ {
			field := typeOfCfg.Field(i)
			jsonTag := field.Tag.Get("json")
			if jsonTag == key {
				f := v.Field(i)
				if f.Kind() == reflect.String {
					f.SetString(value)
					found = true
				} else {
					fmt.Fprintf(os.Stderr, "Error: Unsupported field type for key '%s'\n", key)
					os.Exit(1)
				}
				break
			}
		}

		if !found {
			fmt.Fprintf(os.Stderr, "Error: Config key '%s' not found.\n", key)
			os.Exit(1)
		}

		if err := config.Save(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully set '%s' to '%s'\n", key, value)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
}

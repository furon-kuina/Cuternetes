package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	c8s "github.com/furon-kuina/cuternetes/pkg"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const c8sConfigPath = "../config.yaml"
const specsPath = "specs.yaml"

var c8sConfig c8s.C8sConfig

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Applies configuration",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Printf("apply called with argument %s\n", args[0])
		f, err := os.Open(c8sConfigPath)
		if err != nil {
			log.Fatalf("Failed to open config: %v", err)
		}
		configData, err := io.ReadAll(f)
		if err != nil {
			log.Fatalf("Failed to read config: %v", err)
		}
		var config c8s.C8sConfig
		err = yaml.Unmarshal(configData, &config)
		if err != nil {
			log.Fatalf("Failed to unmarshal config: %v", err)
		}
		fmt.Printf("Loaded config: %+v", config)
		f, err = os.Open(specsPath)
		if err != nil {
			log.Fatalf("Failed to open specs: %v", err)
		}
		specsData, err := io.ReadAll(f)
		fmt.Println(string(specsData))
		if err != nil {
			log.Fatalf("Failed to read specs: %v", err)
		}
		var specs c8s.Spec
		err = yaml.Unmarshal(specsData, &specs)
		if err != nil {
			log.Fatalf("Failed to unmarshal specs: %v", err)
		}
		fmt.Printf("Loaded specs: %+v", specs)
		err = applySpec(specs)
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// applyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// applyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func applySpec(spec c8s.Spec) (err error) {
	defer c8s.Wrap(&err, "applyspec(%q) failed: %v", spec, err)
	specData, err := json.Marshal(spec)
	if err != nil {
		return
	}
	req, err := http.NewRequest("PUT", c8sConfig.ApiServer.Url, bytes.NewBuffer([]byte(specData)))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fmt.Println("Response Status: ", resp.Status)
	fmt.Println("Response Body: ", body)
	return nil
}

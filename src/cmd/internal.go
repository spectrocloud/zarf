// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

// Package cmd contains the CLI commands for Zarf.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alecthomas/jsonschema"
	"github.com/defenseunicorns/zarf/src/internal/agent"
	"github.com/defenseunicorns/zarf/src/internal/api"
	"github.com/defenseunicorns/zarf/src/internal/cluster"
	"github.com/defenseunicorns/zarf/src/internal/packager/git"
	"github.com/defenseunicorns/zarf/src/pkg/message"
	"github.com/defenseunicorns/zarf/src/types"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var internalCmd = &cobra.Command{
	Use:     "internal",
	Aliases: []string{"dev"},
	Hidden:  true,
	Short:   "Internal tools used by zarf",
}

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Runs the zarf agent",
	Long: "NOTE: This command is a hidden command and generally shouldn't be run by a human.\n" +
		"This command starts up a http webhook that Zarf deployments use to mutate pods to conform " +
		"with the Zarf container registry and Gitea server URLs.",
	Run: func(cmd *cobra.Command, args []string) {
		agent.StartWebhook()
	},
}

var generateCLIDocs = &cobra.Command{
	Use:   "generate-cli-docs",
	Short: "Creates auto-generated markdown of all the commands for the CLI",
	Run: func(cmd *cobra.Command, args []string) {
		// Don't include the datestamp in the output
		rootCmd.DisableAutoGenTag = true
		//Generate markdown of the Zarf command (and all of its child commands)
		if err := os.RemoveAll("./docs/4-user-guide/1-the-zarf-cli/100-cli-commands"); err != nil {
			message.Fatalf("Unable to generate the CLI documentation: %s", err.Error())
		}
		if err := os.Mkdir("./docs/4-user-guide/1-the-zarf-cli/100-cli-commands", 0775); err != nil {
			message.Fatalf("Unable to generate the CLI documentation: %s", err.Error())
		}
		if err := doc.GenMarkdownTree(rootCmd, "./docs/4-user-guide/1-the-zarf-cli/100-cli-commands"); err != nil {
			message.Fatalf("Unable to generate the CLI documentation: %s", err.Error())
		} else {
			message.SuccessF("Successfully created the CLI documentation")
		}
	},
}

var configSchemaCmd = &cobra.Command{
	Use:     "config-schema",
	Aliases: []string{"c"},
	Short:   "Generates a JSON schema for the zarf.yaml configuration",
	Run: func(cmd *cobra.Command, args []string) {
		schema := jsonschema.Reflect(&types.ZarfPackage{})
		output, err := json.MarshalIndent(schema, "", "  ")
		if err != nil {
			message.Fatal(err, "Unable to generate the zarf config schema")
		}
		fmt.Print(string(output) + "\n")
	},
}

var apiSchemaCmd = &cobra.Command{
	Use:   "api-schema",
	Short: "Generates a JSON schema from the API types",
	Run: func(cmd *cobra.Command, args []string) {
		schema := jsonschema.Reflect(&types.RestAPI{})
		output, err := json.MarshalIndent(schema, "", "  ")
		if err != nil {
			message.Fatal(err, "Unable to generate the zarf api schema")
		}
		fmt.Print(string(output) + "\n")
	},
}

var createReadOnlyGiteaUser = &cobra.Command{
	Use:   "create-read-only-gitea-user",
	Short: "Creates a read-only user in Gitea",
	Long: "Creates a read-only user in Gitea by using the Gitea API. " +
		"This is called internally by the supported Gitea package component.",
	Run: func(cmd *cobra.Command, args []string) {
		// Load the state so we can get the credentials for the admin git user
		state, err := cluster.NewClusterOrDie().LoadZarfState()
		if err != nil {
			message.Error(err, "Unable to load the Zarf state")
		}

		// Create the non-admin user
		if err = git.New(state.GitServer).CreateReadOnlyUser(); err != nil {
			message.Error(err, "Unable to create a read-only user in the Gitea service.")
		}
	},
}

var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Launch the experimental Zarf UI",
	Run: func(cmd *cobra.Command, args []string) {
		api.LaunchAPIServer()
	},
}

func init() {
	rootCmd.AddCommand(internalCmd)

	internalCmd.AddCommand(agentCmd)
	internalCmd.AddCommand(generateCLIDocs)
	internalCmd.AddCommand(configSchemaCmd)
	internalCmd.AddCommand(apiSchemaCmd)
	internalCmd.AddCommand(createReadOnlyGiteaUser)
	internalCmd.AddCommand(uiCmd)
}

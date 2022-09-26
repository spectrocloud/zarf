package cmd

import (
	"github.com/defenseunicorns/zarf/src/internal/k8s"
	"github.com/spf13/cobra"
)

var (
	connectResourceName string
	connectNamespace    string
	connectResourceType string
	connectLocalPort    int
	connectRemotePort   int
	cliOnly             bool

	connectCmd = &cobra.Command{
		Use:     "connect <REGISTRY|LOGGING|GIT>",
		Aliases: []string{"c"},
		Short:   "Access services or pods deployed in the cluster.",
		Long: "Uses a k8s port-forward to connect to resources within the cluster referenced by your kube-context.\n" +
			"Three default options for this command are <REGISTRY|LOGGING|GIT>. These will connect to the Zarf created resources " +
			"(assuming they were selected when performing the `zarf init` command).\n\n" +
			"Packages can provide service manifests that define their own shortcut connection options. These options will be " +
			"printed to the terminal when the package finishes deploying.\n If you don't remember what connection shortcuts your deployed " +
			"package offers, you can search your cluster for services that have the 'zarf.dev/connect-name' label. The value of that label is " +
			"the name you will pass into the 'zarf connect' command. \n\n" +
			"Even if the packages you deploy don't define their own shortcut connection options, you can use the command flags " +
			"to connect into specific resources. You can read the command flag descriptions below to get a better idea how to connect " +
			"to whatever resource you are trying to connect to.",
		Run: func(cmd *cobra.Command, args []string) {
			var target string
			if len(args) > 0 {
				target = args[0]
			}

			tunnel := k8s.NewTunnel(connectNamespace, connectResourceType, connectResourceName, connectLocalPort, connectRemotePort)
			// If the cliOnly flag is false (default), enable auto-open
			if !cliOnly {
				tunnel.EnableAutoOpen()
			}
			tunnel.Connect(target, true)
		},
	}

	connectListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "List all available connection shortcuts.",
		Run: func(cmd *cobra.Command, args []string) {
			k8s.PrintConnectTable()
		},
	}
)

func init() {
	initViper()

	rootCmd.AddCommand(connectCmd)
	connectCmd.AddCommand(connectListCmd)

	v.SetDefault("connect.namespace", k8s.ZarfNamespace)
	v.SetDefault("connect.resource_type", k8s.SvcResource)

	connectCmd.Flags().StringVar(&connectResourceName, "name", v.GetString("connect.name"), "Specify the resource name.  E.g. name=unicorns or name=unicorn-pod-7448499f4d-b5bk6")
	connectCmd.Flags().StringVar(&connectNamespace, "namespace", v.GetString("connect.namespace"), "Specify the namespace.  E.g. namespace=default")
	connectCmd.Flags().StringVar(&connectResourceType, "type", v.GetString("connect.resource_type"), "Specify the resource type.  E.g. type=svc or type=pod")
	connectCmd.Flags().IntVar(&connectLocalPort, "local-port", v.GetInt("connect.local_port"), "(Optional, autogenerated if not provided) Specify the local port to bind to.  E.g. local-port=42000")
	connectCmd.Flags().IntVar(&connectRemotePort, "remote-port", v.GetInt("connect.remote_port"), "Specify the remote port of the resource to bind to.  E.g. remote-port=8080")
	connectCmd.Flags().BoolVar(&cliOnly, "cli-only", v.GetBool("connect.cli_only"), "Disable browser auto-open")
}

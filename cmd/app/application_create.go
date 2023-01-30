package cmd

import (
	"time"

	"k8s.io/kubectl/pkg/util/templates"
)

// RoutePollTimoutSeconds sets how long app create command waits for route host to be prepopulated
const RoutePollTimeout = 5 * time.Second

var (
	newAppLong = templates.LongDesc(`
		Create a new application by specifying source code, templates, and/or images.

		This command will try to build up the components of an application using images, templates,
		or code that has a public repository. It will look up the images on the local container storage
		(if available), a container image registry, an integrated image stream, or stored templates.
		
		If you specify a source code URL, it will set up a build that takes your source code and converts
		it into an image that can run inside of a pod. Local source must be in a git repository that has a
		remote repository that the server can see. The images will be deployed via a deployment or
		deployment configuration, and a service will be connected to the first public port of the app.
		You may either specify components using the various existing flags or let civo app create autodetect
		what kind of components you have provided.
		If you provide source code, a new build will be automatically triggered.
		You can use 'civo app status' to check the progress.`)

	newAppExample = templates.Examples(`
		# List all local templates and image streams that can be used to create an app
		civo app create --list

		# Create an application based on the source code in the current git repository (with a public remote) and a container image
		civo app create . --image=registry/repo/langimage

		# Create an application from a remote repository using its beta4 branch and tag v1
		civo app create https://github.com/openshift/ruby-hello-world#beta4:v1

		# Create an application from a remote repository and specify a context directory
		civo app create https://github.com/youruser/yourgitrepo --context-dir=src/build

		# Create an application from a remote private repository with GIT TOKEN in env variables
		civo app create https://github.com/youruser/yourgitrepo

	`)

	newAppNoInput = `You must specify one or more images, image streams, templates, or source code locations to create an application.

To list all local templates and image streams, use:

  civo app create -L

To search templates, image streams, and container images that match the arguments provided, use:

  civo app create -S php
  civo app create -S --template=rails
  civo app create -S --image-stream=mysql
  civo app create -S --image=registry.access.redhat.com/ubi8/python-38

  For details on how to use the results from those searches to provide images, image streams, templates, or source code locations as inputs into 'civo create ', use:
  civo app create help 
`
)

// var appName, appSize string
// var wait bool

// var appCreateCmd = &cobra.Command{
// 	Use:     "create",
// 	Aliases: []string{"new", "add"},
// 	Example: "civo app create APP_NAME [flags]",
// 	Short:   "Create a new application",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		utility.EnsureCurrentRegion()

// 		client, err := config.CivoAPIClient()
// 		if err != nil {
// 			utility.Error("Creating the connection to Civo's API failed with %s", err)
// 			os.Exit(1)
// 		}

// 		config, err := client.NewApplicationConfig()
// 		if err != nil {
// 			utility.Error("Unable to create a new config for the app %s", err)
// 			os.Exit(1)
// 		}

// 		if appName != "" {
// 			if utility.ValidNameLength(appName) {
// 				utility.Warning("the name cannot be longer than 63 characters")
// 				os.Exit(1)
// 			}
// 			config.Name = appName
// 		}

// 		if len(args) > 0 {
// 			if utility.ValidNameLength(args[0]) {
// 				utility.Warning("the name cannot be longer than 63 characters")
// 				os.Exit(1)
// 			}
// 			config.Name = args[0]
// 		}

// 		if appSize != "" {
// 			config.Size = appSize
// 		} else {
// 			config.Size = "small"
// 		}

// 		var executionTime string
// 		startTime := utility.StartTime()

// 		var application *civogo.Application
// 		resp, err := client.CreateApplication(config)
// 		if err != nil {
// 			utility.Error("%s", err)
// 			os.Exit(1)
// 		}

// 		if wait {
// 			stillCreating := true
// 			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
// 			s.Prefix = fmt.Sprintf("Creating application (%s)... ", resp.Name)
// 			s.Start()

// 			for stillCreating {
// 				application, err = client.FindApplication(resp.Name)
// 				if err != nil {
// 					utility.Error("%s", err)
// 					os.Exit(1)
// 				}
// 				if application.Status == "ACTIVE" {
// 					stillCreating = false
// 					s.Stop()
// 				} else {
// 					time.Sleep(2 * time.Second)
// 				}
// 			}
// 			executionTime = utility.TrackTime(startTime)
// 		} else {
// 			// we look for the created app to obtain the data that we need
// 			application, err = client.FindApplication(resp.Name)
// 			if err != nil {
// 				utility.Error("App %s", err)
// 				os.Exit(1)
// 			}
// 		}

// 		if common.OutputFormat == "human" {
// 			if executionTime != "" {
// 				fmt.Printf("The app %s has been created in %s\n", utility.Green(application.Name), executionTime)
// 			} else {
// 				fmt.Printf("The app %s has been created\n", utility.Green(application.Name))
// 			}
// 		} else {
// 			ow := utility.NewOutputWriter()
// 			ow.StartLine()
// 			ow.AppendDataWithLabel("id", resp.ID, "ID")
// 			ow.AppendDataWithLabel("name", resp.Name, "Name")
// 			ow.AppendDataWithLabel("network_id", resp.NetworkID, "Network ID")
// 			ow.AppendDataWithLabel("description", resp.Description, "Description")
// 			//ow.AppendDataWithLabel("image", resp.Image, "Image")
// 			ow.AppendDataWithLabel("size", resp.Size, "Size")
// 			ow.AppendDataWithLabel("status", resp.Status, "Status")
// 			// ow.AppendDataWithLabel("process_info", resp.ProcessInfo, "Process Info")
// 			ow.AppendDataWithLabel("domains", strings.Join(resp.Domains, ", "), "Domains")
// 			ow.AppendDataWithLabel("ssh_key_ids", strings.Join(resp.SSHKeyIDs, ", "), "SSH Key IDs")
// 			//ow.AppendDataWithLabel("config", resp.Config, "Config")

// 			if common.OutputFormat == "json" {
// 				ow.WriteSingleObjectJSON(common.PrettySet)
// 			} else {
// 				ow.WriteCustomOutput(common.OutputFields)
// 			}
// 		}
// 	},
// }

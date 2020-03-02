package cmd

// kubernetes list -- list all kubernetes clusters [ls, all]
// kubernetes versions -- list available k3s versions [version, v]
// kubernetes show ID/NAME -- show a Kubernetes cluster by ID or name [get, inspect]
// kubernetes config ID/NAME [--save] -- get or save the ~/.kube/config for a Kubernetes cluster by ID or name [kubeconfig]
// kubernetes create [NAME] [...] -- create a new kubernetes cluster with the specified name and provided options
// kubernetes rename ID/NAME [--name] -- rename Kubernetes cluster
// kubernetes upgrade ID/NAME [--version] -- upgrade Kubernetes cluster's k3s version
// kubernetes scale ID/NAME [--nodes] -- rescale the Kubernetes cluster to a new node count [rescale]
// kubernetes remove ID/NAME -- removes an entire Kubernetes cluster with ID/name entered (use with caution!) [delete, destroy, rm]
// kubernetes_applications list -- list all available kubernetes applications [ls, all]
// kubernetes_applications show NAME -- show a Kubernetes application by name [get, inspect]
// kubernetes_applications add NAME --cluster=... -- add the marketplace application to a Kubernetes cluster by ID or name

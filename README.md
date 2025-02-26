# OpenShift GitOps Cluster Bootstrap CLI

The OpenShift GitOps Cluster Bootstrap CLI is a command line tool that helps you bootstrap an OpenShift GitOps cluster. The CLI is designed to be used to bootrap the folder structure and files needed to create an OpenShift GitOps cluster using ArgoCD.

## How it works

When you run the CLI, it will ask you a series of questions to gather information about your OpenShift GitOps cluster. Let's say you want to create an OpenShift Cluster with the following information:

```yaml
Environment: dev
Stage: testing
Cluster name: my-cluster
```

The first thing we need to create is the environment. For this, just select the "Create Environment" option and provide the name of the environment. In our case, the environment name is `dev`.

The next step is to create the stage. For this, just select the "Create Stage" option, select the "Environment" you just created, and provide the name of the stage. In our case, the stage name is `testing`.

The last step is to create the cluster. For this, just select the "Create Cluster" option, select the "Environment" and "Stage" you just created, and provide the name of the cluster. In our case, the cluster name is `my-cluster`.

For each of the steps, you will be asked if you want to add additional properties. Properties are key-value pairs that you can use to store additional information about the environment, stage, or cluster. For example, you can add properties like `region`, `zone`, `network`, etc. These properties can then be used in the templates to customize the resources that will be created.

After you have created the environment, stage, and cluster, the CLI will create the `PROJECT.yaml` file, which contains the information about the environment, stage, and cluster. The CLI will also create the folder structure for our OpenShift cluster.

```plaintext
├── PROJECT.yaml
├── examples
│   └── templates
│       ├── addons
│       │   └── disco-operator
│       │       ├── kustomization.yaml
│       │       ├── manifest.yaml
│       │       └── patch.yaml
│       ├── appofapps
│       │   ├── kustomization.yaml
│       │   ├── manifest.yaml
│       │   └── values.yaml
│       └── cluster-config
│           ├── kyverno
│           │   └── kustomization.yaml
│           ├── manifest.yaml
│           └── other_test_with_patches
│               ├── kustomization.yaml
│               └── patch.yaml
└── overlays
    └── dev
        └── testing
            └── my-cluster
                ├── app-of-apps
                │   ├── kustomization.yaml
                │   └── values.yaml
                └── cluster-config
                    ├── kyverno
                    │   └── kustomization.yaml
                    └── other_test_with_patches
                        ├── kustomization.yaml
                        └── patch.yaml
```

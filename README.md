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
│       ├── appofapps
│       │   ├── kustomization.yaml
│       │   ├── manifest.yaml
│       │   └── values.yaml
│       └── cluster-config
│           ├── manifest.yaml
│           ├── kyverno
│           │   └── kustomization.yaml
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

## What is a cluster addon?

A cluster addon is a set of resources that can be applied to an OpenShift cluster to extend its functionality. For example, you can create a cluster addon that installs a set of operators, CRDs, and other resources that are needed to run a specific application on the cluster. The cluster addon can be applied to the cluster using ArgoCD.

### How to create a cluster addon?

To create a cluster addon, you need to create a folder with the name of the addon somewhere in the git repository. Inside the folder, you need to create a `manifest.yaml` file that contains the information about the addon.

```yaml
name: my-addon                                      # The name of the addon
properties:                                         # The properties define a set of key-value pairs the user has to enter during cluster creation
  my-property:                                      # The name of the property
    description: My property description            # The description of the property
    type: string                                    # The type of the property (string, int, bool)
    required: true                                  # If the property is required
    default: my-default-value                       # The default value of the property
files:                                              # The files define a set of files that will be created in the cluster folder during cluster creation
  - values.yaml                                     # A reference to the file inside the same folder as the manifest.yaml
  - resources/                                      # A reference to the folder inside the same folder as the manifest.yaml
```

After you have created the `manifest.yaml` file and the files that are needed for the addon, you can add the addon to the `PROJECT.yaml` file. To do this, execute the `ogc` binary and select the "Add Addon" option. The CLI will ask you for the name of the addon and the path to the folder where the addon is located. The CLI will then add the addon to the `PROJECT.yaml` file. The next time you create or update a cluster, the addon will be included in the selection of addons.

### Example

Let's say you want to create a cluster addon that installs the `grafana` stack on the cluster.

1) Create a folder with the name of the addon:

```bash
mkdir -p examples/addons/grafana
```

2) Create all resources that are needed for the addon:

```bash
cat <<EOF> examples/addons/grafana/values.yaml
# https://github.com/grafana/helm-charts/blob/main/charts/grafana/values.yaml
ingress:
  enabled: true
  hosts:
    - grafana.my-domain.local
resources:
  requests:
    cpu: {{ .Properties.request_cpu | default "100m" }}
    memory: {{ .Properties.request_memory | default "128Mi" }}
  limits:
    cpu: {{ .Properties.limit_cpu | default "100m" }}
    memory: {{ .Properties.limit_memory | default "128Mi" }}
EOF
```

3) Create the `manifest.yaml` file and include the values.yaml file:

```bash
cat <<EOF> examples/addons/grafana/manifest.yaml
name: grafana
properties:
  request_cpu:
    description: The CPU request for the grafana pod
    type: string
    required: false
    default: 100m
  request_memory:
    description: The memory request for the grafana pod
    type: string
    required: false
    default: 128Mi
  limit_cpu:
    description: The CPU limit for the grafana pod
    type: string
    required: true
    default: 100m
  limit_memory:
    description: The memory limit for the grafana pod
    type: string
    required: true
    default: 128Mi
files:
  - values.yaml
EOF
```

4) Add the addon to the `PROJECT.yaml` file:

```bash
user@pc % ogc
✔ Add Addon
Addon Name []: grafana
Should this addon be enabled by default? [Y/N]: y
Please provide the path to the location of the addon (the directory must contain a manifest.yaml file) []: examples/addons/grafana/
Are you sure you want to create the addon? [Y/N]: y
```

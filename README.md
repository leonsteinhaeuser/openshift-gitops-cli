# OpenShift GitOps Cluster Bootstrap CLI

The OpenShift GitOps Cluster Bootstrap CLI is a command line tool that helps you bootstrap an OpenShift GitOps cluster. The CLI is designed to be used to bootrap the folder structure and files needed to create an OpenShift GitOps cluster using ArgoCD.

## How it works

When you run the CLI, it will ask you a series of questions to gather information about your OpenShift GitOps cluster. Let's say you want to create an OpenShift Cluster with the following information:

```yaml
Environment: dev
Stage: testing
Cluster name: my-cluster
```

The first thing we need to create is the environment. For this, select the "Manage Environment" and next select "Create". The CLI will ask you several questions to gather information about the environment. For example, you can provide the name of the environment and properties. When you have provided all the information, select "Done". This option will lead you back to the main menu.

The next step is to create the stage. For this, select the "Manage Stage" option, select the "Environment" you just created, and answer the questions the CLI provides. For example, you can provide the name of the stage and properties. When you have provided all the information, select "Done". This option will lead you back to the main menu.

The last step is to create the cluster. For this, select the "Manage Cluster" option followed by the "Create" option. Select the "Environment" and "Stage" you just created and answer the questions the CLI provides. For example, you can provide the name of the cluster, addons, and properties. When you have provided all the information, select "Done". This option will lead you back to the main menu and create the folder structure and files needed to create the bootstrap.

For each of the steps, you will be asked if you want to add additional properties. Properties are key-value pairs that you can use to store additional information about the environment, stage, or cluster. For example, you can add properties like `region`, `zone`, `network`, etc. These properties can then be used in the templates to customize the resources that will be created. If you define the same key in the environment, stage, and cluster, the value of the highest level will be used.

* For example, if you define a property called `region` in the environment and the stage, the value of the stage will be used.
* If you define a property called `region` in the environment and the cluster, the value of the cluster will be used.
* If you define a property called `region` in the stage and the cluster, the value of the cluster will be used.

After you have created the environment, stage, and cluster, the CLI will create the `PROJECT.yaml` file, which contains the information about the environment, stage, and cluster. The CLI will also create the folder structure for our OpenShift cluster.

```plaintext
├── PROJECT.yaml
├── examples
│   ├── addons
│   │   └── disco-operator
│   │       ├── kustomization.yaml
│   │       ├── manifest.yaml
│   │       └── patch.yaml
│   └── templates
│       └── appofapps
│           └── database
│               ├── kustomization.yaml
│               ├── manifest.yaml
│               └── values.yaml
└── overlays
    └── dev                                           # The environment
        └── dev01                                     # The stage
            └── test01                                # The cluster
                ├── app-of-apps
                │   ├── kustomization.yaml
                │   └── values.yaml
                └── cluster-configs                   # The addon group folder
                    └── disco-operator                # The addon folder
                        ├── kustomization.yaml        # The kustomization file
                        └── patch.yaml                # The patch file (imported by the kustomization file)
```

## What is an environment and stage?

An environment in terms of infrastructure is a collection of resources that share the same hardware and network. For example, you can have infrastructure at `aws`, `gcp`, `azure` or even `on-prem`. Each of these environments can be named accordingly. For example, you can have an environment called `aws` which is hosted on `aws`, an environment called `gcp` which is hosted on `gcp`, and an environment called `azure` which is hosted on `azure`.

A stage is a subset of an environment. For example, you can have a `aws` environment with multiple stages like `dev`, `staging`, and `production`, and a `gcp` environment with multiple stages like `dev`, `staging`, and `production`. Each stage can have its own set of resources that are used for different purposes. For example, the `dev` stage can be used for development with unrestricted access to resources, the `staging` stage can be used for testing production-like applications, and the `production` stage can be used for running production applications with restricted access to resources.

## What is a cluster?

A cluster is a set of resources that are used to run an application. For example, you can have a cluster that runs a web application, a database, and a cache. A second cluster can run a web application, a database, and a kafka cluster. Each cluster can have a baseline set of resources that are used to provision the cluster. For example, a policy engine, network policies, ingress controllers, monitoring, tracing, and logging or any other resources that are needed to run the application. We call this type of resources `addons`. Addons are pre-configured resources that can be applied to a cluster to extend its functionality. During the cluster creation, you can select which addons you want to apply to the cluster.

### How to create a cluster?

To create a cluster, execute the `ogc` binary and select the "Create Cluster" option. The CLI will ask you for the name of the `environment`, `stage`, and `cluster` name, followed by `addons` that you want to apply to the cluster and `properties` that you want to set for the cluster.

### Example

Let's say you want to create an OpenShift cluster with the following information:

```bash
Environment: dev
Stage: dev
Cluster name: my-cluster
Addons: [kyverno, monitoring{ingress_host: monitoring.my-domain.local}]
Properties: [gitURL: https://git.example.com/repo/name.git, gitBranch: main]
```

1) Create the environment:

```bash
user@pc % ogc
✔ Manage Environment
✔ Create
Environment Name: aws
Create Property: gitURL
Property Value: https://git.my-url.local
✔ Done
✔ Done
```

2) Create the stage:

```bash
user@pc % ogc
✔ Manage Stage
✔ Create
✔ aws
Stage Name: dev
Create Property: gitBranch
Property Value: main
✔ Done
```

3) Create the cluster:

```bash
✔ Manage Cluster
✔ Create
✔ aws
✔ dev
Cluster Name: my-cluster
✔ Addons
✔ monitoring
✔ Enable
Enable addon monitoring
✔ Done
✔ kyverno
✔ Enable
Enable addon kyverno
✔ Done
✔ cluster-policies
✔ Enable
Enable addon cluster-policies
✔ Settings
✔ enableNetworkPolicies
Value: true
✔ Done
✔ Done
✔ Done
✔ Done
✔ Done
```

When you have created the environment, stage, and cluster, the CLI will create the corresponding entries in the `PROJECT.yaml` and directory structure.

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

The **group** is the name of the parent folder where the addon is located. In our case, the group is called ***cluster-configs***. During the *cluster* creation, this will result in the following folder structure for the addon configs: `<basePath>/<environment>/<stage>/<cluster>/cluster-configs/grafana`

```bash
user@pc % ogc
✔ Manage Addon
✔ Create
Addon Name: grafana
Should this addon be enabled by default? [Y/N]: y
Please provide the path to the location of the addon (the directory must contain a manifest.yaml file): examples/addons/grafana
Create new group: monitoring
Are you sure you want to create the addon? [Y/N]: y
✔ Done
```

## Template / Addon scopes

No matter if you define an addon or a template, you always have access to the following variables:

| Variable | Scope | Description |
| --- | --- | --- |
| `{{ .Environment }}` | addon, tempate | The environment variable returns the name of the environment we are currenlty in |
| `{{ .Stage }}` | addon, tempate | The stage variable returns the name of the stage we are currently in |
| `{{ .Cluster }}` | addon, tempate | The cluster variable returns the name of the cluster we are currently in |
| `{{ .Properties.<key> }}` | addon, tempate | The properties variable returns the value of the property with the key `<key>`. The property keys in addons differ from the property keys in the template, as the addon does not currently have access to the environment, stage or cluster properties. In order for the addon to have properties available, you must define a property key in the `manifest.yaml` file. All properties defined there are then available for your addon template files. |
| `{{ .ClusterProperties.<key> }}` | addon | The cluster properties is a map that contains all properties that are defined for the cluster. |

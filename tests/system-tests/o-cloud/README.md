# O-Cloud System Tests

## Overview

System-level test automation for O-RAN O2IMS-based SNO (Single Node OpenShift) cluster lifecycle management. Tests run against a pre-installed hub cluster and validate provisioning, deprovisioning, day 2 operations, and alarm management through the ORAN O2IMS interface using ClusterTemplates and ProvisioningRequests.

### Prerequisites

* OCP hub cluster with the following operators deployed:
  * ORAN O2IMS operator (`oran-o2ims` namespace)
  * ORAN Hardware Manager Plugin (`oran-hwmgr-plugin` namespace)
  * Advanced Cluster Management (ACM) (`rhacm` namespace)
  * OpenShift GitOps (`openshift-operators` namespace)
* At least two BareMetalHost resources available in the inventory pool
* BMC/Redfish access configured for spoke nodes
* `KUBECONFIG` set to a valid kubeconfig for the hub cluster

#### IBI-specific prerequisites

* `openshift-install` binary available in `PATH` (for seed image generation)
* Lifecycle Agent operator deployed on seed SNO
* SR-IOV operator deployed on seed SNO
* SSH access to spoke nodes (for IBI completion verification)
* Local registry configured for seed images

#### Alarms-specific prerequisites

* O2IMS API accessible via `ECO_OCLOUD_O2IMS_BASE_URL`
* PTP operator deployed on spoke SNO clusters
* Subscriber route available for alarm notifications

### Test suites

| File | Label | Description |
|------|-------|-------------|
| [sno-provisioning-ai.go](tests/sno-provisioning-ai.go) | `ocloud-ai-provisioning` | AI-based SNO provisioning: single cluster success/failure, simultaneous provisioning/deprovisioning with same and different ClusterTemplates |
| [sno-provisioning-ibi.go](tests/sno-provisioning-ibi.go) | `ocloud-ibi-provisioning` | IBI-based SNO provisioning: seed image generation, base image installation via BMC/Redfish, single cluster success/failure |
| [day2-configuration.go](tests/day2-configuration.go) | `ocloud-day2-configuration` | Day 2 operator upgrades: successful upgrade across all SNOs, failed upgrade in all SNOs, failed upgrade in a subset of SNOs |
| [alarms.go](tests/alarms.go) | `ocloud-alarms` | O2IMS alarm API: alarm retrieval after policy violation, alarm cleanup after retention period |

### Internal packages

[**ocloudcommon**](internal/ocloudcommon/)
- All test logic as exported functions called directly from `It()` blocks. Includes provisioning/deprovisioning workflows, cluster instance verification, policy compliance checks, BMH availability, alarm API interactions, and operator upgrade helpers.

[**ocloudconfig**](internal/ocloudconfig/config.go)
- Configuration loading from [`default.yaml`](internal/ocloudconfig/default.yaml) with environment variable overrides (`ECO_OCLOUD_*` prefix via `envconfig`). Includes BMC client construction for both spoke nodes.

[**ocloudinittools**](internal/ocloudinittools/ocloudinittools.go)
- Exports `HubAPIClient`, `OCloudConfig`, and `BMCClient` for use via dot-import in all test and helper files.

[**ocloudparams**](internal/ocloudparams/)
- Suite-level constants (labels, namespaces, timeouts, command templates), cluster instance parameter maps, and reporter configuration (namespaces and CRDs to dump on failure).

### Shared system-tests packages

The suite depends on shared helpers from `tests/system-tests/internal/`:

| Package | Purpose |
|---------|---------|
| `shell` | `ExecuteCmd()` — run shell commands (skopeo, oc, openshift-install) |
| `sshcommand` | `SSHCommand()` — execute commands on spoke nodes via SSH |
| `files` | `CopyFile()` — file copy operations |
| `csv` | `GetCurrentCSVNameFromSubscription()` — OLM CSV helpers |

### Eco-goinfra packages

- [**oran**](https://github.com/rh-ecosystem-edge/eco-goinfra/tree/main/pkg/oran) — ProvisioningRequest, AllocatedNode builders; O2IMS API client (alarms, subscriptions)
- [**siteconfig**](https://github.com/rh-ecosystem-edge/eco-goinfra/tree/main/pkg/siteconfig) — ClusterInstance builder
- [**bmh**](https://github.com/rh-ecosystem-edge/eco-goinfra/tree/main/pkg/bmh) — BareMetalHost builder
- [**bmc**](https://github.com/rh-ecosystem-edge/eco-goinfra/tree/main/pkg/bmc) — BMC client (Redfish boot, power control)
- [**ibi**](https://github.com/rh-ecosystem-edge/eco-goinfra/tree/main/pkg/ibi) — ImageClusterInstall builder
- [**lca**](https://github.com/rh-ecosystem-edge/eco-goinfra/tree/main/pkg/lca) — SeedGenerator builder
- [**ocm**](https://github.com/rh-ecosystem-edge/eco-goinfra/tree/main/pkg/ocm) — ManagedCluster and Policy helpers
- [**olm**](https://github.com/rh-ecosystem-edge/eco-goinfra/tree/main/pkg/olm) — ClusterServiceVersion for operator version tracking
- [**sriov**](https://github.com/rh-ecosystem-edge/eco-goinfra/tree/main/pkg/sriov) — NetworkNodeState for SR-IOV verification during seed generation
- [**reportxml**](https://github.com/rh-ecosystem-edge/eco-goinfra/tree/main/pkg/reportxml) — XML report generation

## Environment Variables

### Hub and cluster identity

| Variable | Default | Description |
|----------|---------|-------------|
| `KUBECONFIG` | _(required)_ | Path to hub cluster kubeconfig |
| `ECO_OCLOUD_CLUSTER_NAME_1` | _(empty)_ | Name of the first cluster |
| `ECO_OCLOUD_CLUSTER_NAME_2` | _(empty)_ | Name of the second cluster |
| `ECO_OCLOUD_HOSTNAME_1` | _(empty)_ | Hostname of the first cluster node |
| `ECO_OCLOUD_HOSTNAME_2` | _(empty)_ | Hostname of the second cluster node |
| `ECO_OCLOUD_NODE_CLUSTER_NAME_1` | _(empty)_ | Name of the first ORAN Node Cluster |
| `ECO_OCLOUD_NODE_CLUSTER_NAME_2` | _(empty)_ | Name of the second ORAN Node Cluster |
| `ECO_OCLOUD_OCLOUD_SITE_ID` | _(empty)_ | ID of the ORAN O-Cloud Site |
| `ECO_OCLOUD_INVENTORY_POOL_NAMESPACE` | _(empty)_ | Namespace of the inventory pool |
| `ECO_OCLOUD_BMH_SPOKE1` | _(empty)_ | BareMetalHost resource name for spoke 1 |
| `ECO_OCLOUD_BMH_SPOKE2` | _(empty)_ | BareMetalHost resource name for spoke 2 |

### ClusterTemplate versions

| Variable | Default | Description |
|----------|---------|-------------|
| `ECO_OCLOUD_TEMPLATE_NAME` | _(empty)_ | Base name of the referenced ClusterTemplate |
| `ECO_OCLOUD_TEMPLATE_VERSION_AI_SUCCESS` | _(empty)_ | ClusterTemplate version for successful AI-based SNO provisioning |
| `ECO_OCLOUD_TEMPLATE_VERSION_AI_FAILURE` | _(empty)_ | ClusterTemplate version for failing AI-based SNO provisioning |
| `ECO_OCLOUD_TEMPLATE_VERSION_SIMULTANEOUS_1` | _(empty)_ | First ClusterTemplate version for multi-cluster provisioning |
| `ECO_OCLOUD_TEMPLATE_VERSION_SIMULTANEOUS_2` | _(empty)_ | Second ClusterTemplate version for multi-cluster provisioning |
| `ECO_OCLOUD_TEMPLATE_VERSION_IBI_SUCCESS` | _(empty)_ | ClusterTemplate version for successful IBI-based SNO provisioning |
| `ECO_OCLOUD_TEMPLATE_VERSION_IBI_FAILURE` | _(empty)_ | ClusterTemplate version for failing IBI-based SNO provisioning |
| `ECO_OCLOUD_TEMPLATE_VERSION_DAY2` | _(empty)_ | ClusterTemplate version for Day 2 operations |
| `ECO_OCLOUD_TEMPLATE_VERSION_SEED` | _(empty)_ | ClusterTemplate version for IBI seed cluster provisioning |

### BMC configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `ECO_OCLOUD_SPOKE1_BMC_USERNAME` | _(empty)_ | BMC username for spoke 1 |
| `ECO_OCLOUD_SPOKE1_BMC_PASSWORD` | _(empty)_ | BMC password for spoke 1 |
| `ECO_OCLOUD_SPOKE1_BMC_HOST` | _(empty)_ | BMC IP address for spoke 1 |
| `ECO_OCLOUD_SPOKE1_BMC_TIMEOUT` | `15s` | BMC operation timeout for spoke 1 |
| `ECO_OCLOUD_SPOKE2_BMC_USERNAME` | _(empty)_ | BMC username for spoke 2 |
| `ECO_OCLOUD_SPOKE2_BMC_PASSWORD` | _(empty)_ | BMC password for spoke 2 |
| `ECO_OCLOUD_SPOKE2_BMC_HOST` | _(empty)_ | BMC IP address for spoke 2 |
| `ECO_OCLOUD_SPOKE2_BMC_TIMEOUT` | `15s` | BMC operation timeout for spoke 2 |

### Network configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `ECO_OCLOUD_INTERFACE_NAME` | _(empty)_ | Network interface name on the spoke nodes |
| `ECO_OCLOUD_INTERFACE_IPV6_1` | _(empty)_ | IPv6 address of the interface for the first cluster |
| `ECO_OCLOUD_INTERFACE_IPV6_2` | _(empty)_ | IPv6 address of the interface for the second cluster |
| `ECO_OCLOUD_DNS_IPV6` | _(empty)_ | IPv6 address of the DNS server |
| `ECO_OCLOUD_NEXT_HOP_IPV6` | _(empty)_ | IPv6 address of the next hop gateway |
| `ECO_OCLOUD_NEXT_HOP_INTERFACE` | _(empty)_ | Network interface for the next hop route |
| `ECO_OCLOUD_SSH_CLUSTER_2` | _(empty)_ | SSH address for the second cluster |
| `ECO_OCLOUD_SSH_KEY` | _(empty)_ | SSH public key for cluster node access |

### IBI (Image Based Install) configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `ECO_OCLOUD_IBI_GENERATE_SEED_IMAGE` | `true` | Generate the seed image for Image Based Install |
| `ECO_OCLOUD_IBI_BASE_IMAGE_PATH` | _(empty)_ | Local path to the IBI base image |
| `ECO_OCLOUD_IBI_BASE_IMAGE_URL` | _(empty)_ | URL to the IBI base image |
| `ECO_OCLOUD_VIRTUAL_MEDIA_ID` | _(empty)_ | Virtual media ID for BMC boot |
| `ECO_OCLOUD_SEED_IMAGE` | _(empty)_ | Seed container image reference |
| `ECO_OCLOUD_SEED_VERSION` | _(empty)_ | Version of the seed image |
| `ECO_OCLOUD_BASE_IMAGE_NAME` | _(empty)_ | Name of the base image |

### Registry and authentication

| Variable | Default | Description |
|----------|---------|-------------|
| `ECO_OCLOUD_REGISTRY_5000` | _(empty)_ | URL for registry on port 5000 |
| `ECO_OCLOUD_REGISTRY_5005` | _(empty)_ | URL for registry on port 5005 |
| `ECO_OCLOUD_LOCAL_REGISTRY_AUTH` | _(empty)_ | Authentication credentials for the local registry |
| `ECO_OCLOUD_PULL_SECRET` | _(empty)_ | Pull secret for container image registries |
| `ECO_OCLOUD_AUTHFILE_PATH` | `/kubeconfig/docker/config.json` | Path to the auth file for Skopeo commands |

### Alarms and O2IMS API

| Variable | Default | Description |
|----------|---------|-------------|
| `ECO_OCLOUD_SUBSCRIBER_URL` | _(empty)_ | URL of the O-Cloud event subscriber |
| `ECO_OCLOUD_SUBSCRIBER_DOMAIN` | _(empty)_ | Domain of the O-Cloud event subscriber |
| `ECO_OCLOUD_O2IMS_BASE_URL` | _(empty)_ | Base URL for the O2IMS API |

## Running O-Cloud Test Suites

### Running all tests

```bash
export KUBECONFIG=/path/to/kubeconfig
export ECO_TEST_FEATURES="system-tests/o-cloud"
# Set all required ECO_OCLOUD_* variables
make run-tests
```

### Running AI provisioning tests

```bash
export KUBECONFIG=/path/to/kubeconfig
export ECO_OCLOUD_TEMPLATE_NAME="my-template"
export ECO_OCLOUD_TEMPLATE_VERSION_AI_SUCCESS="v1.0"
export ECO_OCLOUD_TEMPLATE_VERSION_AI_FAILURE="v1.0-fail"
export ECO_OCLOUD_TEMPLATE_VERSION_SIMULTANEOUS_1="v1.0"
export ECO_OCLOUD_TEMPLATE_VERSION_SIMULTANEOUS_2="v2.0"
# Set cluster, BMH, and network variables
export ECO_TEST_FEATURES="system-tests/o-cloud"
export ECO_TEST_LABELS="ocloud-ai-provisioning"
make run-tests
```

### Running IBI provisioning tests

```bash
export KUBECONFIG=/path/to/kubeconfig
export ECO_OCLOUD_IBI_GENERATE_SEED_IMAGE=true
export ECO_OCLOUD_IBI_BASE_IMAGE_URL="https://registry.example.com/rhcos-ibi.iso"
export ECO_OCLOUD_VIRTUAL_MEDIA_ID="1"
export ECO_OCLOUD_SEED_IMAGE="registry.example.com:5000/seedimage:latest"
export ECO_OCLOUD_SEED_VERSION="4.20.0"
# Set spoke2 BMC, SSH, and network variables
export ECO_TEST_FEATURES="system-tests/o-cloud"
export ECO_TEST_LABELS="ocloud-ibi-provisioning"
make run-tests
```

### Running Day 2 configuration tests

```bash
export KUBECONFIG=/path/to/kubeconfig
export ECO_OCLOUD_TEMPLATE_VERSION_DAY2="v1.0-day2"
export ECO_OCLOUD_AUTHFILE_PATH="/path/to/auth.json"
export ECO_OCLOUD_REGISTRY_5000="registry.example.com:5000"
# Set both spoke BMC, cluster, and network variables
export ECO_TEST_FEATURES="system-tests/o-cloud"
export ECO_TEST_LABELS="ocloud-day2-configuration"
make run-tests
```

### Running alarm tests

```bash
export KUBECONFIG=/path/to/kubeconfig
export ECO_OCLOUD_O2IMS_BASE_URL="https://o2ims.example.com"
export ECO_OCLOUD_SUBSCRIBER_URL="https://subscriber.apps.example.com"
export ECO_OCLOUD_SUBSCRIBER_DOMAIN="apps.example.com"
# Set cluster and template variables for AI provisioning (alarms test provisions a cluster first)
export ECO_TEST_FEATURES="system-tests/o-cloud"
export ECO_TEST_LABELS="ocloud-alarms"
make run-tests
```

### Running directly with ginkgo

```bash
cd tests/system-tests/o-cloud
ginkgo -v -timeout=24h --label-filter="ocloud-ai-provisioning" ./...
```

## Additional Information

Please refer to the [project README](../../../README.md) for global configuration, reporting options, and general test framework documentation.

# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Suite Tests

O-Cloud system tests validate ORAN O2IMS-based SNO (Single Node OpenShift) cluster lifecycle management on a pre-installed hub cluster. Test areas:

- **AI provisioning** (`ocloud-ai-provisioning`): Assisted Installer-based SNO provisioning/deprovisioning, including simultaneous multi-cluster scenarios (same and different ClusterTemplates)
- **IBI provisioning** (`ocloud-ibi-provisioning`): Image Based Installer-based SNO provisioning, including seed image generation and base image installation via BMC/Redfish
- **Day 2 configuration** (`ocloud-day2-configuration`): Operator upgrades across spoke SNO clusters (success, full failure, partial failure)
- **Alarms** (`ocloud-alarms`): O2IMS alarm API retrieval, subscription notifications, and retention-period cleanup

## Running Tests

Tests require a live hub cluster with spoke BMH resources. Set `KUBECONFIG` and the required `ECO_OCLOUD_*` environment variables (see `README.md` for the full list).

```bash
# Run the full o-cloud suite
cd tests/system-tests/o-cloud
ginkgo -v -timeout=24h ./...

# Run a specific test area by label
ginkgo -v -timeout=24h --label-filter="ocloud-ai-provisioning" ./...
ginkgo -v -timeout=24h --label-filter="ocloud-alarms" ./...

# From repo root using the test runner
export ECO_TEST_FEATURES="system-tests/o-cloud"
export ECO_TEST_LABELS="ocloud-ai-provisioning"
make run-tests
```

## Suite Architecture

### Initialization chain

`o_cloud_suite_test.go` imports `tests/internal/inittools` (provides `APIClient` to the hub cluster) and blank-imports the `tests/` package to register all Ginkgo specs.

Test files and helpers access the hub client and config through `ocloudinittools`:
```go
import . "github.com/rh-ecosystem-edge/eco-gotests/tests/system-tests/o-cloud/internal/ocloudinittools"
// Provides: HubAPIClient, OCloudConfig, BMCClient
```

### Package layout

- `tests/` — Ginkgo spec files (thin wrappers that call into `ocloudcommon` helpers). Package name is `o_cloud_system_test`.
- `internal/ocloudcommon/` — All test logic lives here as exported functions (e.g. `VerifySuccessfulSnoProvisioning`, `VerifySuccessfulAlarmRetrieval`). Test files reference these directly in `It()` blocks.
- `internal/ocloudconfig/` — `OCloudConfig` struct, populated from `default.yaml` then overridden by `ECO_OCLOUD_*` env vars via `envconfig`.
- `internal/ocloudparams/` — Constants (namespaces, labels, timeouts, command templates) and package-level variables (cluster instance parameter maps, reporter config). The Ginkgo label is `"ocloud"`.
- `internal/ocloudinittools/` — Singleton init: creates `OCloudConfig`, sets `HubAPIClient`, exposes `BMCClient`. Dot-import pattern.

### Shared system-tests infrastructure

- `tests/system-tests/internal/` — Shared packages used across all system-test suites: `shell` (command execution), `sshcommand`, `files`, `csv` (OLM CSV helpers), `systemtestsconfig` (base config embedded by `OCloudConfig`).

### Test pattern: helpers with Gomega assertions

Unlike most eco-gotests suites, `ocloudcommon` helpers contain Gomega assertions directly (`Expect`, `Eventually`). This is the established pattern here — test files delegate entire test scenarios to single exported functions.

Concurrent verification uses `sync.WaitGroup` with `GinkgoRecover()` in goroutines (see `DeprovisionAiSnoCluster`).

### Key CRDs and resources

Tests interact with these via `eco-goinfra` builders:
- `ProvisioningRequest` (`oran.NewPRBuilder`, `oran.PullPR`) — created with UUID names
- `ClusterInstance` (`siteconfig.PullClusterInstance`)
- `AllocatedNode` (`oran.PullAllocatedNode`)
- `ImageClusterInstall` (`ibi.PullImageClusterInstall`) — IBI only
- `BareMetalHost` (`bmh.Pull`) — verified available before/after provisioning
- `SeedGenerator` (`lca.NewSeedGeneratorBuilder`) — IBI seed image generation
- O2IMS Alarms API (`oranapi.NewClientBuilder`) — token-authenticated REST client

### Configuration

All config flows through `OCloudConfig` (from `ocloudinittools`). Values come from `internal/ocloudconfig/default.yaml` (defaults) overridden by `ECO_OCLOUD_*` environment variables. BMC clients for spoke1/spoke2 are constructed from the BMC host/user/password config fields.

## Conventions

- Every `Describe` block uses `Ordered, ContinueOnFailure` and carries `Label(ocloudparams.Label)`.
- Tests do NOT use `reportxml.ID()` — they use descriptive `It()` strings instead.
- Logging uses `klog.V(ocloudparams.OCloudLogLevel)` (level 90).
- Cluster instance parameters are defined as package-level vars in `ocloudparams/ocloudvars.go`, not inline.
- Deprovisioning verifications (PR deleted, namespace gone, ClusterInstance gone, AllocatedNodes gone) run concurrently via goroutines + WaitGroup.
- Shell commands (skopeo, oc, openshift-install) are executed via `shell.ExecuteCmd` from the shared system-tests package.

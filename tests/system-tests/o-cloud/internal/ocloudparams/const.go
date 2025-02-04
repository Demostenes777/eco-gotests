package ocloudparams

import (
	"time"
)

const (
	// Label represents O-Cloud system tests label that can be used for test cases selection.
	Label = "ocloud"
	
	// DefaultTimeout is the timeout used for test resources creation.
	DefaultTimeout = 900 * time.Second
	
	// OCloudLogLevel configures logging level for O-Cloud related tests.
	OCloudLogLevel = 90

	// AcmNamespace is the namespace for ACM.
	AcmNamespace = "rhacm"
	
	// AcmSubscriptionName is the name of the ACM operator subscription
	AcmSubscriptionName = "acm-operator-subscription"

	// AcmInstanceName is the name of the ACM multicluster hub instance
	AcmInstanceName = "multiclusterhub"

	// OpenshiftGitOpsNamespace is the namespace for the GitOps operator.
	OpenshiftGitOpsNamespace = "openshift-operators"

	// OpenshiftGitOpsSubscriptionName is the name of the GitOps operator subscription.
	OpenshiftGitOpsSubscriptionName = "openshift-gitops-operator-subscription"

	// OCloudO2ImsNamespace is the namespace for the O-Cloud manager operator.
	OCloudO2ImsNamespace = "oran-o2ims"				

	// OCloudO2ImsSubscriptionName is the name of the O-Cloud manager operator subscription.
	OCloudO2ImsSubscriptionName = "oran-o2ims-operator-subscription"

	// OCloudHardwareManagerPluginNamespace is the namespace for the O-Cloud hardware manager plugin operator.
	OCloudHardwareManagerPluginNamespace = "oran-hwmgr-plugin"

	// OCloudHardwareManagerSubscriptionName is the name of the O-Cloud hardware manager plugin operator subscription.
	OCloudHardwareManagerPluginSubscriptionName = "oran-hwmgr-plugin-operator-subscription"

	// PtpNamespace is the namespace for the PTP operator.
	PtpNamespace = "openshift-ptp"

	// PtpOperatorSubscriptionName is the name of the PTP operator subscription.
	PtpOperatorSubscriptionName = "ptp-operator-subscription"

	//nolint:lll
	PodmanTagOperatorUpgrade = "podman tag registry.hub01.oran.telcoqe.eng.rdu2.dc.redhat.com:5000/olm/redhat-operators:v4.18-prerelease-ptp-operator registry.hub01.oran.telcoqe.eng.rdu2.dc.redhat.com:5000/olm/redhat-operators:v4.18-day2"
	PodmanTagSriovUpgrade = "podman tag registry.hub01.oran.telcoqe.eng.rdu2.dc.redhat.com:5000/olm/far-edge-sriov-fec:v4.18-new registry.hub01.oran.telcoqe.eng.rdu2.dc.redhat.com:5000/olm/far-edge-sriov-fec:v4.18-day2"
	PodmanPushOperatorUpgrade = "podman push registry.hub01.oran.telcoqe.eng.rdu2.dc.redhat.com:5000/olm/redhat-operators:v4.18-day2"
	PodmanPushSriovUpgrade = "podman push registry.hub01.oran.telcoqe.eng.rdu2.dc.redhat.com:5000/olm/far-edge-sriov-fec:v4.18-day2"
	PodmanTagOperatorDowngrade = "podman tag registry.hub01.oran.telcoqe.eng.rdu2.dc.redhat.com:5000/olm/redhat-operators:v4.18 registry.hub01.oran.telcoqe.eng.rdu2.dc.redhat.com:5000/olm/redhat-operators:v4.18-day2"
	PodmanTagSriovDowngrade = "podman tag registry.hub01.oran.telcoqe.eng.rdu2.dc.redhat.com:5000/olm/far-edge-sriov-fec:v4.18 registry.hub01.oran.telcoqe.eng.rdu2.dc.redhat.com:5000/olm/far-edge-sriov-fec:v4.18-day2"
	PodmanPushOperatorDowngrade = "podman push registry.hub01.oran.telcoqe.eng.rdu2.dc.redhat.com:5000/olm/redhat-operators:v4.18-day2"
	PodmanPushSriovDowngrade = "podman push registry.hub01.oran.telcoqe.eng.rdu2.dc.redhat.com:5000/olm/far-edge-sriov-fec:v4.18-day2"

	SnoKubeconfigCreate = "oc -n %s get secret %s-admin-kubeconfig -o json | jq -r .data.kubeconfig | base64 -d > tmp/%s/auth/kubeconfig"

	
)

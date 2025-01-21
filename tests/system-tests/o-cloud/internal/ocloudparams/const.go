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

)

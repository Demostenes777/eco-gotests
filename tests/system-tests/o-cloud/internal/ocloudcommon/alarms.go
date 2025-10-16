package ocloudcommon

import (
	"crypto/tls"
	"slices"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	oranapi "github.com/rh-ecosystem-edge/eco-goinfra/pkg/oran/api"
	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/oran/api/filter"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/system-tests/internal/shell"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/system-tests/o-cloud/internal/ocloudparams"
	. "github.com/rh-ecosystem-edge/eco-gotests/tests/system-tests/o-cloud/internal/ocloudinittools"
)

// CreateO2IMSClient creates an O2IMS API client using token authentication and returns it.
func CreateO2IMSClient() (*oranapi.AlarmsClient) {
	By("creating the O2IMS API client")
	token, err := shell.ExecuteCmd("oc create token -n oran-o2ims test-client --duration=24h")
	Expect(err).ToNot(HaveOccurred(), "Failed to create token for O2IMS API")

	o2imsBaseURL := "https://o2ims.apps.hub03.oran.telcoqe.eng.rdu2.dc.redhat.com"
	clientBuilder := oranapi.NewClientBuilder(o2imsBaseURL).
		WithToken(string(token)).
		WithTLSConfig(&tls.Config{MinVersion: tls.VersionTLS12, InsecureSkipVerify: true})

	alarmsClient, err := clientBuilder.BuildAlarms()
	Expect(err).ToNot(HaveOccurred(), "Failed to create the O2IMS API client")

	return alarmsClient
}

// VerifySuccessfulAlarmRetrieval verifies the test case of the successful retrieval of an alarm from the API.
//
//nolint:funlen
func VerifySuccessfulAlarmRetrieval(ctx SpecContext) {
	subscriberURL := "https://oran-subscriber.apps.hub03.oran.telcoqe.eng.rdu2.dc.redhat.com"
	alarmsClient := CreateO2IMSClient()

	VerifyBmhIsAvailable(OCloudConfig.BmhSpoke1, OCloudConfig.InventoryPoolNamespace)
	VerifyBmhIsAvailable(OCloudConfig.BmhSpoke2, OCloudConfig.InventoryPoolNamespace)

	provisioningRequest := VerifyProvisionSnoCluster(
		OCloudConfig.TemplateName,
		OCloudConfig.TemplateVersionAISuccess,
		OCloudConfig.NodeClusterName1,
		OCloudConfig.OCloudSiteID,
		ocloudparams.PolicyTemplateParameters,
		ocloudparams.ClusterInstanceParameters1)

	VerifyOcloudCRsExist(provisioningRequest)

	clusterInstance := VerifyClusterInstanceCompleted(provisioningRequest, ctx)
	nsname := provisioningRequest.Object.Status.Extensions.ClusterDetails.Name


	By("creating a new subscription")
	subscriptionID := uuid.New()
	
	subscription, err := alarmsClient.CreateSubscription(oranapi.AlarmSubscriptionInfo{
		ConsumerSubscriptionId: &subscriptionID,
		Callback: subscriberURL + "/" + subscriptionID.String(),
	})
	Expect(err).ToNot(HaveOccurred(), "Failed to create a new subscription")

	By("filtering subscriptions by the first ConsumerSubscriptionId")
	consumerIDFilter := filter.Equals("consumerSubscriptionId", subscriptionID.String())

	filteredSubscriptions, err := alarmsClient.ListSubscriptions(consumerIDFilter)
	Expect(err).ToNot(HaveOccurred(), "Failed to filter subscriptions")

	By("verifying the filtered results contain the first subscription but not the second")
	containsSubscription1 := slices.ContainsFunc(filteredSubscriptions,
		func(subscription oranapi.AlarmSubscriptionInfo) bool {
			return subscription.ConsumerSubscriptionId.String() == subscriptionID.String()
		})
	
	Expect(containsSubscription1).To(BeTrue(), "First subscription should be found in filtered results")

	By("deleting the test subscriptions")
	err = alarmsClient.DeleteSubscription(*subscription.AlarmSubscriptionId)
	Expect(err).ToNot(HaveOccurred(), "Failed to delete first test subscription")
}

// VerifySuccessfulAlarmsCleanup verifies the test case where the alarms from the database are cleaned up after the retention period.
//
//nolint:funlen
func VerifySuccessfulAlarmsCleanup(ctx SpecContext) {
	Fail("Not implemented")
}

package rdsmanagementcommon

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/ccin2p3/go-freeipa/freeipa"
	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/openshift-kni/eco-gotests/tests/system-tests/rdsmanagement/internal/rdsmanagementinittools"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/rdsmanagement/internal/rdsmanagementparams"
	"golang.org/x/crypto/ssh"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
    kubevirtv1 "kubevirt.io/client-go/api/v1"
    kubevirtclient "kubevirt.io/client-go/kubecli"
)

// VerifyIDMInstallation the IDM installation.
func VerifyIDMInstallation() {
	vmUsername := RDSManagementConfig.IDMConfig.VMUsername
	vmPassword := RDSManagementConfig.IDMConfig.VMPassword
	serverIP := RDSManagementConfig.IDMConfig.IPAddress
	serverURL := RDSManagementConfig.IDMConfig.URL
	ipaAdminUser := RDSManagementConfig.IDMConfig.IPAAdminUser
	ipaAdminPass := RDSManagementConfig.IDMConfig.IPAAdminPass
	givenName := RDSManagementConfig.IDMConfig.TestUserGivenname
	surname := RDSManagementConfig.IDMConfig.TestUserSn
	groupName := RDSManagementConfig.IDMConfig.TestGroup

	By("Verify that SSH login to IDM VM works")

	config := &ssh.ClientConfig{
		User: vmUsername,
		Auth: []ssh.AuthMethod{
			ssh.Password(vmPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	glog.V(rdsmanagementparams.RdsManagementLogLevel).Infof(
		fmt.Sprintf("[%s] Starting SSH connection", serverIP))

	sshConnection, err := ssh.Dial("tcp", serverIP, config)
	Expect(err).NotTo(HaveOccurred(),
		fmt.Sprintf("[%s] Failed to connect the SSH server, err: %v", serverIP, err))
	defer sshConnection.Close()

	glog.V(rdsmanagementparams.RdsManagementLogLevel).Infof(
		fmt.Sprintf("[%s] Establishing SSH session", serverIP))

	sshSession, err := sshConnection.NewSession()
	Expect(err).NotTo(HaveOccurred(),
		fmt.Sprintf("[%s] Failed to create the SSH session, err: %v", serverIP, err))
	defer sshSession.Close()

	glog.V(rdsmanagementparams.RdsManagementLogLevel).Infof(
		fmt.Sprintf("[%s] Running command", serverIP))

	sshOutput, err := sshSession.CombinedOutput("echo hello")
	Expect(err).NotTo(HaveOccurred(),
		fmt.Sprintf("Failed to run command: %s, err: %v", serverIP, err))
	Expect(string(sshOutput)).To(Equal("hello\n"), "Unexpected command output")

	By("Verify that login to IDM web interface is successful")

	glog.V(rdsmanagementparams.RdsManagementLogLevel).Infof(
		fmt.Sprintf("[%s] Starting web-based connection", serverIP))

	tspt := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false, // Check if it works, change to true otherwise
		},
	}

	idmConnection, err := freeipa.Connect(serverURL, tspt, ipaAdminUser, ipaAdminPass)
	Expect(err).NotTo(HaveOccurred(),
		fmt.Sprintf("[%s] Failed to login to server, err: %v", serverIP, err))

	By("Verify that new user accounts can be created")

	glog.V(rdsmanagementparams.RdsManagementLogLevel).Infof(
		"[%s] Adding new user", serverIP)

	r := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	uid := fmt.Sprintf("%s%v", givenName, r.Int())

	userRes, err := idmConnection.UserAdd(&freeipa.UserAddArgs{
		Givenname: givenName,
		Sn:        surname,
	}, &freeipa.UserAddOptionalArgs{
		UID: freeipa.String(uid),
	})

	Expect(err).NotTo(HaveOccurred(),
		fmt.Sprintf("[%s] Failed to add user, err: %v", serverIP, err))

	glog.V(rdsmanagementparams.RdsManagementLogLevel).Infof(
		"Added user %v", userRes.Result.Cn)

	By("Verify that new groups can be created")

	glog.V(rdsmanagementparams.RdsManagementLogLevel).Infof(
		"[%s] Adding new group", serverIP)

	r = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	gid := r.Int()

	groupRes, err := idmConnection.GroupAdd(&freeipa.GroupAddArgs{
		Cn: groupName,
	}, &freeipa.GroupAddOptionalArgs{
		Gidnumber: &gid,
	})

	Expect(err).NotTo(HaveOccurred(),
		fmt.Sprintf("[%s] Failed to add group, err: %v", serverIP, err))

	glog.V(rdsmanagementparams.RdsManagementLogLevel).Infof(
		"Added group %v", groupRes.Result.Cn)
}

// VerifyIDMReplication verifies the IDM replication
func VerifyIDMReplication() {
	
}

// VerifyOCPIntegrationWithIDM verifies the OCP and IDM integration
func VerifyOCPIntegrationWithIDM() {

}

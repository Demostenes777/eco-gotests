package rdsmanagementcommon

import (
	"time"

	. "github.com/onsi/gomega"
	. "github.com/openshift-kni/eco-gotests/tests/system-tests/rdsmanagement/internal/rdsmanagementinittools"
	"golang.org/x/crypto/ssh"
)

// VerifySSHAccess verifies that SSH login to IDM VM is working.
func VerifySSHAccess() {
	server := RDSManagementConfig.IDMConfig.IPAddress
	user := RDSManagementConfig.IDMConfig.VMUsername
	password := RDSManagementConfig.IDMConfig.VMPassword

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	connection, err := ssh.Dial("tcp", server, config)
	Expect(err).NotTo(HaveOccurred(), "Failed to connect to the SSH server")
	defer connection.Close()

	session, err := connection.NewSession()
	Expect(err).NotTo(HaveOccurred(), "Failed to create the SSH session")
	defer session.Close()

	output, err := session.CombinedOutput("echo test")
	Expect(err).NotTo(HaveOccurred(), "Failed to run the command on SSH server")
	Expect(string(output)).To(Equal("test\n"), "Unexpected command output")
}

// VerifyWebAccess verifies that login to IDM web interface is successful.
func VerifyWebAccess() {
	// todo - Add test logic
}

// VerifyNewUserAccountCreation verifies that new user accounts can be created.
func VerifyNewUserAccountCreation() {
	// todo - Add test logic
}

// VerifyNewGroupCreation verifies that new user groups can be created.
func VerifyNewGroupCreation() {
	// todo - Add test logic
}

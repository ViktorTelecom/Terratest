package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/ssh"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestSHConnect(t *testing.T) {
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: "/home/viktor/go/src/terra_test",
	})

	// Clean up resources with "terraform destroy" at the end of the test.
	defer terraform.Destroy(t, terraformOptions)

	// Run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the values of output variables and check they have the expected values.
	// output := terraform.Output(t, terraformOptions, "hello_world")
	// assert.Equal(t, "Hello, World!", output)

	var MyRSAKeyPair = ssh.GenerateRSAKeyPair(t, 2048)

	var Host_1 = ssh.Host{
		Hostname:    terraform.Output(t, terraformOptions, "server_1_public_ip"),
		SshUserName: "ubuntu",
		SshKeyPair:  MyRSAKeyPair,
	}

	var Host_2 = ssh.Host{
		Hostname:    terraform.Output(t, terraformOptions, "server_2_public_ip"),
		SshUserName: "ubuntu",
		SshKeyPair:  MyRSAKeyPair,
	}

	assert.Equal(t, ssh.CheckSshConnectionE(t, Host_1), nil)
	assert.Equal(t, ssh.CheckSshConnectionE(t, Host_2), nil)
}

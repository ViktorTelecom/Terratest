package test

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/retry"
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

	Buffer, _ := ioutil.ReadFile("/home/viktor/.ssh/id_rsa.pub")
	MyPubKey := string(Buffer)
	Buffer, _ = ioutil.ReadFile("/home/viktor/.ssh/id_rsa")
	MyPrivKey := string(Buffer)

	//var MyRSAKeyPair = ssh.GenerateRSAKeyPair(t, 2048)
	MyRSAKeyPair := ssh.KeyPair{MyPubKey, MyPrivKey}

	var Host_1 = ssh.Host{
		Hostname:    terraform.Output(t, terraformOptions, "server_1_public_ip"),
		SshUserName: "ubuntu",
		SshKeyPair:  &MyRSAKeyPair,
	}

	var Host_2 = ssh.Host{
		Hostname:    terraform.Output(t, terraformOptions, "server_2_public_ip"),
		SshUserName: "ubuntu",
		SshKeyPair:  &MyRSAKeyPair,
	}

	// It can take a minute or so for the Instance to boot up, so retry a few times
	maxRetries := 30
	timeBetweenRetries := 5 * time.Second

	hosts := []ssh.Host{Host_1, Host_2}

	for _, host := range hosts {
		description := fmt.Sprintf("SSH to public host %s", host.Hostname)

		expectedText := "Hello, World"
		command := fmt.Sprintf("echo -n '%s'", expectedText)

		retry.DoWithRetry(t, description, maxRetries, timeBetweenRetries, func() (string, error) {
			actualText, err := ssh.CheckSshCommandE(t, host, command)

			if err != nil {
				return "", err
			}

			if strings.TrimSpace(actualText) != expectedText {
				return "", fmt.Errorf("Expected SSH command to return '%s' but got '%s'", expectedText, actualText)
			}

			assert.Equal(t, strings.TrimSpace(actualText), expectedText)

			return "", nil
		})
	}

	//assert.Equal(t, ssh.CheckSshConnectionE(t, Host_1), nil)
	//assert.Equal(t, ssh.CheckSshConnectionE(t, Host_2), nil)
}

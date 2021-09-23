package test

import (
	"flag"
	"fmt"
	"io/ioutil"
	"testing"

	//"os"
	//"path"
	"golang.org/x/crypto/ssh"

	//"github.com/gruntwork-io/terratest/modules/ssh"
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

	var (
		user   = flag.String("u", "", "ubuntu")
		pk     = flag.String("pk", "/home/viktor/.ssh/id_rsa", "Private key file")
		host_1 = flag.String("h", terraform.Output(t, terraformOptions, "server_1_public_ip"), "Host")
		//host_2 = flag.String("h", terraform.Output(t, terraformOptions, "server_2_public_ip"), "Host")
		port = flag.Int("p", 22, "Port")
	)

	flag.Parse()

	key, err := ioutil.ReadFile(*pk)
	if err != nil {
		panic(err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		panic(err)
	}

	config := &ssh.ClientConfig{
		User: *user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	/*
		config = &ssh.ClientConfig{
			User: "ubuntu",
			Auth: []ssh.AuthMethod{
				ssh.Password(""),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
	*/
	addr := fmt.Sprintf("%s:%d", *host_1, *port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		panic(err)
	}

	session, err := client.NewSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	b, err := session.CombinedOutput("uname -a")
	if err != nil {
		panic(err)
	}
	fmt.Print(string(b))

	// Run `terraform output` to get the values of output variables and check they have the expected values.
	assert.Equal(t, err, nil)
}

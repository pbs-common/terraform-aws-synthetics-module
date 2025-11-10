package test

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func testSynthetics(t *testing.T, variant string) {
	t.Parallel()

	terraformDir := fmt.Sprintf("../examples/%s", variant)

	terraformOptions := &terraform.Options{
		TerraformDir: terraformDir,
		LockTimeout:  "5m",
	}

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	expectedName := fmt.Sprintf("%s-%s", "ex-tf-synth", variant)
	name := terraform.Output(t, terraformOptions, "name")
	assert.Equal(t, expectedName, name)

	id := terraform.Output(t, terraformOptions, "id")
	assert.Equal(t, expectedName, id)

	awsRegion := getAWSRegion(t)
	awsAccountID := getAWSAccountID(t)

	expectedARN := fmt.Sprintf("arn:aws:synthetics:%s:%s:canary:%s", awsRegion, awsAccountID, name)
	arn := terraform.Output(t, terraformOptions, "arn")
	assert.Equal(t, expectedARN, arn)

	expectedPartialEngineARN := fmt.Sprintf("arn:aws:lambda:%s:%s:function:cwsyn-%s-", awsRegion, awsAccountID, name)
	engineARN := terraform.Output(t, terraformOptions, "engine_arn")
	assert.Contains(t, engineARN, expectedPartialEngineARN)

	// Only check source_location_arn if present
	sourceLocationARN := terraform.Output(t, terraformOptions, "source_location_arn")
	if sourceLocationARN != "" {
		expectedPartialSourceLocationARN := fmt.Sprintf("arn:aws:lambda:%s:%s:layer:cwsyn-%s-", awsRegion, awsAccountID, name)
		assert.Contains(t, sourceLocationARN, expectedPartialSourceLocationARN)
	} else {
		t.Log("Skipping source_location_arn check since it's not available for zip_file-based canaries")
	}

	expectedStatus := "RUNNING"
	status := terraform.Output(t, terraformOptions, "status")
	assert.Equal(t, expectedStatus, status)

	timeline := terraform.Output(t, terraformOptions, "timeline")
	assert.NotEmpty(t, timeline)
}

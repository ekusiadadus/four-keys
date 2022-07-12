package cli

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/urfave/cli/v2"
)

func TestDefaultAppShouldReturnMetricsWithRepositoryUrlSinceUntil(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	defaltApp := DefaultApp()
	testApp := &cli.App{
		Flags:  defaltApp.Flags,
		Action: defaltApp.Action,
		Writer: output,
	}

	err := testApp.Run([]string{"four-keys", "--repository", "https://github.com/go-git/go-git", "--since", "2020-01-01", "--until", "2020-12-31"})
	if err != nil {
		t.Errorf(err.Error())
	}
	var cliOutput DefaultCliOutput
	json.Unmarshal(output.Bytes(), &cliOutput)
	// intended output
	// {
	//   "option":{"since":"2020-01-01T00:00:00Z","until":"2020-12-31T23:59:59Z"},
	//   "deploymentFrequency":0.00821917808219178,
	//   "leadTimeForChanges":12165952333333333
	// }
	if !isNearBy(cliOutput.DeploymentFrequency, 0.00821917808219178, 0.01) {
		t.Errorf("deploymentFrequency should be near by 0.00821917808219178 but %v", cliOutput.DeploymentFrequency)
	}
	if !isNearBy(float64(cliOutput.LeadTimeForChanges), 12165952333333333, 0.01) {
		t.Errorf("deploymentFrequency should be near by but %v", cliOutput.DeploymentFrequency)
	}
}

func TestDefaultAppShouldRunWithoutOption(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	defaltApp := DefaultApp()
	testApp := &cli.App{
		Flags:  defaltApp.Flags,
		Action: defaltApp.Action,
		Writer: output,
	}

	err := testApp.Run([]string{"four-keys"})
	if err != nil {
		t.Errorf(err.Error())
	}

}

// isNearBy checks actual is in range of [expected*(1-epsilon), expected*(1+epsiolon)]
func isNearBy(actual float64, expected float64, epsilon float64) bool {
	return actual >= expected*(1-epsilon) && actual <= expected*(1+epsilon)
}

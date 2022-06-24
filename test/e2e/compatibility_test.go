package e2e

import (
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const shipyard = `apiVersion: "spec.keptn.sh/0.2.0"
kind: "Shipyard"
metadata:
  name: "shipyard"
spec:
  stages:
    - name: "dev"
      sequences:
        - name: "delivery"
          tasks:
            - name: "deployment"
              properties:
                deploymentstrategy: "user_managed"
`

func Test_CreateAndDeleteProject(t *testing.T) {
	if !isE2ETestingAllowed() {
		t.Skip("Skipping Test_ActionTriggered, not allowed by environment")
	}

	// Just test if we can connect to the cluster
	clientset, err := keptnutils.GetClientset(false)
	require.NoError(t, err)
	assert.NotNil(t, clientset)

	// Create a new Keptn api for the use of the E2E test
	keptnAPI := NewKeptAPI(readKeptnConnectionDetailsFromEnv())

	err = keptnAPI.CreateProject("e2e-gitea-test-project", []byte(shipyard))
	require.NoError(t, err)

	// TODO: test if the project exists in gitea

	err = keptnAPI.DeleteProject("e2e-gitea-test-project")
	require.NoError(t, err)

	// TODO: test if the project has been deleted from gitea
}

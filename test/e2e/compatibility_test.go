package e2e

import (
	"code.gitea.io/sdk/gitea"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"keptn-sandbox/keptn-gitea-provisioner/pkg/provisioner"
	"net/http"
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

	// Create the Gitea client from the environment
	giteaDetails := readGiteaConnectionDetailsFromEnv()
	client, err := gitea.NewClient(giteaDetails.Endpoint,
		gitea.SetBasicAuth(giteaDetails.Username, giteaDetails.Password),
	)
	require.NoError(t, err, "unable to connect to gitea")

	// Note: while projectName can be chosen freely, the projectUser is dependent on the Keptn namespace
	projectName := "e2e-gitea-test-project"
	projectUser := provisioner.DefaultKeptnNamespace

	// Create a repository
	err = keptnAPI.CreateProject(projectName, []byte(shipyard))
	require.NoError(t, err)

	// Check if the repository exists in the Gitea upstream server
	_, r, err := client.GetRepo(projectUser, projectName)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, r.StatusCode)

	// Delete the repository in Keptn
	err = keptnAPI.DeleteProject(projectName)
	require.NoError(t, err)

	// Repository must also be deleted from upstream Gitea server
	_, r, err = client.GetRepo(projectUser, projectName)
	require.Error(t, err)
	require.Equal(t, http.StatusNotFound, r.StatusCode)
}

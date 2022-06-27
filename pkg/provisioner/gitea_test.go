package provisioner

import (
	"code.gitea.io/sdk/gitea"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"keptn-sandbox/keptn-gitea-provisioner/pkg/keptn"
	"keptn-sandbox/keptn-gitea-provisioner/pkg/provisioner/fake"
	"net/http"
	"testing"
)

func createResponse(statusCode int) *gitea.Response {
	return &gitea.Response{
		Response: &http.Response{
			StatusCode: statusCode,
		},
	}
}

func TestGiteaProvisioner_GetUsername(t *testing.T) {
	prefixes := []string{"", "keptn-"}
	for _, prefix := range prefixes {
		giteaProvisioner := GiteaProvisioner{
			UsernamePrefix: prefix,
		}

		name := giteaProvisioner.GetUsername("example-namespace")
		assert.Equal(t, prefix+"example-namespace", name)
	}
}

func TestGiteaProvisioner_GetUsernameWithoutNamspace(t *testing.T) {
	giteaProvisioner := GiteaProvisioner{}
	assert.Equal(t, DefaultKeptnNamespace, giteaProvisioner.GetUsername(""))
}

func TestGiteaProvisioner_GetProjectName(t *testing.T) {
	prefixes := []string{"", "keptn_project-"}
	for _, prefix := range prefixes {
		giteaProvisioner := GiteaProvisioner{
			ProjectPrefix: prefix,
		}

		name := giteaProvisioner.GetProjectName("example-project")
		assert.Equal(t, prefix+"example-project", name)
	}
}

func TestGiteaProvisioner_CreateUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	giteaClient := fake.NewMockGiteaClient(mockCtrl)
	giteaProvisioner := GiteaProvisioner{
		client: giteaClient,
	}

	giteaClient.EXPECT().GetUserInfo("keptn").Times(1).Return(nil, createResponse(http.StatusNotFound), nil)
	giteaClient.EXPECT().AdminCreateUser(gomock.Any()).Times(1).Return(nil, createResponse(http.StatusCreated), nil)

	user, err := giteaProvisioner.CreateUser("keptn")
	require.NoError(t, err)
	require.Equal(t, "keptn", user)
}

func TestGiteaProvisioner_CreateRepository(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	giteaClient := fake.NewMockGiteaClient(mockCtrl)
	giteaProvisioner := GiteaProvisioner{
		client: giteaClient,
	}

	namespace := "keptn"
	repository := gitea.Repository{
		CloneURL: "http://some-gitea.repo:3000/keptn/repository",
	}

	giteaClient.EXPECT().AdminCreateRepo(namespace, gomock.Any()).Times(1).Return(&repository, createResponse(http.StatusCreated), nil)

	repo, err := giteaProvisioner.CreateRepository(namespace, "repository")
	require.NoError(t, err)
	require.Equal(t, repo, repository.CloneURL)
}

func TestGiteaProvisioner_UnableToCreateRepository(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	giteaClient := fake.NewMockGiteaClient(mockCtrl)
	giteaProvisioner := GiteaProvisioner{
		client: giteaClient,
	}

	namespace := "keptn"
	giteaClient.EXPECT().AdminCreateRepo(namespace, gomock.Any()).Times(1).Return(nil, createResponse(http.StatusConflict), nil)

	repo, err := giteaProvisioner.CreateRepository(namespace, "repository")
	require.Error(t, err)
	require.ErrorIs(t, ErrRepositoryAlreadyExists, err)
	require.Equal(t, "", repo)
}

func TestGiteaProvisioner_CreateToken(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	giteaClient := fake.NewMockGiteaClient(mockCtrl)
	giteaProvisioner := GiteaProvisioner{
		newClientFunc: func(url string, options ...gitea.ClientOption) (GiteaClient, error) {
			assert.Len(t, options, 2)
			return giteaClient, nil
		},
	}

	expectedParameters := gitea.CreateAccessTokenOption{
		Name: "some-keptn-project",
	}
	expectedToken := &gitea.AccessToken{
		Token: "12345670091-1230542347",
	}

	giteaClient.EXPECT().CreateAccessToken(expectedParameters).Times(1).Return(expectedToken, createResponse(http.StatusCreated), nil)

	token, err := giteaProvisioner.CreateToken("test", "some-keptn-project")
	require.NoError(t, err)
	require.Equal(t, "12345670091-1230542347", token)
}

func TestGiteaProvisioner_DeleteRepository(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	giteaClient := fake.NewMockGiteaClient(mockCtrl)
	giteaProvisioner := GiteaProvisioner{
		client: giteaClient,
		newClientFunc: func(url string, options ...gitea.ClientOption) (GiteaClient, error) {
			assert.Len(t, options, 2)
			return giteaClient, nil
		},
	}

	giteaClient.EXPECT().DeleteRepo("some-username", "project1").Times(1).Return(createResponse(http.StatusOK), nil)
	giteaClient.EXPECT().DeleteAccessToken("project1").Times(1).Return(nil, nil)
	giteaClient.EXPECT().ListMyRepos(gitea.ListReposOptions{}).Times(1).Return([]*gitea.Repository{}, createResponse(http.StatusOK), nil)
	giteaClient.EXPECT().AdminDeleteUser("some-username").Times(1).Return(createResponse(http.StatusNoContent), nil)

	err := giteaProvisioner.DeleteRepository("some-username", "project1")
	require.NoError(t, err)
}

func TestGiteaProvisioner_ProvisionRepository(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	giteaClient := fake.NewMockGiteaClient(mockCtrl)
	giteaProvisioner := GiteaProvisioner{
		client: giteaClient,
		newClientFunc: func(url string, options ...gitea.ClientOption) (GiteaClient, error) {
			assert.Len(t, options, 2)
			return giteaClient, nil
		},
	}

	repository := gitea.Repository{
		CloneURL: "http://some-gitea.repo:3000/user/some-keptn-project",
	}

	tokenOptions := gitea.CreateAccessTokenOption{
		Name: "some-keptn-project",
	}
	expectedToken := &gitea.AccessToken{
		Token: "12345670091-1230542347",
	}

	giteaClient.EXPECT().GetUserInfo("user").Times(1).Return(nil, createResponse(http.StatusNotFound), nil)
	giteaClient.EXPECT().AdminCreateUser(gomock.Any()).Times(1).Return(nil, createResponse(http.StatusCreated), nil)
	giteaClient.EXPECT().AdminCreateRepo("user", gomock.Any()).Times(1).Return(&repository, createResponse(http.StatusCreated), nil)
	giteaClient.EXPECT().CreateAccessToken(tokenOptions).Times(1).Return(expectedToken, createResponse(http.StatusCreated), nil)

	provisionRepository, err := giteaProvisioner.ProvisionRepository("user", "some-keptn-project")
	require.NoError(t, err)

	expectedResult := keptn.ProvisionResponse{
		GitRemoteURL: "http://some-gitea.repo:3000/user/some-keptn-project",
		GitToken:     "12345670091-1230542347",
		GitUser:      "user",
	}

	require.Equal(t, expectedResult, *provisionRepository)
}

func TestGiteaProvisioner_ProvisionRepositoryConflict(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	giteaClient := fake.NewMockGiteaClient(mockCtrl)
	giteaProvisioner := GiteaProvisioner{
		client: giteaClient,
		newClientFunc: func(url string, options ...gitea.ClientOption) (GiteaClient, error) {
			assert.Len(t, options, 2)
			return giteaClient, nil
		},
	}

	giteaClient.EXPECT().GetUserInfo("user").Times(1).Return(nil, createResponse(http.StatusNotFound), nil)
	giteaClient.EXPECT().AdminCreateUser(gomock.Any()).Times(1).Return(nil, createResponse(http.StatusCreated), nil)
	giteaClient.EXPECT().AdminCreateRepo("user", gomock.Any()).Times(1).Return(nil, createResponse(http.StatusConflict), nil)
	giteaClient.EXPECT().CreateAccessToken(gomock.Any()).Times(0)

	provisionRepository, err := giteaProvisioner.ProvisionRepository("user", "some-keptn-project")
	require.ErrorIs(t, err, ErrRepositoryAlreadyExists)
	require.Nil(t, provisionRepository)
}

func TestNewGiteaProvisioner(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	giteaClient := fake.NewMockGiteaClient(mockCtrl)
	mockBuilder := func(url string, options ...gitea.ClientOption) (GiteaClient, error) {
		return giteaClient, nil
	}

	giteaProvisioner, err := NewGiteaProvisioner(
		"http://gitea.endpoint:3000/",
		"admin",
		"secret",
		&GiteaProvisionerOptions{
			UsernamePrefix:  "user-",
			UserEmailDomain: "domain.local",
			ProjectPrefix:   "project-",
			TokenPrefix:     "token-",
			ClientBuilder:   mockBuilder,
		})

	require.NoError(t, err)
	assert.Equal(t, giteaClient, giteaProvisioner.client)
	assert.Equal(t, "http://gitea.endpoint:3000/", giteaProvisioner.endpoint)
	assert.Equal(t, "user-", giteaProvisioner.UsernamePrefix)
	assert.Equal(t, "domain.local", giteaProvisioner.UserEmailDomain)
	assert.Equal(t, "project-", giteaProvisioner.ProjectPrefix)
	assert.Equal(t, "token-", giteaProvisioner.TokenPrefix)
}

func TestGiteaProvisioner_ProvisionRepositoryEmptyProject(t *testing.T) {
	giteaProvisioner := GiteaProvisioner{}
	repo, err := giteaProvisioner.ProvisionRepository("", "")
	require.ErrorIs(t, err, ErrInvalidRequest)
	require.Nil(t, repo)
}

func TestGiteaProvisioner_DeleteRepositoryEmptyProject(t *testing.T) {
	giteaProvisioner := GiteaProvisioner{}
	err := giteaProvisioner.DeleteRepository("", "")
	require.ErrorIs(t, err, ErrInvalidRequest)
}

func TestGiteaProvisioner_CreateUserGiteaProblem(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	giteaClient := fake.NewMockGiteaClient(mockCtrl)
	giteaProvisioner := GiteaProvisioner{
		client: giteaClient,
	}

	giteaClient.EXPECT().GetUserInfo("keptn").Times(1).Return(nil, createResponse(http.StatusNotFound), nil)
	giteaClient.EXPECT().AdminCreateUser(gomock.Any()).Times(1).Return(nil, createResponse(http.StatusUnprocessableEntity), nil)

	user, err := giteaProvisioner.CreateUser("keptn")
	require.Error(t, err)
	require.Equal(t, "", user)
}

func TestGiteaProvisioner_CreateRepositoryGiteaProblem(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	giteaClient := fake.NewMockGiteaClient(mockCtrl)
	giteaProvisioner := GiteaProvisioner{
		client: giteaClient,
	}

	namespace := "keptn"
	repository := gitea.Repository{
		CloneURL: "http://some-gitea.repo:3000/keptn/repository",
	}

	giteaClient.EXPECT().AdminCreateRepo(namespace, gomock.Any()).Times(1).Return(&repository, createResponse(http.StatusUnprocessableEntity), nil)

	repo, err := giteaProvisioner.CreateRepository(namespace, "repository")
	require.Error(t, err)
	require.Equal(t, "", repo)
}

func TestGiteaProvisioner_CreateTokenGiteaProblem(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	giteaClient := fake.NewMockGiteaClient(mockCtrl)
	giteaProvisioner := GiteaProvisioner{
		newClientFunc: func(url string, options ...gitea.ClientOption) (GiteaClient, error) {
			assert.Len(t, options, 2)
			return giteaClient, nil
		},
	}

	expectedParameters := gitea.CreateAccessTokenOption{
		Name: "some-keptn-project",
	}
	expectedToken := &gitea.AccessToken{
		Token: "12345670091-1230542347",
	}

	giteaClient.EXPECT().CreateAccessToken(expectedParameters).Times(1).Return(expectedToken, createResponse(http.StatusUnprocessableEntity), nil)

	token, err := giteaProvisioner.CreateToken("test", "some-keptn-project")
	require.Error(t, err)
	require.Equal(t, "", token)
}

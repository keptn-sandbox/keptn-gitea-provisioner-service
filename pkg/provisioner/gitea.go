package provisioner

import (
	"errors"
	"fmt"
	"net/http"

	"code.gitea.io/sdk/gitea"

	"keptn-sandbox/keptn-gitea-provisioner/pkg/utils"
)

// DefaultPasswordLength indicates the length of the generated passwords
const DefaultPasswordLength = 32

// DefaultKeptnNamespace is used when no additional keptn namespace was defined in the request
const DefaultKeptnNamespace = "keptn"

// DefaultUserEmailDomain is the default E-Mail domain used for users
const DefaultUserEmailDomain = "keptn-gitea-auto-provisioner.local"

// GiteaClient represents the interface of the Gitea client that is needed for the provisioner
type GiteaClient interface {
	GetUserInfo(user string) (*gitea.User, *gitea.Response, error)
	AdminCreateUser(opt gitea.CreateUserOption) (*gitea.User, *gitea.Response, error)
	AdminCreateRepo(username string, opt gitea.CreateRepoOption) (*gitea.Repository, *gitea.Response, error)
	DeleteRepo(username string, repository string) (*gitea.Response, error)
	CreateAccessToken(opt gitea.CreateAccessTokenOption) (*gitea.AccessToken, *gitea.Response, error)
	DeleteAccessToken(value interface{}) (*gitea.Response, error)
}

//go:generate mockgen -destination=fake/gitea_mock.go -package=fake . GiteaClient

// The GiteaProvisioner structure implements the Provisioner interface and provides functionality for creating, deleting
// the different resources in a Gitea
type GiteaProvisioner struct {
	endpoint        string
	credentials     gitea.ClientOption
	client          GiteaClient
	newClientFunc   func(url string, options ...gitea.ClientOption) (GiteaClient, error)
	UsernamePrefix  string
	UserEmailDomain string
	ProjectPrefix   string
	TokenPrefix     string
}

// GiteaProvisionerOptions defines additional options than can be specified when creating a GiteaProvisioner
type GiteaProvisionerOptions struct {
	UsernamePrefix  string
	UserEmailDomain string
	ProjectPrefix   string
	TokenPrefix     string
	ClientBuilder   func(url string, options ...gitea.ClientOption) (GiteaClient, error)
}

// NewGiteaProvisioner creates a new gitea provisioner service with the given credentials and options
func NewGiteaProvisioner(giteaEndpoint string, adminUsername string, adminPassword string, options *GiteaProvisionerOptions) (*GiteaProvisioner, error) {
	clientBuilder := func(url string, options ...gitea.ClientOption) (GiteaClient, error) {
		return gitea.NewClient(url, options...)
	}

	if options.ClientBuilder != nil {
		clientBuilder = options.ClientBuilder
	}

	clientCredentials := gitea.SetBasicAuth(adminUsername, adminPassword)
	giteaClient, err := clientBuilder(giteaEndpoint, clientCredentials)
	if err != nil {
		return nil, fmt.Errorf("unable to create Gitea Client: %w", err)
	}

	provisioner := GiteaProvisioner{
		endpoint:      giteaEndpoint,
		credentials:   clientCredentials,
		client:        giteaClient,
		newClientFunc: clientBuilder,
	}

	// If options are set, apply them to the provisioner
	if options != nil {
		provisioner.UsernamePrefix = options.UsernamePrefix
		provisioner.UserEmailDomain = options.UserEmailDomain
		provisioner.ProjectPrefix = options.ProjectPrefix
		provisioner.TokenPrefix = options.TokenPrefix
	}

	// Make sure the e-mail domain is set, because otherwise account creation will fail
	if provisioner.UserEmailDomain == "" {
		provisioner.UserEmailDomain = DefaultUserEmailDomain
	}

	return &provisioner, nil
}

// CreateUser creates a user if it doesn't exist already for the given Keptn namespace
func (h *GiteaProvisioner) CreateUser(namespace string) (string, error) {

	// Generate a user
	username := h.GetUsername(namespace)
	password := utils.GenerateRandomString(DefaultPasswordLength)

	// Check if user
	user, r, err := h.client.GetUserInfo(username)
	if err != nil && r == nil {
		return "", fmt.Errorf("unable to get user info for user %s: %w", username, err)
	}

	// If no user was found, we have to create the user
	if user == nil || r.StatusCode == http.StatusNotFound {
		passwordChangePolicy := false

		_, r, err := h.client.AdminCreateUser(gitea.CreateUserOption{
			LoginName:          username,
			Username:           username,
			FullName:           username,
			Email:              fmt.Sprintf("%s@%s", username, h.UserEmailDomain),
			Password:           password,
			MustChangePassword: &passwordChangePolicy,
			SendNotify:         false,
		})

		if err != nil && r == nil {
			return "", fmt.Errorf("unable to create user %s: %w", username, err)
		}

		// Possible status codes: 400, 403, 422
		if r.StatusCode != http.StatusCreated {
			return "", fmt.Errorf("unable to create user %s, received unexpected status code: %d",
				username, r.StatusCode,
			)
		}
	}

	return username, nil
}

// CreateToken creates an access token that has read/write privileges for the given project
func (h *GiteaProvisioner) CreateToken(namespace string, project string) (string, error) {
	// Note: we must change the client to use a different user:
	userClient, err := h.newClientFunc(h.endpoint, h.credentials, gitea.SetSudo(h.GetUsername(namespace)))
	if err != nil {
		return "", fmt.Errorf("unable to create gitea client: %w", err)
	}

	token, r, err := userClient.CreateAccessToken(gitea.CreateAccessTokenOption{
		Name: h.GetAccessTokenName(project),
	})
	if err != nil {
		return "", fmt.Errorf("unable to create access token: %w", err)
	}

	if r.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("recieved unkown http status code: %d", r.StatusCode)
	}

	return token.Token, nil
}

// CreateRepository creates a repository in the Gitea server
func (h *GiteaProvisioner) CreateRepository(namespace string, project string) (string, error) {
	projectName := h.GetProjectName(project)
	projectDesc := fmt.Sprintf(
		"Repository was automatically provisioned by keptn-gitea-Provisioner for project %s",
		projectName,
	)

	// Note: Keptn requires a completely empty git repository where the default branch is set to master
	repo, r, err := h.client.AdminCreateRepo(h.GetUsername(namespace), gitea.CreateRepoOption{
		Name:          projectName,
		Description:   projectDesc,
		Private:       true,
		IssueLabels:   "",
		AutoInit:      false,
		Template:      false,
		Gitignores:    "",
		License:       "",
		Readme:        "",
		DefaultBranch: "master",
		TrustModel:    gitea.TrustModelDefault,
	})

	// Error while talking to gitea, upstream failed or something else
	if err != nil && r == nil {
		return "", fmt.Errorf("unable to create project \"%s\": %w", projectName, err)
	}

	// Project already exists, relay the status code only
	if r.StatusCode == http.StatusConflict {
		return "", ErrRepositoryAlreadyExists
	}

	// Possible status codes: 403, 404, 422
	if r.StatusCode != http.StatusCreated {
		return "", fmt.Errorf(
			"recieved unexpected status code %d while creating repository %s for namespace %s",
			r.StatusCode, project, namespace,
		)
	}

	return repo.CloneURL, nil
}

// GetUsername returns the username that is used by the gitea upstream server to identify a Keptn namespace
func (h *GiteaProvisioner) GetUsername(namespace string) string {
	// Use default Keptn namespace if no one is defined, to avoid creating users that
	// have more or less an empty name if no prefix was defined
	if namespace == "" {
		namespace = DefaultKeptnNamespace
	}

	return fmt.Sprintf("%s%s", h.UsernamePrefix, namespace)
}

// GetProjectName returns the name of the project in the gitea upstream repository
func (h *GiteaProvisioner) GetProjectName(project string) string {
	return fmt.Sprintf("%s%s", h.ProjectPrefix, project)
}

// GetAccessTokenName returns the name of the access token that as read/write privileges for the specified project
func (h *GiteaProvisioner) GetAccessTokenName(project string) string {
	return fmt.Sprintf("%s%s", h.TokenPrefix, project)
}

// DeleteRepository deletes a given repository and all associated resources that where created with that repository
func (h *GiteaProvisioner) DeleteRepository(namespace string, project string) error {

	username := h.GetUsername(namespace)
	accessToken := h.GetAccessTokenName(project)

	r, err := h.client.DeleteRepo(username, project)
	if err != nil && r == nil {
		return fmt.Errorf("unable to delete the repository: %w", err)
	}

	// Project does not exist, relay the status code only
	if r.StatusCode == http.StatusNotFound {
		return ErrRepositoryDoesNotExist
	}

	// Note: to delete a access token we have to use sudo mode:
	userClient, err := h.newClientFunc(h.endpoint, h.credentials, gitea.SetSudo(username))
	if err != nil {
		return fmt.Errorf("unable create gitea client: %w", err)
	}

	_, err = userClient.DeleteAccessToken(accessToken)
	if err != nil {
		return fmt.Errorf("unable to delete the access token: ")
	}

	// TODO: If the user doesn't have any repositories we might want to delete the user
	return nil
}

// ProvisionRepository provisions the Gitea repository with the given Keptn namespace and project name, this includes also
// the creation of needed resources such as users and access tokens.
func (h *GiteaProvisioner) ProvisionRepository(namespace string, project string) (*ProvisionResponse, error) {

	if project == "" {
		return nil, fmt.Errorf("unable to create project with an empty name")
	}

	if _, err := h.CreateUser(namespace); err != nil {
		return nil, fmt.Errorf("unable to create user: %s\n", err.Error())
	}

	repository, err := h.CreateRepository(namespace, project)
	if err != nil {

		if errors.Is(err, ErrRepositoryAlreadyExists) {
			return nil, ErrRepositoryAlreadyExists
		}

		return nil, fmt.Errorf("unable to create repository: %w", err)
	}

	username := h.GetUsername(namespace)
	token, err := h.CreateToken(namespace, project)
	if err != nil {
		return nil, fmt.Errorf("unable to create token: %w", err)
	}

	return &ProvisionResponse{
		GitRemoteURL: repository,
		GitToken:     token,
		GitUser:      username,
	}, nil
}

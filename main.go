package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kelseyhightower/envconfig"

	"keptn-sandbox/keptn-gitea-provisioner/pkg/provisioner"
)

var /*const*/ env envConfig

type envConfig struct {
	// Port on which the provisioner listens on
	Port int `envconfig:"RCV_PORT" default:"8080"`
	// The GiteaEndpoint is a required environment variable that describes the URL of the Gitea endpoint
	GiteaEndpoint string `envconfig:"GITEA_ENDPOINT" required:"true"`
	// GiteaUser is a required environment variable, which should be the username of the admin user
	GiteaUser string `envconfig:"GITEA_USER" required:"true"`
	// GiteaPassword is a required environment variable, which should be the password of the admin user
	GiteaPassword string `envconfig:"GITEA_PASSWORD" required:"true"`
	// UsernamePrefix defines the prefix that should be used when creating users in Gitea
	UsernamePrefix string `envconfig:"USERNAME_PREFIX"`
	// UserEmailDomain defines the prefix that should be used when creating users in Gitea
	UserEmailDomain string `envconfig:"USER_EMAIL_DOMAIN"`
	// ProjectPrefix defines the prefix that should be used when creating projects in Gitea
	ProjectPrefix string `envconfig:"PROJECT_PREFIX"`
	// TokenPrefix defines the prefix that should be used when creating tokens in Gitea
	TokenPrefix string `envconfig:"TOKEN_PREFIX"`
}

func main() {
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("Failed to process env var: %s", err)
	}

	giteaOptions := provisioner.GiteaProvisionerOptions{
		UsernamePrefix:  env.UsernamePrefix,
		UserEmailDomain: env.UserEmailDomain,
		ProjectPrefix:   env.ProjectPrefix,
		TokenPrefix:     env.TokenPrefix,
	}

	repoProvisioner, err := provisioner.NewGiteaProvisioner(env.GiteaEndpoint, env.GiteaUser, env.GiteaPassword, &giteaOptions)
	if err != nil {
		log.Fatalf("Unable to create gitea provisioner: %s", err)
	}

	provisionerHandler := provisioner.ProvisionHandler{
		Provisioner: repoProvisioner,
	}

	http.HandleFunc("/repository", provisionerHandler.HandleProvisionRepoRequest)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", env.Port), nil); err != nil {
		log.Fatalf("Failed to serve git provisioning service at endpoint: %s", err)
	}

	os.Exit(0)
}

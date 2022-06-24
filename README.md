# keptn-gitea-provisioner-service

:warning: Work in progress :warning:
---

![GitHub release (latest by date)](https://img.shields.io/github/v/release/keptn-sandbox/keptn-gitea-provisioner-service)
![Go Report Card](https://goreportcard.com/badge/github.com/keptn-sandbox/keptn-gitea-provisioner-service)

This repository contains a reference implementation for a Keptn service that is able to auto-provision git repositories
in Gitea. This is done done by utilizing the extension points for [automatic git provisioning](https://keptn.sh/docs/0.16.x/api/git_provisioning/) in Keptn.

## Compatibility Matrix

| Keptn Version* | [Keptn-Service-Template-Go Docker Image](https://hub.docker.com/r/keptn-sandbox/keptn-gitea-provisioner-service/tags) |
|:--------------:|:---------------------------------------------------------------------------------------------------------------------:|
|     0.16.x     |                                  keptn-sandbox/keptn-gitea-provisioner-service:0.1.0                                  |

\* This is the Keptn version we aim to be compatible with. Other versions should work too, but there is no guarantee.

## Quickstart

* The keptn-gitea-provisioner-service can be easily installed with helm and is able to create the credentials for the Gitea
instance:
  ```bash
  #!/bin/bash
  NAMESPACE=default
   GITEA_ADMIN_USERNAME=#Define a username for the admin
   GITEA_ADMIN_PASSWORD=#Define a password for the admin
  
  helm install -n ${NAMESPACE} keptn-gitea-provisioner-service chart/ \
    --set admin.create=true \
    --set admin.username=${GITEA_ADMIN_USERNAME} \
    --set admin.password=${GITEA_ADMIN_PASSWORD}
    
  ```
  
  *Note*: You can re-use by omitting the set parameters of the helm installation


* If there is no Gitea instance installed, an appropriate instance can be created with the following bash script, otherwise it's sufficient to 
just create the kubernetes secret.
  ```bash
  #!/bin/bash
  NAMESPACE=default
  
  # Add the gitea helm charts and install gitea to the cluster
  helm repo add gitea-charts https://dl.gitea.io/charts/
  helm repo update
  
  # Install gitea in the 
  helm install gitea gitea-charts/gitea \
      --set memcached.enabled=false \
      --set postgresql.enabled=false \
      --set gitea.config.database.DB_TYPE=sqlite3 \
      --set gitea.admin.existingSecret=gitea-admin-secret \
      --set gitea.config.server.OFFLINE_MODE=true \
      --set gitea.config.server.ROOT_URL=http://gitea-http.${NAMESPACE}:3000/
  ```

* Keptn must be configured to use the keptn-gitea-provisioner-service to automatically provision git repositories:
  ```bash
  #!/bin/bash
  NAMESPACE=default
  
  helm upgrade -n keptn keptn \
    --set "control-plane.features.automaticProvisioningURL=http://keptn-gitea-provisioner-service.${NAMESPACE}"
  ```

### Deploy in your Kubernetes cluster

To deploy the current version of the *keptn-gitea-provisioner-service* in your Keptn Kubernetes cluster use the [`helm chart`](chart/Chart.yaml) file,
for example:

```console
helm install -n keptn keptn-gitea-provisioner-service chart/
```

This should install the `keptn-gitea-provisioner-service` together with a Keptn `distributor` into the `keptn` namespace, which you can verify using

```console
kubectl -n keptn get deployment keptn-gitea-provisioner-service -o wide
kubectl -n keptn get pods -l run=keptn-gitea-provisioner-service
```

### Up- or Downgrading

Adapt and use the following command in case you want to up- or downgrade your installed version (specified by the `$VERSION` placeholder):

```console
helm upgrade -n keptn --set image.tag=$VERSION keptn-gitea-provisioner-service chart/
```

### Uninstall

To delete a deployed *keptn-gitea-provisioner-service*, use the file `deploy/*.yaml` files from this repository and delete the Kubernetes resources:

```console
helm uninstall -n keptn keptn-gitea-provisioner-service
```

## Development

Development can be conducted using any GoLang compatible IDE/editor (e.g., Jetbrains GoLand, VSCode with Go plugins).

It is recommended to make use of branches as follows:

* `main` contains the latest potentially unstable version
* `release-*` contains a stable version of the service (e.g., `release-0.1.0` contains version 0.1.0)
* create a new branch for any changes that you are working on, e.g., `feature/my-cool-stuff` or `bug/overflow`
* once ready, create a pull request from that branch back to the `main`/`master` branch

When writing code, it is recommended to follow the coding style suggested by the [Golang community](https://github.com/golang/go/wiki/CodeReviewComments).

### Common tasks

* Build the binary: `go build -ldflags '-linkmode=external' -v -o keptn-gitea-provisioner-service`
* Run tests: `go test -race -v ./...`
* Build the docker image: `docker build . -t keptn-sandbox/keptn-gitea-provisioner-service:dev` (Note: Ensure that you use the correct DockerHub account/organization)
* Run the docker image locally: `docker run --rm -it -p 8080:8080 keptn-sandbox/keptn-gitea-provisioner-service:dev`
* Push the docker image to DockerHub: `docker push keptn-sandbox/keptn-gitea-provisioner-service:dev` (Note: Ensure that you use the correct DockerHub account/organization)
* Deploy the service using `kubectl`: `kubectl apply -f deploy/`
* Delete/undeploy the service using `kubectl`: `kubectl delete -f deploy/`
* Watch the deployment using `kubectl`: `kubectl -n keptn get deployment keptn-gitea-provisioner-service -o wide`
* Get logs using `kubectl`: `kubectl -n keptn logs deployment/keptn-gitea-provisioner-service -f`
* Watch the deployed pods using `kubectl`: `kubectl -n keptn get pods -l run=keptn-gitea-provisioner-service`
* Deploy the service using [Skaffold](https://skaffold.dev/): `skaffold run --default-repo=your-docker-registry --tail` (Note: Replace `your-docker-registry` with your container image registry (defaults to ghcr.io/keptn-sandbox/keptn-gitea-provisioner-service); also make sure to adapt the image name in [skaffold.yaml](skaffold.yaml))


## Automation

### GitHub Actions: Automated Pull Request Review

This repo uses [reviewdog](https://github.com/reviewdog/reviewdog) for automated reviews of Pull Requests. 

You can find the details in [.github/workflows/reviewdog.yml](.github/workflows/reviewdog.yml).

### GitHub Actions: Unit Tests

This repo has automated unit tests for pull requests. 

You can find the details in [.github/workflows/CI.yml](.github/workflows/CI.yml).

### GH Actions/Workflow: Build Docker Images

This repo uses GH Actions and Workflows to test the code and automatically build docker images.

Docker Images are automatically pushed based on the configuration done in [.ci_env](.ci_env) and the two [GitHub Secrets](https://github.com/keptn-sandbox/keptn-gitea-provisioner-service/settings/secrets/actions)
* `REGISTRY_USER` - your DockerHub username
* `REGISTRY_PASSWORD` - a DockerHub [access token](https://hub.docker.com/settings/security) (alternatively, your DockerHub password)

## How to release a new version of this service

It is assumed that the current development takes place in the `main` branch (either via Pull Requests or directly).

Once you're ready, go to the Actions tab on GitHub, select Pre-Release or Release, and run the action.


## License

Please find more information in the [LICENSE](LICENSE) file.

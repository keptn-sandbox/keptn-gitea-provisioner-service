# keptn-gitea-provisioner-service

:warning: Work in progress :warning:
---

![GitHub release (latest by date)](https://img.shields.io/github/v/release/keptn-sandbox/keptn-gitea-provisioner-service)
[![Go Report Card](https://goreportcard.com/badge/github.com/keptn-sandbox/keptn-gitea-provisioner-service)](https://goreportcard.com/report/github.com/keptn-sandbox/keptn-gitea-provisioner-service)

This implements a keptn-gitea-provisioner-service for Keptn. If you want to learn more about Keptn visit us on [keptn.sh](https://keptn.sh)

## Compatibility Matrix

*Please fill in your versions accordingly*

| Keptn Version* | [Keptn-Service-Template-Go Docker Image](https://hub.docker.com/r/keptn-sandbox/keptn-gitea-provisioner-service/tags) |
|:--------------:|:---------------------------------------------------------------------------------------------------------------:|
|   0.6 - 0.8    |                                  keptn-sandbox/keptn-gitea-provisioner-service:0.8.3                                  |
|     0.10.x     |                                 keptn-sandbox/keptn-gitea-provisioner-service:0.10.0                                  |
|     0.11.x     |                                 keptn-sandbox/keptn-gitea-provisioner-service:0.11.4                                  |
|     0.12.x     |                                 keptn-sandbox/keptn-gitea-provisioner-service:0.12.2                                  |
|     0.13.x     |                                 keptn-sandbox/keptn-gitea-provisioner-service:0.13.0                                  |
|     0.14.x     |                                 keptn-sandbox/keptn-gitea-provisioner-service:0.14.0                                  |

\* This is the Keptn version we aim to be compatible with. Other versions should work too, but there is no guarantee.

**Note**: Versions compatible with Keptn 0.14.x and newer are not backward compatible due to a change in NATS cluster name
(see https://github.com/keptn/keptn/releases/tag/0.14.1 for more info about the breaking change).

## Installation

The *keptn-gitea-provisioner-service* can be installed as a part of [Keptn's uniform](https://keptn.sh).

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

* `main`/`master` contains the latest potentially unstable version
* `release-*` contains a stable version of the service (e.g., `release-0.1.0` contains version 0.1.0)
* create a new branch for any changes that you are working on, e.g., `feature/my-cool-stuff` or `bug/overflow`
* once ready, create a pull request from that branch back to the `main`/`master` branch

When writing code, it is recommended to follow the coding style suggested by the [Golang community](https://github.com/golang/go/wiki/CodeReviewComments).

### Where to start

If you don't care about the details, your first entrypoint is [eventhandlers.go](eventhandlers.go). Within this file 
 you can add implementation for pre-defined Keptn Cloud events.
 
To better understand all variants of Keptn CloudEvents, please look at the [Keptn Spec](https://github.com/keptn/spec).
 
If you want to get more insights into processing those CloudEvents or even defining your own CloudEvents in code, please 
 look into [main.go](main.go) (specifically `processKeptnCloudEvent`), [chart/values.yaml](chart/values.yaml),
 consult the [Keptn docs](https://keptn.sh/docs/) as well as existing [Keptn Core](https://github.com/keptn/keptn) and
 [Keptn Contrib](https://github.com/keptn-contrib/) services.

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


### Testing Cloud Events

We have dummy cloud-events in the form of [RFC 2616](https://ietf.org/rfc/rfc2616.txt) requests in the [test-events/](test-events/) directory. These can be easily executed using third party plugins such as the [Huachao Mao REST Client in VS Code](https://marketplace.visualstudio.com/items?itemName=humao.rest-client).

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

It is assumed that the current development takes place in the `main`/`master` branch (either via Pull Requests or directly).

Once you're ready, go to the Actions tab on GitHub, select Pre-Release or Release, and run the action.


## License

Please find more information in the [LICENSE](LICENSE) file.

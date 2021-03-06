# keptn-gitea-provisioner-service

![GitHub release (latest by date)](https://img.shields.io/github/v/release/keptn-sandbox/keptn-gitea-provisioner-service)
![Go Report Card](https://goreportcard.com/badge/github.com/keptn-sandbox/keptn-gitea-provisioner-service)

This repository contains a reference implementation for a [Keptn extension which auto-provisions git upstream repositories](https://keptn.sh/docs/0.16.x/api/git_provisioning/) in Gitea.

## Compatibility Matrix

| Keptn Version* | [Keptn-Service-Template-Go Docker Image](https://hub.docker.com/r/keptn-sandbox/keptn-gitea-provisioner-service/tags) |
|:--------------:|:---------------------------------------------------------------------------------------------------------------------:|
|   0.15, 0.16   |                                  keptn-sandbox/keptn-gitea-provisioner-service:0.1.0                                  |
|     0.17.0     |                                  keptn-sandbox/keptn-gitea-provisioner-service:0.1.1                                  |

\* This is the Keptn version we aim to be compatible with. Other versions should work too, but there is no guarantee.

## Quickstart

* The keptn-gitea-provisioner-service can be easily installed with helm and is able to create the credentials for the Gitea
instance:
  ```bash
  #!/bin/bash
  VERSION=0.1.0
  NAMESPACE=default
  GITEA_ENDPOINT="http://gitea-http.${NAMESPACE}:3000/"
   GITEA_ADMIN_USERNAME=#Define a username for the admin
   GITEA_ADMIN_PASSWORD=#Define a password for the admin
  
  helm install keptn-gitea-provisioner-service https://github.com/keptn-sandbox/keptn-gitea-provisioner-service/releases/download/${VERSION}/keptn-gitea-provisioner-service-${VERSION}.tgz \
        --set gitea.admin.create=true \
        --set gitea.admin.username=${GITEA_ADMIN_USERNAME} \
        --set gitea.admin.password=${GITEA_ADMIN_PASSWORD} \
        --set gitea.endpoint=${GITEA_ENDPOINT} \
        --wait
  ```
  
  *Note*: You can re-use existing credentials by omitting the set parameters of the helm installation; For a full list 
          of options that can be set see the chart [helm documentation](chart/README.md).

* In case you need a simple Gitea instance, you can create it with following bash script:
  ```bash
  #!/bin/bash
  NAMESPACE=default #Should be configured to match the GITEA_ENDPOINT environment variable when installing the provisioner
  
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

* Keptn must be configured to use the keptn-gitea-provisioner-service to automatically provision git repositories. The flag is different for Keptn version 0.16 and version 0.17 and onwards:
  ```bash
  #!/bin/bash
  NAMESPACE=default
  
  # Keptn 0.17
  helm upgrade -n keptn keptn keptn/keptn \
    --set "features.automaticProvisioning.serviceURL=http://keptn-gitea-provisioner-service.${NAMESPACE}"
  
  # Keptn 0.16
  helm upgrade -n keptn keptn keptn/keptn \
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

## Architecture

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

## License

Please find more information in the [LICENSE](LICENSE) file.

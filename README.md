# keptn-gitea-provisioner-service

![GitHub release (latest by date)](https://img.shields.io/github/v/release/keptn-sandbox/keptn-gitea-provisioner-service)
![Go Report Card](https://goreportcard.com/badge/github.com/keptn-sandbox/keptn-gitea-provisioner-service)

This repository contains a reference implementation for a Keptn service that is able to auto-provision git repositories
in Gitea. This is done by utilizing the extension points for [automatic git provisioning](https://keptn.sh/docs/0.16.x/api/git_provisioning/) in Keptn.

## Compatibility Matrix

| Keptn Version* | [Keptn-Service-Template-Go Docker Image](https://hub.docker.com/r/keptn-sandbox/keptn-gitea-provisioner-service/tags) |
|:--------------:|:---------------------------------------------------------------------------------------------------------------------:|
|     0.15.1     |                                  keptn-sandbox/keptn-gitea-provisioner-service:0.1.0                                  |

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
  
  *Note*: You can re-use existing credentials by omitting the set parameters of the helm installation; For a full list 
          of options that can be set see the chart [helm documentation](chart/README.md).

* If there is no Gitea instance installed, an appropriate instance can be created with the following bash script:
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

## License

Please find more information in the [LICENSE](LICENSE) file.

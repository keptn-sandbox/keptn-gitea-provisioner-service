keptn-gitea-provisioner-service
===========

Helm Chart for the keptn keptn-gitea-provisioner-service


## Configuration

The following table lists the configurable parameters of the keptn-gitea-provisioner-service chart and their default values.

| Parameter                       | Description                                                                        | Default                                                   |
|---------------------------------|------------------------------------------------------------------------------------|-----------------------------------------------------------|
| `image.repository`              | Container image name                                                               | `"ghcr.io/keptn-sandbox/keptn-gitea-provisioner-service"` |
| `image.pullPolicy`              | Kubernetes image pull policy                                                       | `"IfNotPresent"`                                          |
| `image.tag`                     | Container tag                                                                      | `""`                                                      |
| `service.enabled`               | Creates a kubernetes service for the keptn-gitea-provisioner-service               | `true`                                                    |
 | `gitea.endpoint`                | The endpoint URL of the Gitea server                                               | `http://gitea-http.default:3000/`                         |
 | `gitea.admin.create`            | Set to `true` if the admin user & password credentials should be saved to a secret | `false`                                                   |
| `gitea.admin.username`          | The username of the Gitea admin user                                               | ` `                                                       |
| `gitea.admin.password`          | The password of the Gitea admin user                                               | ` `                                                       |
| `gitea.options.usernamePrefix`  | A prefix that is used by the provisioner when creating users in Gitea              | ` `                                                       |
| `gitea.options.userEmailDomain` | The E-Mail domain that is used when creating users                                 | ` `                                                       |
| `gitea.options.projectPrefix`   | A prefix that is used by the provisioner when creating project in Gitea            | ` `                                                       |
| `gitea.options.tokenPrefix`     | A prefix that is used by the provisioner when creating tokens in Gitea             | ` `                                                       |
| `imagePullSecrets`              | Secrets to use for container registry credentials                                  | `[]`                                                      |
| `podAnnotations`                | Annotations to add to the created pods                                             | `{}`                                                      |
| `podSecurityContext`            | Set the pod security context (e.g. fsgroups)                                       | `{}`                                                      |
| `securityContext`               | Set the security context (e.g. runasuser)                                          | `{}`                                                      |
| `resources`                     | Resource limits and requests                                                       | `{}`                                                      |
| `nodeSelector`                  | Node selector configuration                                                        | `{}`                                                      |
| `tolerations`                   | Tolerations for the pods                                                           | `[]`                                                      |
| `affinity`                      | Affinity rules                                                                     | `{}`                                                      |


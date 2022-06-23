package provisioner

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

// ErrRepositoryAlreadyExists indicates that the repository already exists
var /*const*/ ErrRepositoryAlreadyExists = errors.New("the repository already exists")

// ErrRepositoryDoesNotExist indicates that the repository does not exist
var /*const*/ ErrRepositoryDoesNotExist = errors.New("the repository does not exist")

// ProvisionRequest represents the request body of Keptn which is used when requesting a new git repository
type ProvisionRequest struct {
	Project   string `json:"project"`
	Namespace string `json:"namespace"`
}

// ProvisionResponse represents the response body that Keptn is expecting as a response from the provisioning request
type ProvisionResponse struct {
	GitRemoteURL string `json:"gitRemoteURL"`
	GitToken     string `json:"gitToken"`
	GitUser      string `json:"gitUser"`
}

// Provisioner contains a set of methods that a repository provisioner must implement
type Provisioner interface {
	// DeleteRepository deletes the repository and all associated resources (e.g.: token)
	DeleteRepository(namespace string, project string) error
	// ProvisionRepository creates all required resources for the given request
	ProvisionRepository(namespace string, project string) (*ProvisionResponse, error)
}

// The ProvisionHandler provides the HandleProvisionRepoRequest method which can be used within a HTTPListener to process
// repository provision and deletion requests from Keptn
type ProvisionHandler struct {
	Provisioner Provisioner
}

// HandleProvisionRepoRequest handles a GET or POST http request and provisions or deletes the defined repository in the request
func (p *ProvisionHandler) HandleProvisionRepoRequest(w http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case http.MethodPost:
		p.handleProvisionRepository(w, req)
		break

	case http.MethodDelete:
		p.handleDeleteRepository(w, req)
		break

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

// decodeRequestBody decodes the body of the given http request into a keptn.ProvisionRequest or throws an error
// if the body cannot be decoded correctly
func (p *ProvisionHandler) decodeRequestBody(request *http.Request) (*ProvisionRequest, error) {
	decodedRequest := new(ProvisionRequest)

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&request)
	if err != nil {
		return nil, fmt.Errorf("encountered error while decoding request body: %w", err)
	}

	return decodedRequest, nil
}

// handleProvisionRepository processes the request of provisioning a repository and will generate the following status code:
//	- 201	If the repository, token and optionally a user have been created successfully
//	- 400 	If the request body can not be decoded
//  - 409	If the repository already exists on the Gitea server
//  - 503 	If the upstream Gitea repository is not available
func (p *ProvisionHandler) handleProvisionRepository(w http.ResponseWriter, req *http.Request) {
	request, err := p.decodeRequestBody(req)
	if err != nil {
		log.Printf("Unable to process request body: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := p.Provisioner.ProvisionRepository(request.Namespace, request.Project)
	if err != nil {
		if errors.Is(err, ErrRepositoryAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			return
		}

		log.Printf("Unable to create repository: %s\n", err.Error())
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		log.Printf("Unable to marshal reponse: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully provisioned repo %s in namespace %s\n", request.Project, request.Namespace)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(responseJson)
	if err != nil {
		log.Printf("Encountered error while writing response body: %s", err.Error())
	}
}

// handleDeleteRepository processes the request of deleting a repository and will generate the following status codes:
//   - 204  If the repository has been deleted successfully
//   - 400 	If the request body can not be decoded
//   - 404 	If the given repository cannot be found
//   - 503  If the upstream Gitea repository is not available
func (p *ProvisionHandler) handleDeleteRepository(w http.ResponseWriter, req *http.Request) {
	request, err := p.decodeRequestBody(req)
	if err != nil {
		log.Printf("Unable to process request body: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("Deleting repository %s in namspace %s\n", request.Project, request.Namespace)

	err = p.Provisioner.DeleteRepository(request.Namespace, request.Project)
	if err != nil {
		if errors.Is(err, ErrRepositoryDoesNotExist) {
			log.Printf("Unable to delete repository, does not exist!\n")
			w.WriteHeader(http.StatusNotFound)
			return
		}

		log.Printf("Unable to delete repository: %s\n", err.Error())
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

package keptn

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

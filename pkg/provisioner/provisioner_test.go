package provisioner

import (
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"keptn-sandbox/keptn-gitea-provisioner/pkg/keptn"
	"keptn-sandbox/keptn-gitea-provisioner/pkg/provisioner/fake"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProvisionHandler_CreateRepository(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	provisioner := fake.NewMockGitProvisioner(mockCtrl)
	handler := ProvisionHandler{
		Provisioner: provisioner,
	}

	request, _ := http.NewRequest(http.MethodPost, "/repository",
		strings.NewReader(`{"namespace":"keptn","project":"test"}`),
	)
	response := httptest.NewRecorder()

	provisioner.EXPECT().ProvisionRepository("keptn", "test").Times(1).Return(&keptn.ProvisionResponse{
		GitRemoteURL: "http://some.git.server:9999/user-keptn/repository-test",
		GitToken:     "8399p4q8cbunq983N489VNB2Q89T7B09",
		GitUser:      "user-keptn",
	}, nil)

	handler.HandleProvisionRepoRequest(response, request)
	require.Equal(t, http.StatusCreated, response.Code)

	expectedResult := map[string]string{
		"gitRemoteURL": "http://some.git.server:9999/user-keptn/repository-test",
		"gitToken":     "8399p4q8cbunq983N489VNB2Q89T7B09",
		"gitUser":      "user-keptn",
	}

	var responseBody map[string]string
	err := json.Unmarshal(response.Body.Bytes(), &responseBody)
	require.NoError(t, err)
	require.Equal(t, expectedResult, responseBody)
}

func TestProvisionHandler_CreateRepositoryError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	provisioner := fake.NewMockGitProvisioner(mockCtrl)
	handler := ProvisionHandler{
		Provisioner: provisioner,
	}

	request, _ := http.NewRequest(http.MethodPost, "/repository",
		strings.NewReader(`{"namespace":"keptn","project":"test"}`),
	)
	response := httptest.NewRecorder()

	provisioner.EXPECT().ProvisionRepository("keptn", "test").Times(1).Return(
		nil, ErrRepositoryAlreadyExists,
	)

	handler.HandleProvisionRepoRequest(response, request)
	require.Equal(t, http.StatusConflict, response.Code)
	require.Equal(t, response.Body.Len(), 0)
}

func TestProvisionHandler_DeleteRepository(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	provisioner := fake.NewMockGitProvisioner(mockCtrl)
	handler := ProvisionHandler{
		Provisioner: provisioner,
	}

	request, _ := http.NewRequest(http.MethodDelete, "/repository",
		strings.NewReader(`{"namespace":"keptn","project":"test"}`),
	)
	response := httptest.NewRecorder()

	provisioner.EXPECT().DeleteRepository("keptn", "test").Times(1).Return(nil)

	handler.HandleProvisionRepoRequest(response, request)
	require.Equal(t, http.StatusNoContent, response.Code)
	require.Equal(t, response.Body.Len(), 0)
}

func TestProvisionHandler_DeleteRepositoryError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	provisioner := fake.NewMockGitProvisioner(mockCtrl)
	handler := ProvisionHandler{
		Provisioner: provisioner,
	}

	request, _ := http.NewRequest(http.MethodDelete, "/repository",
		strings.NewReader(`{"namespace":"keptn","project":"test"}`),
	)
	response := httptest.NewRecorder()

	provisioner.EXPECT().DeleteRepository("keptn", "test").Times(1).Return(ErrRepositoryDoesNotExist)

	handler.HandleProvisionRepoRequest(response, request)
	require.Equal(t, http.StatusNotFound, response.Code)
	require.Equal(t, response.Body.Len(), 0)
}

func TestProvisionHandler_InvalidRequestBody(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	provisioner := fake.NewMockGitProvisioner(mockCtrl)
	provisioner.EXPECT().ProvisionRepository("", "").Times(1).Return(nil, ErrInvalidRequest)
	provisioner.EXPECT().DeleteRepository("", "").Times(1).Return(ErrInvalidRequest)

	handler := ProvisionHandler{
		provisioner,
	}

	tests := []struct {
		Test    string
		method  string
		content string
		code    int
	}{
		{
			Test:    "CreateRepository_EmptyJSON",
			method:  http.MethodPost,
			content: `{}`,
			code:    http.StatusUnprocessableEntity,
		},
		{
			Test:    "CreateRepository_InvalidBody",
			method:  http.MethodPost,
			content: `asiovsadifvbsapüoi`,
			code:    http.StatusBadRequest,
		},
		{
			Test:    "DeleteRepository_EmptyJSON",
			method:  http.MethodDelete,
			content: `{}`,
			code:    http.StatusUnprocessableEntity,
		},
		{
			Test:    "DeleteRepository_InvalidBody",
			method:  http.MethodDelete,
			content: `asiovsadifvbsapüoi`,
			code:    http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.Test, func(t *testing.T) {
			request, _ := http.NewRequest(test.method, "/repository",
				strings.NewReader(test.content),
			)
			response := httptest.NewRecorder()

			handler.HandleProvisionRepoRequest(response, request)
			assert.Equal(t, test.code, response.Code)
			assert.Equal(t, response.Body.Len(), 0)
		})
	}
}

func TestProvisionHandler_InvalidMethod(t *testing.T) {
	handler := ProvisionHandler{}
	request, _ := http.NewRequest(http.MethodGet, "/repository", nil)
	response := httptest.NewRecorder()

	handler.HandleProvisionRepoRequest(response, request)
	assert.Equal(t, http.StatusMethodNotAllowed, response.Code)
	assert.Equal(t, response.Body.Len(), 0)
}

func TestProvisionHandler_UnresponsiveUpstream(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	provisioner := fake.NewMockGitProvisioner(mockCtrl)
	provisioner.EXPECT().ProvisionRepository("keptn", "test").Times(1).Return(nil, fmt.Errorf("upstream error"))
	provisioner.EXPECT().DeleteRepository("keptn", "test").Times(1).Return(fmt.Errorf("upstream error"))

	handler := ProvisionHandler{
		provisioner,
	}

	tests := []struct {
		Test   string
		method string
	}{
		{
			Test:   "CreateRepository",
			method: http.MethodPost,
		},
		{
			Test:   "DeleteRepository",
			method: http.MethodDelete,
		},
	}

	for _, test := range tests {
		t.Run(test.Test, func(t *testing.T) {
			request, _ := http.NewRequest(test.method, "/repository",
				strings.NewReader(`{"namespace":"keptn","project":"test"}`),
			)
			response := httptest.NewRecorder()

			handler.HandleProvisionRepoRequest(response, request)
			assert.Equal(t, http.StatusFailedDependency, response.Code)
			assert.Equal(t, response.Body.Len(), 0)
		})
	}

}

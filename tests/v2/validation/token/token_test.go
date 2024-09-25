//go:build (validation || infra.any || cluster.any || sanity) && !stress

package token

import (
	"fmt"
	"github.com/rancher/rancher/tests/v2/actions/kubeapi/tokens"
	"github.com/rancher/shepherd/clients/rancher"
	fv3 "github.com/rancher/shepherd/clients/rancher/generated/management/v3"
	"github.com/rancher/shepherd/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

const (
	initialTokenDesc = "my-token"
	updatedTokenDesc = "changed-token"
	localClusterID   = "local"
)

type TokenTestSuite struct {
	suite.Suite
	client  *rancher.Client
	clients []*rancher.Client
	session *session.Session
}

func (t *TokenTestSuite) TearDownSuite() {
	t.session.Cleanup()
}

func (t *TokenTestSuite) SetupSuite() {
	testSession := session.NewSession()
	t.session = testSession

	client, err := rancher.NewClient("", t.session)
	require.NoError(t.T(), err)

	t.client = client

	// Initialize multiple Rancher clients
	t.clients, err = client.RancherClients()
	t.clients = append(t.clients, t.client)
	require.NoError(t.T(), err)
	require.NotEmpty(t.T(), t.clients)
}

func (t *TokenTestSuite) TestPatchToken() {
	for i, client := range t.clients {
		t.Run(fmt.Sprintf("RancherInstance_%d", i), func() {
			tokenToCreate := &fv3.Token{Description: initialTokenDesc}
			createdToken, err := client.Management.Token.Create(tokenToCreate)
			require.NoError(t.T(), err)

			assert.Equal(t.T(), initialTokenDesc, createdToken.Description)

			patchedToken, unstructuredRes, err := tokens.PatchToken(client, localClusterID, createdToken.Name, "replace", "/description", updatedTokenDesc)
			require.NoError(t.T(), err)

			assert.Equal(t.T(), updatedTokenDesc, patchedToken.Description)

			uc := unstructuredRes.UnstructuredContent()
			if val, ok := uc["groupPrincipals"]; ok {
				assert.NotEmpty(t.T(), val)
			}
		})
	}
}

func TestTokenTestSuite(t *testing.T) {
	suite.Run(t, new(TokenTestSuite))
}

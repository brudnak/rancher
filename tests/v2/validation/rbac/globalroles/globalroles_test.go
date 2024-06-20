package globalroles

import (
	"fmt"
	"github.com/rancher/shepherd/clients/rancher"
	"github.com/rancher/shepherd/pkg/session"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

const (
	branch              = "dev"
	url                 = "https://git.rancher.io/system-charts"
	newUserDefaultFalse = false
)

type GlobalRoleTestSuite struct {
	suite.Suite
	client  *rancher.Client
	session *session.Session
}

func (g *GlobalRoleTestSuite) SetupSuite() {
	testSession := session.NewSession()
	g.session = testSession

	client, err := rancher.NewClient("", testSession)
	require.NoError(g.T(), err)

	g.client = client
}

func (g *GlobalRoleTestSuite) TearDownSuite() {
	g.session.Cleanup()
}

func (g *GlobalRoleTestSuite) TestGlobalRoleCreate1() {
	// Create a new global role that permits creating catalogs
	dadFish, err := createGlobalRoleForCatalogs(g.client, newUserDefaultFalse)
	require.NoError(g.T(), err)
	fmt.Println(dadFish)

	//// Create a new user
	//user, token := createUser(g.client)
	//defer deleteUser(g.client, user)
	//
	//// Check that the user cannot create catalogs
	//name := randomName()
	//validateCreateCatalog(token, name, branch, url, false)
	//
	//// Assign the global role to the user
	//createGlobalRoleBinding(g.client, gr.Name, user.Name)
	//
	//// Check that the user has the global role
	//grbList, err := g.client.GlobalRoleBinding.List(metav1.ListOptions{
	//	FieldSelector: "globalRoleName=" + gr.Name + ",userName=" + user.Name,
	//})
	//require.NoError(g.T(), err)
	//require.Len(g.T(), grbList.Items, 1)
	//
	//// Check that the user can create catalogs
	//validateCreateCatalog(token, name, branch, url, true)
}

func TestGlobalRoleTestSuite(t *testing.T) {
	suite.Run(t, new(GlobalRoleTestSuite))
}

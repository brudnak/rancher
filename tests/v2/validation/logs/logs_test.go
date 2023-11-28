//go:build (validation || infra.any || cluster.any || sanity) && !stress && !extended

package logs

import (
	"context"
	"github.com/rancher/rancher/tests/framework/extensions/cloudcredentials"
	"github.com/rancher/rancher/tests/framework/pkg/config"
	"testing"

	v1 "github.com/rancher/rancher/pkg/apis/provisioning.cattle.io/v1"
	"github.com/rancher/rancher/tests/framework/clients/rancher"
	management "github.com/rancher/rancher/tests/framework/clients/rancher/generated/management/v3"
	"github.com/rancher/rancher/tests/framework/pkg/session"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	fleetLocal   = "fleet-local"
	localCluster = "local"
)

type LogsTestSuite struct {
	suite.Suite
	client       *rancher.Client
	session      *session.Session
	cluster      *management.Cluster
	clusterV1    *v1.Cluster
	awsEC2Config *cloudcredentials.AmazonEC2CredentialConfig
}

func (l *LogsTestSuite) TearDownSuite() {
	l.session.Cleanup()
}

func (l *LogsTestSuite) SetupSuite() {
	testSession := session.NewSession()
	l.session = testSession

	client, err := rancher.NewClient("", testSession)
	require.NoError(l.T(), err)

	ec2Config := new(cloudcredentials.AmazonEC2CredentialConfig)
	config.LoadConfig(cloudcredentials.AmazonEC2CredentialConfigurationFileKey, ec2Config)

	l.awsEC2Config = ec2Config
	l.client = client

	require.NoError(l.T(), err)

	kubeClient, err := l.client.GetKubeAPIProvisioningClient()
	require.NoError(l.T(), err)

	l.clusterV1, err = kubeClient.Clusters(fleetLocal).Get(context.TODO(), localCluster, metav1.GetOptions{})
	require.NoError(l.T(), err)
}

func (l *LogsTestSuite) auditLogs() {

	l.Run("EC2 Log Check", func() {
		checkLogs(l.T(), l.client, l.clusterV1, localCluster, l.awsEC2Config.AccessKey)
	})
}

func (l *LogsTestSuite) TestAuditLogs() {
	subSession := l.session.NewSession()
	defer subSession.Cleanup()
	l.auditLogs()
}

func TestLogsTestSuite(t *testing.T) {
	suite.Run(t, new(LogsTestSuite))
}

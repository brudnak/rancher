package logs

import (
	"github.com/rancher/rancher/tests/framework/extensions/cloudcredentials/aws"
	"github.com/stretchr/testify/assert"
	"log"
	"regexp"
	"strings"
	"testing"

	v1 "github.com/rancher/rancher/pkg/apis/provisioning.cattle.io/v1"
	"github.com/rancher/rancher/tests/framework/clients/rancher"
	"github.com/rancher/rancher/tests/framework/extensions/kubectl"
	"github.com/stretchr/testify/require"
)

func checkLogs(t *testing.T, client *rancher.Client, clusterV1 *v1.Cluster, clusterID, searchTerm string) {

	rancherPodCommand := []string{"kubectl", "get", "pods", "-n", "cattle-system"}
	rancherPods, err := kubectl.Command(client, clusterV1, nil, clusterID, rancherPodCommand)
	require.NoError(t, err)

	pattern := `rancher-\S+`

	r, err := regexp.Compile(pattern)
	if err != nil {
		log.Println("Error compiling regex:", err)
		return
	}

	matches := r.FindAllString(rancherPods, -1)

	var filteredMatches []string
	for _, match := range matches {
		if !containsWebhook(match) {
			filteredMatches = append(filteredMatches, match)
		}
	}

	_, err = aws.CreateAWSCloudCredentials(client)
	if err != nil {
		return
	}

	cmd := []string{"kubectl", "logs", "", "-n", "cattle-system", "--all-containers", "--tail=200"}
	for _, match := range filteredMatches {
		cmd[2] = match
		logs, err := kubectl.Command(client, clusterV1, nil, clusterID, cmd)
		require.NoError(t, err)
		assert.Equal(t, false, strings.Contains(logs, searchTerm))
	}
}

func containsWebhook(s string) bool {
	return regexp.MustCompile(`rancher-webhook`).MatchString(s)
}

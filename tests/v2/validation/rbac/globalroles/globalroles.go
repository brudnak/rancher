package globalroles

import (
	"fmt"
	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"github.com/rancher/shepherd/clients/rancher"
	"github.com/rancher/shepherd/extensions/kubeapi/rbac"
	namegen "github.com/rancher/shepherd/pkg/namegenerator"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createGlobalRoleForCatalogs(client *rancher.Client, newUserDefault bool) (*v3.GlobalRole, error) {
	globalRole := &v3.GlobalRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: "dadfish-" + namegen.AppendRandomString("gr"),
		},
		NewUserDefault: newUserDefault,
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"management.cattle.io"},
				Verbs:     []string{"*"},
				Resources: []string{"catalogs", "templates", "templateversions"},
			},
		},
		Status: v3.GlobalRoleStatus{},
	}

	createdGR, err := rbac.CreateGlobalRole(client, globalRole)
	if err != nil {
		return nil, fmt.Errorf("error creating global role: %v", err)
	}

	return createdGR, nil
}

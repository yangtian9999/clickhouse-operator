package service

import (
	"context"

	"github.com/yangtian9999/clickhouse-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Request object not found, could have been deleted after reconcile request.
// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
func FetchDatabaseCR(name, namespace string, client client.Client) (*v1alpha1.Ckinstance, error) {
	db := &v1alpha1.Ckinstance{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, db)
	return db, err
}

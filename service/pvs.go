package service

import (
	"github.com/yangtian9999/clickhouse-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

//Returns the deployment object for the Database
func NewDatabasePvc(db *v1alpha1.Ckinstance, scheme *runtime.Scheme) *corev1.PersistentVolumeClaim {
	ls := map[string]string{"owner": "ckoperator", "cr": db.Spec.DatabaseName}
	pv := &corev1.PersistentVolumeClaim{
		ObjectMeta: v1.ObjectMeta{
			Name:      db.Spec.DatabaseName,
			Namespace: db.Spec.Namespace,
			Labels:    ls,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(db.Spec.DatabaseStorageRequest),
				},
			},
			StorageClassName: &db.Spec.DatabaseStorageClassName,
		},
	}
	controllerutil.SetControllerReference(db, pv, scheme)
	return pv
}

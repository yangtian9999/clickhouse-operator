package service

import (
	"github.com/yangtian9999/clickhouse-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

//buildDBDeployment returns the deployment object for the Database
func NewDatabaseDeployment(db *v1alpha1.Ckinstance, scheme *runtime.Scheme) *appsv1.Deployment {
	ls := map[string]string{"owner": "ckoperator", "cr": db.Spec.DatabaseName}
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      db.Spec.DatabaseName,
			Namespace: db.Spec.Namespace,
			Labels:    ls,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: db.Spec.Shards,
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RecreateDeploymentStrategyType,
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           db.Spec.Image,
						Name:            db.Spec.DatabaseName,
						ImagePullPolicy: "IfNotPresent",
						Ports: []corev1.ContainerPort{{
							ContainerPort: db.Spec.TcpPort,
							Protocol:      "TCP",
						}},
						VolumeMounts: []corev1.VolumeMount{{
							Name:      db.Spec.DatabaseName,
							MountPath: "/etc/clickhouse-server/config.d/",
							// MountPath: "/etc/clickhouse-server/",
						}},
					}},
					Volumes: []corev1.Volume{{
						Name: db.Spec.DatabaseName,
						VolumeSource: corev1.VolumeSource{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
								ClaimName: db.Spec.DatabaseName,
							},
						},
					}},
				},
			},
		},
	}
	controllerutil.SetControllerReference(db, dep, scheme)
	return dep
}

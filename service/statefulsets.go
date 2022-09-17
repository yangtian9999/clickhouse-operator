package service

import (
	"fmt"

	"github.com/yangtian9999/clickhouse-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

//buildDBDeployment returns the deployment object for the Database
func NewStatefulsSet(db *v1alpha1.Ckinstance, scheme *runtime.Scheme) *appsv1.StatefulSet {
	ls := map[string]string{"owner": "ckoperator", "cr": db.Spec.DatabaseName}
	var terminationGracePeriodSeconds int64 = 30
	sta := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      db.Spec.DatabaseName,
			Namespace: db.Spec.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: fmt.Sprintf("%s-headless", db.Spec.DatabaseName),
			Replicas:    db.Spec.Shards,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
				Type: appsv1.RollingUpdateStatefulSetStrategyType,
			},
			// PodManagementPolicy: appsv1.OrderedReadyPodManagement,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					TerminationGracePeriodSeconds: &terminationGracePeriodSeconds,
					Containers: []corev1.Container{{
						Image:           db.Spec.Image,
						Name:            "clickhouse-server",
						ImagePullPolicy: "IfNotPresent",
						Ports: []corev1.ContainerPort{{
							ContainerPort: db.Spec.TcpPort,
							Name:          "tcpport",
						}},
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "clickhouse-data",
							MountPath: "/etc/clickhouse-server/config.d/",
							// MountPath: "/etc/clickhouse-server/",
						}},
					}},
					// Volumes: []corev1.Volume{{
					// 	Name: "click",
					// 	VolumeSource: corev1.VolumeSource{
					// 		PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					// 			ClaimName: "click",
					// 		},
					// 	},
					// }},
				},
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
				// 	ObjectMeta: metav1.ObjectMeta{
				// 		Name:      "clickhouse-data",
				// 		Namespace: db.Spec.Namespace,
				// 	},
				// 	Spec: corev1.PersistentVolumeClaimSpec{
				// 		AccessModes: []corev1.PersistentVolumeAccessMode{
				// 			corev1.ReadWriteOnce,
				// 		},
				// 		Resources: corev1.ResourceRequirements{
				// 			Requests: corev1.ResourceList{
				// 				corev1.ResourceStorage: resource.MustParse(db.Spec.DatabaseStorageRequest),
				// 			},
				// 		},
				// 	},
				// 	StorageClassName: &db.Spec.DatabaseStorageClassName,
			},
		},
	}
	controllerutil.SetControllerReference(db, sta, scheme)
	return sta
}

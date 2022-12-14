/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"encoding/json"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	ckopv1alpha1 "github.com/yangtian9999/clickhouse-operator/api/v1alpha1"
	"github.com/yangtian9999/clickhouse-operator/service"
)

// CkinstanceReconciler reconciles a Ckinstance object
type CkinstanceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=ckop.yt9999.io,resources=ckinstances,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ckop.yt9999.io,resources=ckinstances/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ckop.yt9999.io,resources=ckinstances/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=statefulsets/status,verbs=get
//+kubebuilder:rbac:groups=core,resources=events;pods;pods/exec;persistentvolumeclaims;persistentvolumes;configmaps;secrets;services,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Ckinstance object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *CkinstanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	instance := &ckopv1alpha1.Ckinstance{}
	if err := r.Client.Get(ctx, req.NamespacedName, instance); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		// return ctrl.Result{}, nil
	}

	// If ck cr was deleted, return nil
	if instance.DeletionTimestamp != nil {
		return ctrl.Result{}, nil
	}

	// Get ck deployment resource
	// oldDeployment := &appsv1.Deployment{}
	// oldCkInstance := &ckopv1alpha1.Ckinstance{}
	oldStatefulSet := &appsv1.StatefulSet{}

	//TODO: ??????StatefulSet?????????????????????
	if err := r.Client.Get(ctx, req.NamespacedName, oldStatefulSet); err != nil {
		if errors.IsNotFound(err) {
			// if err := r.Client.Create(context.TODO(), service.NewDatabasePvc(instance, r.Scheme)); err != nil {
			// 	return ctrl.Result{}, err
			// }

			if err := r.Client.Create(ctx, service.NewStatefulsSet(instance, r.Scheme)); err != nil {
				return ctrl.Result{}, err
			}

			data, _ := json.Marshal(instance.Spec)
			if instance.Annotations != nil {
				instance.Annotations["spec"] = string(data)
			} else {
				instance.Annotations = map[string]string{"spec": string(data)}
			}
			if err := r.Client.Update(ctx, instance); err != nil {
				return ctrl.Result{}, err
			}

		} else {
			return ctrl.Result{}, err
		}
		// } else {
		// 	//TODO: ??????StatefulSet??????????????????
		// 	// oldSpec := redisv1.RedisSpec{}
		// 	oldSpec := ckopv1alpha1.CkinstanceSpec{}
		// 	if err := json.Unmarshal([]byte(instance.Annotations["spec"]), &oldSpec); err != nil {
		// 		return ctrl.Result{}, nil
		// 	}

		// 	// ????????????????????????????????????????????? ???????????????????????????????????????
		// 	if !reflect.DeepEqual(oldSpec, instance.Spec) {
		// 		// ??????StatefulSet, ?????????Spec
		// 		newStatefulSet := statefulset.New(instance)
		// 		oldStatefulset.Spec = newStatefulSet.Spec
		// 		if err := r.Client.Update(ctx, oldStatefulset); err != nil {
		// 			return ctrl.Result{}, err
		// 		}

		// 		// ??????service
		// 		newService := service.New(instance)
		// 		oldService := &corev1.Service{}
		// 		if err := r.Client.Get(ctx, req.NamespacedName, oldService); err != nil {
		// 			return ctrl.Result{}, err
		// 		}
		// 		// ????????????Service???Spec???????????????????????????????????????????????????
		// 		clusterIP := oldService.Spec.ClusterIP
		// 		oldService.Spec = newService.Spec
		// 		oldService.Spec.ClusterIP = clusterIP // Service???ClusterIP, 10.254.x.x
		// 		if err := r.Client.Update(ctx, oldService); err != nil {
		// 			return ctrl.Result{}, err
		// 		}

		// 		// ??????????????? Annotations
		// 		data, _ := json.Marshal(instance.Spec)
		// 		if instance.Annotations != nil {
		// 			instance.Annotations["spec"] = string(data)
		// 		} else {
		// 			instance.Annotations = map[string]string{"spec": string(data)}
		// 		}
		// 		if err := r.Client.Update(ctx, instance); err != nil {
		// 			return ctrl.Result{}, err
		// 		}
		// 	}
	}

	return ctrl.Result{}, nil

}

// SetupWithManager sets up the controller with the Manager.
func (r *CkinstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ckopv1alpha1.Ckinstance{}).
		Complete(r)
}

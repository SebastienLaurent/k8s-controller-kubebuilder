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
	"reflect"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	ref "k8s.io/client-go/tools/reference"

	samplev1 "github.com/SebastienLaurent/k8s-controller-kubebuilder/api/v1"
)

// ModuleReconciler reconciles a Module object
type ModuleReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=sample.alien4cloud,resources=modules,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=sample.alien4cloud,resources=modules/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=sample.alien4cloud,resources=modules/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Module object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *ModuleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var module samplev1.Module
	if err := r.Get(ctx, req.NamespacedName, &module); err != nil {
		if errors.IsNotFound(err) {
			log.Info("Deleting module")
			return ctrl.Result{}, nil
		}
		log.Error(err, "unable to fetch Pod")
		return ctrl.Result{}, err
	}

	log.V(1).Info("Reconcile Module:  ", "request", req)

	if module.Status.Sidecar == nil {
		log.Info("Creating Sidecar Pod")

		sidePod, err := r.launchSidecar(ctx, module, req.Namespace)
		if err != nil {
			return ctrl.Result{}, err
		}

		if err := r.updateReference(ctx, module, sidePod); err != nil {
			return ctrl.Result{}, err
		}

	} else {
		sidePod := &corev1.Pod{}
		err := r.Get(ctx, client.ObjectKey{Namespace: module.Status.Sidecar.Namespace, Name: module.Status.Sidecar.Name}, sidePod)
		if err != nil && !errors.IsNotFound(err) {
			log.Error(err, "unable to get sidecar pod")
			return ctrl.Result{}, err
		}
		if err != nil {
			log.Info("Recreate Sidecar Pod")

			sidePod, err := r.launchSidecar(ctx, module, req.Namespace)
			if err != nil {
				return ctrl.Result{}, err
			}

			if err := r.updateReference(ctx, module, sidePod); err != nil {
				return ctrl.Result{}, err
			}

		} else {
			if reflect.DeepEqual(sidePod.Spec.Containers[0], module.Spec.Sidecar) {
				log.Info("Spec is equal")
			} else {
				log.Info("Must recheck the pod")
			}
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ModuleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&samplev1.Module{}).
		Owns(&corev1.Pod{}).
		Complete(r)
}

func (r *ModuleReconciler) updateReference(ctx context.Context, module samplev1.Module, pod *corev1.Pod) error {
	log := log.FromContext(ctx)

	podRef, err := ref.GetReference(r.Scheme, pod)
	if err != nil {
		log.Error(err, "Cant get sidecar reference")
		return err
	}
	module.Status.Sidecar = podRef

	if err := r.Status().Update(ctx, &module); err != nil {
		log.Error(err, "unable to update module status")
		return err
	}
	return nil
}

func (r *ModuleReconciler) launchSidecar(ctx context.Context, module samplev1.Module, namespace string) (*corev1.Pod, error) {
	sidePod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Labels:      make(map[string]string),
			Annotations: make(map[string]string),
			Name:        "sidecar",
			Namespace:   namespace,
		},
		Spec: *module.Spec.Sidecar.DeepCopy(),
	}

	if err := ctrl.SetControllerReference(&module, sidePod, r.Scheme); err != nil {
		return nil, err
	}

	if err := r.Create(ctx, sidePod); err != nil {
		return nil, err
	}

	return sidePod, nil
}

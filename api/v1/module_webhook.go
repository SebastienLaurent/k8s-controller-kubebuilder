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

package v1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var modulelog = logf.Log.WithName("module-resource")

func (r *Module) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-sample-alien4cloud-v1-module,mutating=true,failurePolicy=fail,sideEffects=None,groups=sample.alien4cloud,resources=modules,verbs=create;update,versions=v1,name=mmodule.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Module{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Module) Default() {
	modulelog.Info("default", "name", r.Name)

	if r.Spec.Cu == "" {
		r.Spec.Cu = "test"
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-sample-alien4cloud-v1-module,mutating=false,failurePolicy=fail,sideEffects=None,groups=sample.alien4cloud,resources=modules,verbs=create;update,versions=v1,name=vmodule.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Module{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Module) ValidateCreate() error {
	modulelog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Module) ValidateUpdate(old runtime.Object) error {
	modulelog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Module) ValidateDelete() error {
	modulelog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
/*

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

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var testlog = logf.Log.WithName("test-resource")

func (r *Test) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-stressed-jpang-dev-v1alpha1-test,mutating=true,failurePolicy=fail,groups=stressed.jpang.dev,resources=tests,verbs=create;update,versions=v1alpha1,name=mtest.kb.io

var _ webhook.Defaulter = &Test{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Test) Default() {
	testlog.Info("default", "name", r.Name)

	if r.Spec.Replicas == nil {
		r.Spec.Replicas = new(int32)
		*r.Spec.Replicas = 1
	}

	if r.Spec.Image == "" {
		r.Spec.Image = "jpangms/stress-ng:latest"
	}

	if r.Spec.JobFile == "" {
		r.Spec.JobFile = defaultJobFile
	}
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-stressed-jpang-dev-v1alpha1-test,mutating=false,failurePolicy=fail,groups=stressed.jpang.dev,resources=tests,versions=v1alpha1,name=vtest.kb.io

var _ webhook.Validator = &Test{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Test) ValidateCreate() error {
	testlog.Info("validate create", "name", r.Name)

	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Test) ValidateUpdate(old runtime.Object) error {
	testlog.Info("validate update", "name", r.Name)

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Test) ValidateDelete() error {
	testlog.Info("validate delete", "name", r.Name)

	return nil
}

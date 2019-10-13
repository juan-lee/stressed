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

package controllers

import (
	"context"
	"crypto/sha256"
	"fmt"
	"reflect"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	stressedv1alpha1 "github.com/juan-lee/stressed/api/v1alpha1"
)

// TestReconciler reconciles a Test object
type TestReconciler struct {
	client.Client
	Log logr.Logger
}

// +kubebuilder:rbac:groups=stressed.jpang.dev,resources=tests,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stressed.jpang.dev,resources=tests/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete

func (r *TestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stressedv1alpha1.Test{}).
		Complete(r)
}

func (r *TestReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("test", req.NamespacedName)

	instance := &stressedv1alpha1.Test{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	cm, err := r.reconcileConfigMap(ctx, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	privileged := true
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-deployment",
			Namespace: instance.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: stressedv1alpha1.GroupVersion.String(),
					Kind:       "Test",
					Name:       instance.Name,
					UID:        instance.UID,
				},
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: instance.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"deployment": instance.Name + "-deployment"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"deployment": instance.Name + "-deployment"}},
				Spec: corev1.PodSpec{
					NodeSelector: instance.Spec.NodeSelector,
					Containers: []corev1.Container{
						{
							SecurityContext: &corev1.SecurityContext{
								Privileged: &privileged,
							},
							Name:  instance.Name,
							Image: instance.Spec.Image,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "jobfile",
									ReadOnly:  true,
									MountPath: "/stress/jobs",
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "CONFIG_HASH",
									Value: fmt.Sprintf("%x", sha256.Sum256([]byte(cm.Data["jobfile"]))),
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "jobfile",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: instance.Name + "-config",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	log.Info("Reconciling Deployment", "deploy", deploy)

	found := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: deploy.Name, Namespace: deploy.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating Deployment", "Namespace", deploy.Namespace, "Name", deploy.Name)
		err = r.Create(ctx, deploy)
		if err != nil {
			return ctrl.Result{}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	}

	if !reflect.DeepEqual(deploy.Spec, found.Spec) {
		found.Spec = deploy.Spec
		log.Info("Updating Deployment", "Namespace", deploy.Namespace, "Name", deploy.Name)
		err = r.Update(ctx, found)
		if err != nil {
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

func (r *TestReconciler) reconcileConfigMap(ctx context.Context, instance *stressedv1alpha1.Test) (*corev1.ConfigMap, error) {
	log := r.Log.WithValues("test", fmt.Sprintf("%s/%s", instance.Namespace, instance.Name))

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-config",
			Namespace: instance.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: stressedv1alpha1.GroupVersion.String(),
					Kind:       "Test",
					Name:       instance.Name,
					UID:        instance.UID,
				},
			},
		},
		Data: map[string]string{
			"jobfile": instance.Spec.JobFile,
		},
	}
	log.Info("Reconciling ConfigMap", "cm", cm)

	found := &corev1.ConfigMap{}
	err := r.Get(ctx, types.NamespacedName{Name: cm.Name, Namespace: cm.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating ConfigMap", "Namespace", cm.Namespace, "Name", cm.Name)
		err = r.Create(ctx, cm)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	if !reflect.DeepEqual(cm.Data, found.Data) {
		found.Data = cm.Data
		log.Info("Updating ConfigMap", "Namespace", cm.Namespace, "Name", cm.Name)
		err = r.Update(ctx, found)
		if err != nil {
			return nil, err
		}
		return found, nil
	}
	return cm, nil
}

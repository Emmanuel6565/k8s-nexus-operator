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
	"fmt"

	cachev1 "github.com/Emmanuel6565/k8s-nexus-operator/api/v1"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NexusReconciler reconciles a Nexus object
type NexusReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=cache.foundry.io,resources=nexus,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cache.foundry.io,resources=nexus/status,verbs=get;update;patch

func (r *NexusReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("nexus", req.NamespacedName)

	// your logic here
	nexus := &cachev1.Nexus{}
	newDeployment := r.NexusDeployment(nexus)
	err := r.Create(ctx,newDeployement)
	fmt.Println(err)
	return ctrl.Result{}, nil
}

// retrun Deployment for Nexus instance
func (r *NexusReconciler) NexusDeployment(nexus *cachev1.Nexus) *appsv1.Deployment {
	labeltemplate := "nexus"
	nexusfullname := nexus.Spec.Name
	namespace := nexus.Spec.Namespace
	nexusdeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      &nexusfullname,
			Namespace: &namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &nexus.Spec.Replicascount,
			Selector: &metav1.LabelSelector{
				MatchLabels: labeltemplate,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labeltemplate,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: "ff",
							Image: &nexus.Spec.Imagename,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name: "nexus-data",
									MountPath: "/nexus-data",
								},
								{
									Name: "nexus-backup",
									MountPath: "/nexus-data/backup",
								},	
							},
						},
					}
				},
			},
		},
	}
}

// return Nexus service Nodeport
func (r *NexusReconciler) NexusService(nexus *cachev1.Nexus) *corev1.Service {

}

// return nexus configmap
func (r *NexusReconciler) NexusConfigmap(nexus *cachev1.Nexus) *corev1.ConfigMap {

}

// return a deployment for nexus backup
func (r *NexusReconciler) Nexusdeploymentbackup(nexus *cachev1.Nexus) *appsv1.Deployment {

}

// create service account for nexus
func (r *NexusReconciler) Serviceaccount() *corev1.ServiceAccount {

}

func (r *NexusReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1.Nexus{}).
		Complete(r)
}

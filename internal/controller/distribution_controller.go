/*
Copyright 2025.

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

package controller

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	secdistv1 "github.com/etzba/secret-distributor/api/v1"
	"github.com/etzba/secret-distributor/pkg/logger"
)

// DistributionReconciler reconciles a Distribution object
type DistributionReconciler struct {
	Logger *logger.Log
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=secdist.etzba.com,resources=distributions,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=secdist.etzba.com,resources=distributions/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=secdist.etzba.com,resources=distributions/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Distribution object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *DistributionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = logf.FromContext(ctx)

	log := logger.New()
	r.Logger = log

	r.Logger.Info(fmt.Sprintf("Found change in namespace %s. Start reconcile", req.Namespace))

	var resource secdistv1.Distribution

	if err := r.Get(context.Background(), req.NamespacedName, &resource); err != nil {
		if errors.IsNotFound(err) {
			r.Logger.Info("distribution resource is not found. skipping..")
			return ctrl.Result{Requeue: false, RequeueAfter: 0}, nil
		}
		r.Logger.Error("could not fetch resource"+resource.Kind, nil)
		return ctrl.Result{Requeue: true, RequeueAfter: time.Minute}, err
	}

	secret, err := r.readFromSecretAndCreateNewSecret(resource.Spec.SecretName, req.Namespace)
	if err != nil {
		r.Logger.Error("Failed to get secret by name", err)
		return ctrl.Result{Requeue: false}, nil
	}

	if err := r.Get(ctx, req.NamespacedName, secret); err != nil {
		if errors.IsNotFound(err) {
			r.Logger.Info("Not found. Creating secret")
			if err := r.Create(ctx, secret); err != nil {
				if errors.IsAlreadyExists(err) {
					if err := r.Update(ctx, secret); err != nil {
						if errors.IsInvalid(err) {
							r.Logger.Error("Invalid update", err)
						} else {
							r.Logger.Error("Unable to update secret", err)
						}
						return ctrl.Result{}, nil
					}
				}
				r.Logger.Error("Could not create secret", err)
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}

		if err := r.Update(ctx, secret); err != nil {
			if errors.IsInvalid(err) {
				r.Logger.Error("Invalid update", err)
			} else {
				r.Logger.Error("Unable to update secret", err)
			}
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DistributionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&secdistv1.Distribution{}).
		Named("distribution").
		Complete(r)
}

package configmap

/*
Copyright 2020 Crunchy Data Solutions, Inc.
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

import (
	"context"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/crunchydata/postgres-operator/internal/controller/filter"
	"github.com/crunchydata/postgres-operator/internal/kubeapi"
	"github.com/crunchydata/postgres-operator/internal/operator"
)

// ReconcileConfigMap holds resources for the ConfigMap reconciler
type ReconcileConfigMap struct {
	Client    client.Client
	ClientSet *kubeapi.Client
}

// Reconcile reconciles a ConfigMap in a namespace managed by the PostgreSQL Operator
func (r *ReconcileConfigMap) Reconcile(context context.Context,
	request reconcile.Request) (reconcile.Result, error) {

	// Run handleConfigMapSync, passing it the namespace and name of the configMap that
	// needs to be synced
	if err := r.reconcileConfigMap(request.Namespace, request.Name); err != nil {
		log.Errorf("ConfigMap Reconciler: error syncing ConfigMap '%s', will now requeue: %v",
			request.Name, err)
		return reconcile.Result{
			Requeue: true,
		}, err
	}

	return reconcile.Result{}, nil
}

// SetupWithManager adds the ConfigMap controller to the provided runtime manager.
func (r *ReconcileConfigMap) SetupWithManager(mgr manager.Manager) error {

	if err := builder.ControllerManagedBy(mgr).
		For(&corev1.ConfigMap{}).
		WithEventFilter(filter.GetManagedNamespaceFilter(context.TODO(), mgr.GetClient())).
		WithEventFilter(filter.GetPGHAConfigMapFilter()).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: *operator.Pgo.Pgo.ConfigMapWorkerCount,
		}).
		Complete(r); err != nil {
		return err
	}

	return nil
}

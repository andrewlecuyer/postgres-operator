package runtime

/*
Copyright 2020 Crunchy Data
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
	"time"

	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/crunchydata/postgres-operator/internal/controller/configmap"
	"github.com/crunchydata/postgres-operator/internal/kubeapi"
	"github.com/crunchydata/postgres-operator/internal/operator"
	crunchydatascheme "github.com/crunchydata/postgres-operator/pkg/generated/clientset/versioned/scheme"
)

// NewPGORuntimeManager creates a new controller runtime manager for the PostgreSQL Operator.  The
// manager returned is configured specifically for the PostgreSQL Operator, and includes any
// controllers that it will be responsible for managing.
func NewPGORuntimeManager() (manager.Manager, error) {

	// add pgo custom resources to the default scheme
	if err := crunchydatascheme.AddToScheme(scheme.Scheme); err != nil {
		return nil, err
	}

	// create controller runtime manager
	syncPeriod := time.Duration(*operator.Pgo.Pgo.ControllerGroupRefreshInterval) * time.Second
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{
		SyncPeriod: &syncPeriod,
		Scheme:     scheme.Scheme,
	})
	if err != nil {
		return nil, err
	}

	// create a client for kube resources
	clientSet, err := kubeapi.NewClient()
	if err != nil {
		return nil, err
	}

	// add all PostgreSQL Operator controllers to the runtime manager
	if err := addControllersToManager(mgr, clientSet); err != nil {
		return nil, err
	}

	return mgr, nil
}

// addControllersToManager adds all PostgreSQL Operator controllers to the provided controller
// runtime manager.
func addControllersToManager(mgr manager.Manager, clientSet *kubeapi.Client) error {
	r := &configmap.ReconcileConfigMap{
		Client:    mgr.GetClient(),
		ClientSet: clientSet,
	}
	if err := r.SetupWithManager(mgr); err != nil {
		return err
	}
	return nil
}

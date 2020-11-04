package filter

/*
Copyright 2017 - 2020 Crunchy Data
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

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	crtclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/crunchydata/postgres-operator/internal/config"
	"github.com/crunchydata/postgres-operator/internal/operator"
)

var (
	pghaConfigMapFilter = predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			return isPGHAConfigMap(e.Object)
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			return isPGHAConfigMap(e.ObjectNew)
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return isPGHAConfigMap(e.Object)
		},
	}
)

func GetPGHAConfigMapFilter() predicate.Funcs {
	return pghaConfigMapFilter
}

func GetManagedNamespaceFilter(context context.Context,
	client crtclient.Client) predicate.Funcs {
	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			return isManagedNamespace(context, client, e.Object.GetNamespace())
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			return isManagedNamespace(context, client, e.ObjectNew.GetNamespace())
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return isManagedNamespace(context, client, e.Object.GetNamespace())
		},
	}
}

func isManagedNamespace(context context.Context, client crtclient.Client,
	namespace string) bool {

	ns := &corev1.Namespace{}
	client.Get(context, crtclient.ObjectKey{
		Name: namespace,
	}, ns)

	labels := ns.GetLabels()
	if labels != nil && labels[config.LABEL_VENDOR] == config.LABEL_CRUNCHY &&
		labels[config.LABEL_PGO_INSTALLATION_NAME] == operator.InstallationName {
		return true
	}
	return false
}

func isPGHAConfigMap(o client.Object) bool {
	_, ok := o.GetLabels()[config.LABEL_PGHA_CONFIGMAP]
	return ok
}

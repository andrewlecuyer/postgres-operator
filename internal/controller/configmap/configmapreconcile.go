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
	"sync"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crunchydata/postgres-operator/internal/config"
	cfg "github.com/crunchydata/postgres-operator/internal/operator/config"
	crv1 "github.com/crunchydata/postgres-operator/pkg/apis/crunchydata.com/v1"
)

// handleConfigMapSync is responsible for syncing a configMap resource that has obtained from
// the ConfigMap controller's worker queue
func (r *ReconcileConfigMap) reconcileConfigMap(namespace, configMapName string) error {

	log.Debugf("ConfigMap Reconciler: handling a configmap sync for configMap %s in namespace "+
		"(namespace %s)", namespace, configMapName)

	configMap := &corev1.ConfigMap{}
	if err := r.Client.Get(context.TODO(), runtimeclient.ObjectKey{
		Namespace: namespace,
		Name:      configMapName,
	}, configMap); err != nil {
		return err
	}

	clusterName := configMap.GetObjectMeta().GetLabels()[config.LABEL_PG_CLUSTER]

	cluster := &crv1.Pgcluster{}
	err := r.Client.Get(context.TODO(), runtimeclient.ObjectKey{
		Namespace: namespace,
		Name:      clusterName,
	}, cluster)
	if err != nil {
		// If the pgcluster is not found, then simply log an error and return.  This should not
		// typically happen, but in the event of an orphaned configMap with no pgcluster we do
		// not want to keep re-queueing the same item.  If any other error is encountered then
		// return that error.
		if kerrors.IsNotFound(err) {
			log.Errorf("ConfigMap Reconciler: cannot find pgcluster for configMap %s (namespace %s),"+
				"ignoring", configMapName, namespace)
			return nil
		}
		return err
	}

	// if an upgrade is pending for the cluster, then don't attempt to sync and just return
	if cluster.GetAnnotations()[config.ANNOTATION_IS_UPGRADED] == config.ANNOTATIONS_FALSE {
		log.Debugf("ConfigMap Reconciler: syncing of configMap %s (namespace %s) disabled pending the "+
			"upgrade of cluster %s", configMapName, namespace, clusterName)
		return nil
	}

	// disable syncing when the cluster isn't currently initialized
	if cluster.Status.State != crv1.PgclusterStateInitialized {
		return nil
	}

	r.syncPGHAConfig(r.createPGHAConfigs(configMap, clusterName,
		cluster.GetObjectMeta().GetLabels()[config.LABEL_PGHA_SCOPE]))

	return nil
}

// createConfigurerMap creates the configs needed to sync the PGHA configMap
func (r *ReconcileConfigMap) createPGHAConfigs(configMap *corev1.ConfigMap,
	clusterName, clusterScope string) []cfg.Syncer {

	var configSyncers []cfg.Syncer

	configSyncers = append(configSyncers, cfg.NewDCS(configMap, r.ClientSet, clusterScope))

	localDBConfig, err := cfg.NewLocalDB(configMap, r.ClientSet.Config, r.ClientSet)
	// Just log the error and don't add to the map so a sync can still be attempted with
	// any other configurers
	if err != nil {
		log.Error(err)
	} else {
		configSyncers = append(configSyncers, localDBConfig)
	}

	return configSyncers
}

// syncAllConfigs takes a map of configurers and runs their sync functions concurrently
func (r *ReconcileConfigMap) syncPGHAConfig(configSyncers []cfg.Syncer) {

	var wg sync.WaitGroup

	for _, configSyncer := range configSyncers {

		wg.Add(1)

		go func(syncer cfg.Syncer) {
			if err := syncer.Sync(); err != nil {
				log.Error(err)
			}
			wg.Done()
		}(configSyncer)
	}

	wg.Wait()
}

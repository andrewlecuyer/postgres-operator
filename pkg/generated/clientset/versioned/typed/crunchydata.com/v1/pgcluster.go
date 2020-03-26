/*
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"time"

	v1 "github.com/crunchydata/postgres-operator/apis/crunchydata.com/v1"
	scheme "github.com/crunchydata/postgres-operator/pkg/generated/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// PgclustersGetter has a method to return a PgclusterInterface.
// A group's client should implement this interface.
type PgclustersGetter interface {
	Pgclusters(namespace string) PgclusterInterface
}

// PgclusterInterface has methods to work with Pgcluster resources.
type PgclusterInterface interface {
	Create(*v1.Pgcluster) (*v1.Pgcluster, error)
	Update(*v1.Pgcluster) (*v1.Pgcluster, error)
	UpdateStatus(*v1.Pgcluster) (*v1.Pgcluster, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error
	Get(name string, options metav1.GetOptions) (*v1.Pgcluster, error)
	List(opts metav1.ListOptions) (*v1.PgclusterList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Pgcluster, err error)
	PgclusterExpansion
}

// pgclusters implements PgclusterInterface
type pgclusters struct {
	client rest.Interface
	ns     string
}

// newPgclusters returns a Pgclusters
func newPgclusters(c *CrunchydataV1Client, namespace string) *pgclusters {
	return &pgclusters{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the pgcluster, and returns the corresponding pgcluster object, and an error if there is any.
func (c *pgclusters) Get(name string, options metav1.GetOptions) (result *v1.Pgcluster, err error) {
	result = &v1.Pgcluster{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("pgclusters").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Pgclusters that match those selectors.
func (c *pgclusters) List(opts metav1.ListOptions) (result *v1.PgclusterList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.PgclusterList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("pgclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested pgclusters.
func (c *pgclusters) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("pgclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a pgcluster and creates it.  Returns the server's representation of the pgcluster, and an error, if there is any.
func (c *pgclusters) Create(pgcluster *v1.Pgcluster) (result *v1.Pgcluster, err error) {
	result = &v1.Pgcluster{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("pgclusters").
		Body(pgcluster).
		Do().
		Into(result)
	return
}

// Update takes the representation of a pgcluster and updates it. Returns the server's representation of the pgcluster, and an error, if there is any.
func (c *pgclusters) Update(pgcluster *v1.Pgcluster) (result *v1.Pgcluster, err error) {
	result = &v1.Pgcluster{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("pgclusters").
		Name(pgcluster.Name).
		Body(pgcluster).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *pgclusters) UpdateStatus(pgcluster *v1.Pgcluster) (result *v1.Pgcluster, err error) {
	result = &v1.Pgcluster{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("pgclusters").
		Name(pgcluster.Name).
		SubResource("status").
		Body(pgcluster).
		Do().
		Into(result)
	return
}

// Delete takes name of the pgcluster and deletes it. Returns an error if one occurs.
func (c *pgclusters) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("pgclusters").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *pgclusters) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("pgclusters").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched pgcluster.
func (c *pgclusters) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Pgcluster, err error) {
	result = &v1.Pgcluster{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("pgclusters").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}

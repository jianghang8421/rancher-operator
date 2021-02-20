/*
Copyright 2021 Rancher Labs, Inc.

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

// Code generated by main. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	"github.com/rancher/lasso/pkg/client"
	"github.com/rancher/lasso/pkg/controller"
	v1 "github.com/rancher/rancher-operator/pkg/apis/rke.cattle.io/v1"
	"github.com/rancher/wrangler/pkg/generic"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type RKEBootstrapTemplateHandler func(string, *v1.RKEBootstrapTemplate) (*v1.RKEBootstrapTemplate, error)

type RKEBootstrapTemplateController interface {
	generic.ControllerMeta
	RKEBootstrapTemplateClient

	OnChange(ctx context.Context, name string, sync RKEBootstrapTemplateHandler)
	OnRemove(ctx context.Context, name string, sync RKEBootstrapTemplateHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() RKEBootstrapTemplateCache
}

type RKEBootstrapTemplateClient interface {
	Create(*v1.RKEBootstrapTemplate) (*v1.RKEBootstrapTemplate, error)
	Update(*v1.RKEBootstrapTemplate) (*v1.RKEBootstrapTemplate, error)

	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1.RKEBootstrapTemplate, error)
	List(namespace string, opts metav1.ListOptions) (*v1.RKEBootstrapTemplateList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.RKEBootstrapTemplate, err error)
}

type RKEBootstrapTemplateCache interface {
	Get(namespace, name string) (*v1.RKEBootstrapTemplate, error)
	List(namespace string, selector labels.Selector) ([]*v1.RKEBootstrapTemplate, error)

	AddIndexer(indexName string, indexer RKEBootstrapTemplateIndexer)
	GetByIndex(indexName, key string) ([]*v1.RKEBootstrapTemplate, error)
}

type RKEBootstrapTemplateIndexer func(obj *v1.RKEBootstrapTemplate) ([]string, error)

type rKEBootstrapTemplateController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewRKEBootstrapTemplateController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) RKEBootstrapTemplateController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &rKEBootstrapTemplateController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromRKEBootstrapTemplateHandlerToHandler(sync RKEBootstrapTemplateHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1.RKEBootstrapTemplate
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1.RKEBootstrapTemplate))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *rKEBootstrapTemplateController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1.RKEBootstrapTemplate))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateRKEBootstrapTemplateDeepCopyOnChange(client RKEBootstrapTemplateClient, obj *v1.RKEBootstrapTemplate, handler func(obj *v1.RKEBootstrapTemplate) (*v1.RKEBootstrapTemplate, error)) (*v1.RKEBootstrapTemplate, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *rKEBootstrapTemplateController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *rKEBootstrapTemplateController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *rKEBootstrapTemplateController) OnChange(ctx context.Context, name string, sync RKEBootstrapTemplateHandler) {
	c.AddGenericHandler(ctx, name, FromRKEBootstrapTemplateHandlerToHandler(sync))
}

func (c *rKEBootstrapTemplateController) OnRemove(ctx context.Context, name string, sync RKEBootstrapTemplateHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromRKEBootstrapTemplateHandlerToHandler(sync)))
}

func (c *rKEBootstrapTemplateController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *rKEBootstrapTemplateController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *rKEBootstrapTemplateController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *rKEBootstrapTemplateController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *rKEBootstrapTemplateController) Cache() RKEBootstrapTemplateCache {
	return &rKEBootstrapTemplateCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *rKEBootstrapTemplateController) Create(obj *v1.RKEBootstrapTemplate) (*v1.RKEBootstrapTemplate, error) {
	result := &v1.RKEBootstrapTemplate{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *rKEBootstrapTemplateController) Update(obj *v1.RKEBootstrapTemplate) (*v1.RKEBootstrapTemplate, error) {
	result := &v1.RKEBootstrapTemplate{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *rKEBootstrapTemplateController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *rKEBootstrapTemplateController) Get(namespace, name string, options metav1.GetOptions) (*v1.RKEBootstrapTemplate, error) {
	result := &v1.RKEBootstrapTemplate{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *rKEBootstrapTemplateController) List(namespace string, opts metav1.ListOptions) (*v1.RKEBootstrapTemplateList, error) {
	result := &v1.RKEBootstrapTemplateList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *rKEBootstrapTemplateController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *rKEBootstrapTemplateController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v1.RKEBootstrapTemplate, error) {
	result := &v1.RKEBootstrapTemplate{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type rKEBootstrapTemplateCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *rKEBootstrapTemplateCache) Get(namespace, name string) (*v1.RKEBootstrapTemplate, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v1.RKEBootstrapTemplate), nil
}

func (c *rKEBootstrapTemplateCache) List(namespace string, selector labels.Selector) (ret []*v1.RKEBootstrapTemplate, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.RKEBootstrapTemplate))
	})

	return ret, err
}

func (c *rKEBootstrapTemplateCache) AddIndexer(indexName string, indexer RKEBootstrapTemplateIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1.RKEBootstrapTemplate))
		},
	}))
}

func (c *rKEBootstrapTemplateCache) GetByIndex(indexName, key string) (result []*v1.RKEBootstrapTemplate, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v1.RKEBootstrapTemplate, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v1.RKEBootstrapTemplate))
	}
	return result, nil
}

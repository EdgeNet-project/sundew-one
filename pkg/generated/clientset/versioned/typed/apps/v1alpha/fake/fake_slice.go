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

package fake

import (
	"context"

	v1alpha "github.com/EdgeNet-project/edgenet/pkg/apis/apps/v1alpha"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeSlices implements SliceInterface
type FakeSlices struct {
	Fake *FakeAppsV1alpha
	ns   string
}

var slicesResource = schema.GroupVersionResource{Group: "apps.edgenet.io", Version: "v1alpha", Resource: "slices"}

var slicesKind = schema.GroupVersionKind{Group: "apps.edgenet.io", Version: "v1alpha", Kind: "Slice"}

// Get takes name of the slice, and returns the corresponding slice object, and an error if there is any.
func (c *FakeSlices) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha.Slice, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(slicesResource, c.ns, name), &v1alpha.Slice{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha.Slice), err
}

// List takes label and field selectors, and returns the list of Slices that match those selectors.
func (c *FakeSlices) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha.SliceList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(slicesResource, slicesKind, c.ns, opts), &v1alpha.SliceList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha.SliceList{ListMeta: obj.(*v1alpha.SliceList).ListMeta}
	for _, item := range obj.(*v1alpha.SliceList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested slices.
func (c *FakeSlices) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(slicesResource, c.ns, opts))

}

// Create takes the representation of a slice and creates it.  Returns the server's representation of the slice, and an error, if there is any.
func (c *FakeSlices) Create(ctx context.Context, slice *v1alpha.Slice, opts v1.CreateOptions) (result *v1alpha.Slice, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(slicesResource, c.ns, slice), &v1alpha.Slice{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha.Slice), err
}

// Update takes the representation of a slice and updates it. Returns the server's representation of the slice, and an error, if there is any.
func (c *FakeSlices) Update(ctx context.Context, slice *v1alpha.Slice, opts v1.UpdateOptions) (result *v1alpha.Slice, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(slicesResource, c.ns, slice), &v1alpha.Slice{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha.Slice), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeSlices) UpdateStatus(ctx context.Context, slice *v1alpha.Slice, opts v1.UpdateOptions) (*v1alpha.Slice, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(slicesResource, "status", c.ns, slice), &v1alpha.Slice{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha.Slice), err
}

// Delete takes name of the slice and deletes it. Returns an error if one occurs.
func (c *FakeSlices) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(slicesResource, c.ns, name), &v1alpha.Slice{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeSlices) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(slicesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha.SliceList{})
	return err
}

// Patch applies the patch and returns the patched slice.
func (c *FakeSlices) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha.Slice, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(slicesResource, c.ns, name, pt, data, subresources...), &v1alpha.Slice{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha.Slice), err
}

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

// FakeNodeContributions implements NodeContributionInterface
type FakeNodeContributions struct {
	Fake *FakeAppsV1alpha
	ns   string
}

var nodecontributionsResource = schema.GroupVersionResource{Group: "apps.edgenet.io", Version: "v1alpha", Resource: "nodecontributions"}

var nodecontributionsKind = schema.GroupVersionKind{Group: "apps.edgenet.io", Version: "v1alpha", Kind: "NodeContribution"}

// Get takes name of the nodeContribution, and returns the corresponding nodeContribution object, and an error if there is any.
func (c *FakeNodeContributions) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha.NodeContribution, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(nodecontributionsResource, c.ns, name), &v1alpha.NodeContribution{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha.NodeContribution), err
}

// List takes label and field selectors, and returns the list of NodeContributions that match those selectors.
func (c *FakeNodeContributions) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha.NodeContributionList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(nodecontributionsResource, nodecontributionsKind, c.ns, opts), &v1alpha.NodeContributionList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha.NodeContributionList{ListMeta: obj.(*v1alpha.NodeContributionList).ListMeta}
	for _, item := range obj.(*v1alpha.NodeContributionList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested nodeContributions.
func (c *FakeNodeContributions) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(nodecontributionsResource, c.ns, opts))

}

// Create takes the representation of a nodeContribution and creates it.  Returns the server's representation of the nodeContribution, and an error, if there is any.
func (c *FakeNodeContributions) Create(ctx context.Context, nodeContribution *v1alpha.NodeContribution, opts v1.CreateOptions) (result *v1alpha.NodeContribution, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(nodecontributionsResource, c.ns, nodeContribution), &v1alpha.NodeContribution{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha.NodeContribution), err
}

// Update takes the representation of a nodeContribution and updates it. Returns the server's representation of the nodeContribution, and an error, if there is any.
func (c *FakeNodeContributions) Update(ctx context.Context, nodeContribution *v1alpha.NodeContribution, opts v1.UpdateOptions) (result *v1alpha.NodeContribution, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(nodecontributionsResource, c.ns, nodeContribution), &v1alpha.NodeContribution{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha.NodeContribution), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeNodeContributions) UpdateStatus(ctx context.Context, nodeContribution *v1alpha.NodeContribution, opts v1.UpdateOptions) (*v1alpha.NodeContribution, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(nodecontributionsResource, "status", c.ns, nodeContribution), &v1alpha.NodeContribution{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha.NodeContribution), err
}

// Delete takes name of the nodeContribution and deletes it. Returns an error if one occurs.
func (c *FakeNodeContributions) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(nodecontributionsResource, c.ns, name), &v1alpha.NodeContribution{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeNodeContributions) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(nodecontributionsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha.NodeContributionList{})
	return err
}

// Patch applies the patch and returns the patched nodeContribution.
func (c *FakeNodeContributions) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha.NodeContribution, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(nodecontributionsResource, c.ns, name, pt, data, subresources...), &v1alpha.NodeContribution{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha.NodeContribution), err
}

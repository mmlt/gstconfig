/*
Copyright 2021.

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
	"github.com/mmlt/gstconfig/pkg/gst"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/json"

	clusteropsv1 "github.com/mmlt/gstconfig/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// GSTConfigReconciler reconciles a GSTConfig object
type GSTConfigReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	// ClusterLabelSelector selects Secrets that represent clusters
	ClusterLabelSelector labels.Selector

	// Repo is a place to store the Getting Started Tool cluster definitions.
	Repo repoer

	// Invocation counters
	reconTally int
}

// Repoer stores Getting Started Tool config data.
type repoer interface {
	Read(ctx context.Context, name string, data *gst.Config) error
	Write(ctx context.Context, name string, data *gst.Config) error
}

//+kubebuilder:rbac:groups=clusterops.mmlt.nl,resources=gstconfigs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=clusterops.mmlt.nl,resources=gstconfigs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=clusterops.mmlt.nl,resources=gstconfigs/finalizers,verbs=update
//+kubebuilder:rbac:resources=secrets,verbs=get;list;watch

// Reconcile
func (r *GSTConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx).WithName("Reconcile")
	ctx = logf.IntoContext(ctx, log)

	r.reconTally++
	log.V(1).Info("Start Reconcile", "tally", r.reconTally)
	defer log.V(1).Info("End Reconcile", "tally", r.reconTally)

	// Get the Custom Resource (deep copy).
	cr := &clusteropsv1.GSTConfig{}
	if err := r.Get(ctx, req.NamespacedName, cr); err != nil {
		log.V(2).Info("unable to get kind GSTConfig (retried)", "error", err)
		return ctrl.Result{}, ignoreNotFound(err)
	}

	/*	// Ignore environments that do not match selector.
		// (implemented as client side filtering, for server side see https://github.com/kubernetes-sigs/controller-runtime/issues/244)
		if len(r.Selector) > 0 {
			v, ok := cr.Labels[label]
			if !ok || v != r.Selector {
				log.V(2).Info("ignored, label selector doesn't match", "label", label, "value", v, "selector", r.Selector)
				return noRequeue, nil
			}
		}*/

	secrets := &corev1.SecretList{}
	err := r.List(ctx, secrets, &client.ListOptions{
		Namespace:     req.Namespace,
		LabelSelector: r.ClusterLabelSelector,
	})
	if err != nil {
		log.V(2).Info("unable to get cluster Secrets (retried)", "error", err)
		return ctrl.Result{}, ignoreNotFound(err)
	}

	targetState, err := mapSecretsToClusterConfig(secrets)
	if err != nil {
		return ctrl.Result{}, err
	}

	const filename = "clusterdefinitions.json"

	currentState := &gst.Config{}
	err = r.Repo.Read(ctx, filename, currentState)
	if err != nil {
		return ctrl.Result{}, err
	}

	newState, noChange := diff(currentState, targetState)
	if noChange {
		return ctrl.Result{}, r.updateStatusConditions(ctx, cr, noChange, nil)
	}

	err = r.Repo.Write(ctx, filename, newState)

	return ctrl.Result{}, r.updateStatusConditions(ctx, cr, noChange, err)
}

// SetupWithManager sets up the controller with the Manager.
func (r *GSTConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&clusteropsv1.GSTConfig{}).
		Complete(r)
}

// UpdateStatusConditions update CR status.conditions based on 'ready' and 'err' input.
func (r *GSTConfigReconciler) updateStatusConditions(ctx context.Context, cr *clusteropsv1.GSTConfig, ready bool, err error) error {
	c := metav1.Condition{Type: "Ready"}

	if err != nil {
		c.Status = "False"
		c.Reason = "Error"
		c.Message = err.Error()
	} else {
		if ready {
			c.Status = "True"
			c.Reason = "Synced"
		} else {
			c.Status = "False"
			c.Reason = "Syncing"
		}
	}
	meta.SetStatusCondition(&cr.Status.Conditions, c)

	return r.Status().Update(ctx, cr)
}

// MapSecretsToClusterConfig takes a list of Secrets containing kubeconfigs and returns a ClusterDefinitions object.
// A Secret contains a "cluster" and an "authInfo" field with the JSON representations similar to a kube config.
func mapSecretsToClusterConfig(secrets *corev1.SecretList) (*gst.Config, error) {
	var endpoints []gst.ClusterEndpoint
	for _, secret := range secrets.Items {
		cs, ok := secret.StringData["cluster"]
		if !ok {
			return nil, fmt.Errorf("Secret %s/%s does not contain a 'cluster' field", secret.Namespace, secret.Name)
		}

		cluster := clientcmdapi.Cluster{}
		err := json.Unmarshal([]byte(cs), &cluster)
		if err != nil {
			return nil, err
		}

		endpoints = append(endpoints, gst.ClusterEndpoint{
			Endpoint: cluster.Server,
		})
	}

	return &gst.Config{
		Clusters: endpoints,
	}, nil
}

// Diff compares the current and target state.
// It returns true if current == target state or false + a new state if current != target state.
func diff(current, target *gst.Config) (*gst.Config, bool) {
	//TODO compare current and target state by endpoints ignoring endpoints with deprecated==true

	return target, true
}

// IgnoreNotFound makes NotFound errors disappear.
// We generally want to ignore (not requeue) NotFound errors, since we'll get a reconciliation request once the
// object exists, and re-queuing in the meantime won't help.
func ignoreNotFound(err error) error {
	if apierrors.IsNotFound(err) {
		return nil
	}
	return err
}

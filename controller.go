package main

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// reconcilePod reconciles Pods
type reconcilePod struct {
	client client.Client
}

func (r *reconcilePod) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := log.FromContext(ctx)

	rs := &corev1.Pod{}
	err := r.client.Get(ctx, request.NamespacedName, rs)
	if apierrors.IsNotFound(err) {
		log.Error(nil, "Could not find Pod")
		return reconcile.Result{}, nil
	}

	if err != nil {
		return reconcile.Result{}, fmt.Errorf("could not fetch Pod: %+w", err)
	}
	log.Info("Reconciling Pod", "name", rs.Name, "namespace", rs.Namespace)
	return reconcile.Result{}, nil
}

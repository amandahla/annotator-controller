package main

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const AnnotationKey = "annotator-controller/processed"

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
		return reconcile.Result{Requeue: true}, fmt.Errorf("could not fetch Pod: %+w", err)
	}

	if rs.Annotations == nil {
		rs.Annotations = map[string]string{}
	}
	if v, ok := rs.Annotations[AnnotationKey]; ok {
		if strings.EqualFold(v, "true") {
			log.Info("Skipping pod, already annotated", "name", rs.Name, "namespace", rs.Namespace)
			return reconcile.Result{}, nil
		}
	}

	rs.Annotations[AnnotationKey] = "true"
	err = r.client.Update(ctx, rs)
	log.Info("Annotating", "name", rs.Name, "namespace", rs.Namespace)
	if err != nil {
		return reconcile.Result{Requeue: true}, fmt.Errorf("could not write Pod: %+w", err)
	}

	return reconcile.Result{}, nil
}

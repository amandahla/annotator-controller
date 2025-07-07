package main

import (
	"context"
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

func init() {
	log.SetLogger(zap.New())
}

func main() {
	cfg := config.GetConfigOrDie()
	fmt.Println("setting the manager")
	mgr, err := manager.New(cfg, manager.Options{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("registering the controller")
	c, err := controller.New("annotator-controller", mgr, controller.Options{
		Reconciler: &reconcilePod{client: mgr.GetClient()},
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("watching pods")
	err = c.Watch(source.Kind(mgr.GetCache(), &corev1.Pod{}, &handler.TypedEnqueueRequestForObject[*corev1.Pod]{}))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("starting the manager")
	err = mgr.Start(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "problem running manager: %v\n", err)
		os.Exit(1)
	}
}

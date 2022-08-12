// Copyright 2022 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package k8s

import (
	"os"
	"sync"

	tsuruv1 "github.com/tsuru/tsuru/provision/kubernetes/pkg/client/clientset/versioned/typed/tsuru/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type k8sClientGetter struct {
	cs  *kubernetes.Clientset
	tcs *tsuruv1.TsuruV1Client

	mu sync.RWMutex
}

func (g *k8sClientGetter) GetClientset() (*kubernetes.Clientset, error) {
	g.mu.RLock()
	cs := g.cs
	g.mu.RUnlock()

	if cs == nil {
		var err error
		cs, err = newClientset()
		if err != nil {
			return nil, err
		}

		g.mu.Lock()
		g.cs = cs
		g.mu.Unlock()
	}

	return g.cs, nil
}

func (g *k8sClientGetter) GetTsuruClient() (*tsuruv1.TsuruV1Client, error) {
	g.mu.RLock()
	tcs := g.tcs
	g.mu.RUnlock()

	if tcs == nil {
		var err error
		tcs, err = newTsuruClient()
		if err != nil {
			return nil, err
		}

		g.mu.Lock()
		g.tcs = tcs
		g.mu.Unlock()
	}

	return g.tcs, nil
}

func newClientset() (*kubernetes.Clientset, error) {
	config, err := newConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

func newTsuruClient() (*tsuruv1.TsuruV1Client, error) {
	config, err := newConfig()
	if err != nil {
		return nil, err
	}

	return tsuruv1.NewForConfig(config)
}

func newConfig() (*rest.Config, error) {
	config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		return nil, err
	}

	return config, nil
}

var Default k8sClientGetter

func GetClientset() (*kubernetes.Clientset, error) { return Default.GetClientset() }

func GetTsuruClient() (*tsuruv1.TsuruV1Client, error) { return Default.GetTsuruClient() }

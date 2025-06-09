package controller

import (
	"context"
	"os"

	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// k8sNamespaceFilename file inside the pod that shows in which namespace the pod is running
const k8sNamespaceFilename = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

// getSecretByName get a secret from secret distributor namespace by name
func (r *DistributionReconciler) getSecretByName(secretName string) (*v1.Secret, error) {
	sdNamesapce, err := getPodNamespace()
	if err != nil {
		return nil, err
	}

	secret := &v1.Secret{}
	if err := r.Get(context.Background(), client.ObjectKey{
		Namespace: sdNamesapce,
		Name:      secretName,
	}, secret); err != nil {
		return nil, err
	}
	return secret, nil
}

// getPodNamespace will get secret-distributor namespace from a file inside pod (only in kubernetes)
func getPodNamespace() (string, error) {
	bytes, err := os.ReadFile(k8sNamespaceFilename)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

package controller

import (
	"context"
	"os"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// k8sNamespaceFilename file inside the pod that shows in which namespace the pod is running
const k8sNamespaceFilename = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

// readFromSecretAndCreateNewSecret get a secret from secret distributor namespace by name
func (r *DistributionReconciler) readFromSecretAndCreateNewSecret(secretName, targetNamespace string) (*v1.Secret, error) {
	distributerNamesapce, err := getPodNamespace()
	if err != nil {
		return nil, err
	}

	secret := &v1.Secret{}
	if err := r.Get(context.Background(), client.ObjectKey{
		Namespace: distributerNamesapce,
		Name:      secretName,
	}, secret); err != nil {
		return nil, err
	}

	newSecret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secret.Name,
			Namespace: targetNamespace,
			Labels:    secret.Labels,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		Data:       secret.Data,
		StringData: secret.StringData,
		Type:       secret.Type,
	}
	return newSecret, nil
}

// getPodNamespace will get secret-distributor namespace from a file inside pod (only in kubernetes)
func getPodNamespace() (string, error) {
	bytes, err := os.ReadFile(k8sNamespaceFilename)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

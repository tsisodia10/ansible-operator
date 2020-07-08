package tools

import (
	miqv1alpha1 "github.com/ansible-operator/pkg/apis/ansible/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DefaultAnsibleSecret(cr *ansiblev1alpha1.AnsiblePlaybook) *corev1.Secret {
	labels := map[string]string{
		"app": cr.Spec.PlaybookName,
	}
	secret := map[string]string{
		"username": "root",
		"password": generatePassword(),
		"hostname": "localhost",
	}

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ansibleSecretName(cr),
			Namespace: cr.ObjectMeta.Namespace,
			Labels:    labels,
		},
		StringData: secret,
	}
}

func ansibleSecretName(cr *miqv1alpha1.ManageIQ) string {
	secretName := "ansible-secrets"
	if cr.Spec.DatabaseSecret != "" {
		secretName = cr.Spec.DatabaseSecret
	}

	return secretName
}

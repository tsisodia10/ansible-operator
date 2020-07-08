package tools

import (
	"strconv"
	"strings"

	ansiblev1alpha1 "github.com/ansible-operator/pkg/apis/ansible/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func OrchestratorJob(cr *ansiblev1alpha1.AnsiblePlaybook, scheme *runtime.Scheme) (*appsv1.Job, controllerutil.MutateFn, error) {
	delaySecs, err := strconv.Atoi(cr.Spec.OrchestratorInitialDelay)
	if err != nil {
		return nil, nil, err
	}
	pullPolicy := corev1.PullIfNotPresent
	if strings.Contains(cr.Spec.OrchestratorImageTag, "latest") {
		pullPolicy = corev1.PullAlways
	}

	deploymentLabels := map[string]string{
		"name": "orchestrator",
		"app":  cr.Spec.AppName,
	}

	container := corev1.Container{
		Name:            "orchestrator",
		Image:           cr.Spec.OrchestratorImageNamespace + "/" + cr.Spec.OrchestratorImageName + ":" + cr.Spec.OrchestratorImageTag,
		ImagePullPolicy: pullPolicy,
		LivenessProbe: &corev1.Probe{
			Handler: corev1.Handler{
				Exec: &corev1.ExecAction{
					Command: []string{},
				},
			},
			InitialDelaySeconds: int32(delaySecs),
			TimeoutSeconds:      3,
		},
		Env: []corev1.EnvVar{
			corev1.EnvVar{
				Name:  "ALLOW_INSECURE_SESSION",
				Value: "true",
			},
			corev1.EnvVar{
				Name:  "ANSIBLEPLAYBOOK_NAME",
				Value: cr.Spec.PlaybookName,
			},
			corev1.EnvVar{
				Name:  "REPOSITORY_TYPE",
				Value: cr.Spec.RepositoryType,
			},
			corev1.EnvVar{
				Name:  "REPOSITORY_URL",
				Value: cr.Spec.RepositoryURL,
			},
			corev1.EnvVar{
				Name:  "INVENTORY",
				Value: cr.Spec.Inventory,
			},
			corev1.EnvVar{
				Name:  "DATABASE_REGION",
				Value: cr.Spec.DatabaseRegion,
			},
			corev1.EnvVar{
				Name: "HOST_CREDENTIAL",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{Name: ansibleSecretName(cr)},
						Key:                  "password",
					},
				}, Value: cr.Spec.ImagePullSecret,
			},
			corev1.EnvVar{
				Name: "ENCRYPTION_KEY",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{Name: "app-secrets"},
						Key:                  "encryption-key",
					},
				},
			},
			corev1.EnvVar{
				Name: "JOB_NAME",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{FieldPath: "metadata.name"},
				},
			},
			corev1.EnvVar{
				Name: "JOB_UID",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{FieldPath: "metadata.uid"},
				},
			},
		},
	}

	job := &corev1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orchestrator",
			Namespace: cr.ObjectMeta.Namespace,
		},
		Spec: corev1.JobSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: deploymentLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: deploymentLabels,
					Name:   "orchestrator",
				},
				Spec: corev1.PodSpec{},
			},
		},
	}

	f := func() error {
		if err := controllerutil.SetControllerReference(cr, job, scheme); err != nil {
			return err
		}
		addAppLabel(cr.Spec.AppName, &job.ObjectMeta)

		job.Spec.Template.Spec.Containers = []corev1.Container{container}

		if cr.Spec.ImagePullSecret != "" {
			pullSecret := []corev1.LocalObjectReference{
				corev1.LocalObjectReference{Name: cr.Spec.ImagePullSecret},
			}
			job.Spec.Template.Spec.ImagePullSecrets = pullSecret

			c := &job.Spec.Template.Spec.Containers[0]
			pullSecretEnv := corev1.EnvVar{
				Name:  "IMAGE_PULL_SECRET",
				Value: cr.Spec.ImagePullSecret,
			}
			c.Env = append(c.Env, pullSecretEnv)
		}

		return nil
	}

	return job, f, nil
}

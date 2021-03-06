package kubermatic

import (
	"fmt"

	operatorv1alpha1 "github.com/kubermatic/kubermatic/api/pkg/crd/operator/v1alpha1"
	"github.com/kubermatic/kubermatic/api/pkg/resources/reconciling"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
)

func masterControllerManagerPodLabels() map[string]string {
	return map[string]string{
		nameLabel:    "master-controller-manager",
		versionLabel: "v1",
	}
}

func MasterControllerManagerDeploymentCreator(cfg *operatorv1alpha1.KubermaticConfiguration) reconciling.NamedDeploymentCreatorGetter {
	return func() (string, reconciling.DeploymentCreator) {
		return masterControllerManagerDeploymentName, func(d *appsv1.Deployment) (*appsv1.Deployment, error) {
			d.Spec.Replicas = pointer.Int32Ptr(2)
			d.Spec.Selector = &metav1.LabelSelector{
				MatchLabels: masterControllerManagerPodLabels(),
			}

			d.Spec.Template.Labels = d.Spec.Selector.MatchLabels
			d.Spec.Template.Annotations = map[string]string{
				"prometheus.io/scrape": "true",
				"prometheus.io/port":   "8085",
				"fluentbit.io/parser":  "glog",
			}

			d.Spec.Template.Spec.ServiceAccountName = serviceAccountName
			d.Spec.Template.Spec.ImagePullSecrets = []corev1.LocalObjectReference{
				{
					Name: dockercfgSecretName,
				},
			}

			args := []string{
				"-v=4",
				"-logtostderr",
				"-internal-address=0.0.0.0:8085",
				"-kubeconfig=/opt/.kube/kubeconfig",
			}

			volumes := []corev1.Volume{
				{
					Name: "kubeconfig",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: kubeconfigSecretName,
						},
					},
				},
			}

			volumeMounts := []corev1.VolumeMount{
				{
					MountPath: "/opt/.kube/",
					Name:      "kubeconfig",
					ReadOnly:  true,
				},
			}

			if cfg.Spec.Datacenters != "" {
				args = append(args, "-datacenters=/opt/datacenters/datacenters.yaml")

				volumes = append(volumes, corev1.Volume{
					Name: "datacenters",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: datacentersSecretName,
						},
					},
				})

				volumeMounts = append(volumeMounts, corev1.VolumeMount{
					MountPath: "/opt/datacenters/",
					Name:      "datacenters",
					ReadOnly:  true,
				})
			}

			d.Spec.Template.Spec.Volumes = volumes
			d.Spec.Template.Spec.InitContainers = []corev1.Container{projectsMigratorContainer(cfg)}
			d.Spec.Template.Spec.Containers = []corev1.Container{
				{
					Name:    "controller-manager",
					Image:   cfg.Spec.MasterController.Image,
					Command: []string{"master-controller-manager"},
					Args:    args,
					Ports: []corev1.ContainerPort{
						{
							Name:          "metrics",
							ContainerPort: 8085,
							Protocol:      corev1.ProtocolTCP,
						},
					},
					VolumeMounts: volumeMounts,
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("50m"),
							corev1.ResourceMemory: resource.MustParse("128Mi"),
						},
						Limits: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("100m"),
							corev1.ResourceMemory: resource.MustParse("256Mi"),
						},
					},
				},
			}

			return d, nil
		}
	}
}

func projectsMigratorContainer(cfg *operatorv1alpha1.KubermaticConfiguration) corev1.Container {
	return corev1.Container{
		Name:    "projects-migrator",
		Image:   cfg.Spec.MasterController.Image,
		Command: []string{"projects-migrator"},
		Args: []string{
			"-v=2",
			"-logtostderr",
			"-kubeconfig=/opt/.kube/kubeconfig",
			fmt.Sprintf("-dry-run=%v", cfg.Spec.MasterController.ProjectsMigrator.DryRun),
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				MountPath: "/opt/.kube/",
				Name:      "kubeconfig",
				ReadOnly:  true,
			},
		},
	}
}

func MasterControllerManagerPDBCreator(cfg *operatorv1alpha1.KubermaticConfiguration) reconciling.NamedPodDisruptionBudgetCreatorGetter {
	name := "kubermatic-master-controller-manager-v1"

	return func() (string, reconciling.PodDisruptionBudgetCreator) {
		return name, func(pdb *policyv1beta1.PodDisruptionBudget) (*policyv1beta1.PodDisruptionBudget, error) {
			min := intstr.FromInt(1)

			pdb.Spec.MinAvailable = &min
			pdb.Spec.Selector = &metav1.LabelSelector{
				MatchLabels: masterControllerManagerPodLabels(),
			}

			return pdb, nil
		}
	}
}

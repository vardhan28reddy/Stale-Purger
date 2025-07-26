package k8s

import (
	"context"

	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type KubeClient struct {
	KubeClient kubernetes.Interface
	Logger     *logrus.Entry
}

func (k KubeClient) ListNamespaces(ctx context.Context) (*corev1.NamespaceList, error) {
	return k.KubeClient.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
}

func (k KubeClient) ListPods(ctx context.Context, namespace string) (*corev1.PodList, error) {
	return k.KubeClient.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
}

func (k KubeClient) DeletePod(ctx context.Context, namespace, podName string) error {
	return k.KubeClient.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
}

func (k KubeClient) GetDeployment(ctx context.Context, namespace, deploymentName string) (*appsv1.Deployment, error) {
	return k.KubeClient.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
}

func (k KubeClient) GetReplicaSet(ctx context.Context, namespace, replicaSetName string) (*appsv1.ReplicaSet, error) {
	return k.KubeClient.AppsV1().ReplicaSets(namespace).Get(context.TODO(), replicaSetName, metav1.GetOptions{})
}

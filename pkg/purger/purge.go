package purger

import (
	"Stale-purger/pkg/consts"
	"Stale-purger/pkg/db"
	"Stale-purger/pkg/k8s"
	"Stale-purger/pkg/utils"
	"context"
	"database/sql"
	"time"

	"github.com/sirupsen/logrus"
)

type PodInfo struct {
	PodName string
	Phase   string
	Reason  string
}

type StalePodsInfo struct {
	Info       map[string][]PodInfo
	PostgresDB *sql.DB
}

// Function to check if all desired replicas are available
func replicasAreStable(ctx context.Context, kubeClient k8s.KubeClient, namespace, resourceName, resourceType string, logger *logrus.Entry) (bool, error) {
	switch resourceType {
	case string(consts.DeploymentResource):
		deploy, err := kubeClient.GetDeployment(ctx, namespace, resourceName)
		logger.WithFields(logrus.Fields{"namespace": namespace, "resource": resourceType, "availableReplicas": deploy.Status.AvailableReplicas, "desiredReplicas": *deploy.Spec.Replicas}).Debug("Replica information")
		return deploy.Status.AvailableReplicas == *deploy.Spec.Replicas, err
	case string(consts.ReplicaSetResource):
		rs, err := kubeClient.GetReplicaSet(ctx, namespace, resourceName)
		logger.WithFields(logrus.Fields{"namespace": namespace, "resource": resourceType, "availableReplicas": rs.Status.AvailableReplicas, "desiredReplicas": *rs.Spec.Replicas}).Debug("Replica information")
		return rs.Status.AvailableReplicas == *rs.Spec.Replicas, err
	default:
		return false, nil
	}
}

// PurgeStalePods Function to find and delete stale pods, If AvailableReplicas == DesiredReplicas.
func PurgeStalePods(ctx context.Context, stalePodsInfo StalePodsInfo, kubeClient k8s.KubeClient, namespace string, logger *logrus.Entry) (StalePodsInfo, error) {
	pods, err := kubeClient.ListPods(context.TODO(), namespace)
	if err != nil {
		return StalePodsInfo{}, err
	}
	for _, pod := range pods.Items {
		// Check if the pod is in one of the target consts
		if utils.StalePodStates[pod.Status.Phase] {
			dbPodIndo := &db.PodInfo{
				PodName:   pod.Name,
				Namespace: namespace,
				NodeName:  pod.Spec.NodeName,
				Reason:    pod.Status.Reason,
				Status:    string(pod.Status.Phase),
			}
			// Find pod owner (Deployment or ReplicaSet)
			if len(pod.OwnerReferences) == 0 {
				dbPodIndo.ActionType = string(consts.SkipOnNoOwnerActionType)
				if err := db.NewStalePurgerDB(stalePodsInfo.PostgresDB).SaveStalePodInfo(dbPodIndo); err != nil {
					return StalePodsInfo{}, err
				}
				logger.WithFields(logrus.Fields{"namespace": namespace, "pod": pod.Name}).Debug("Pod has no owner, skipping...")
				continue
			}
			resourceType := pod.OwnerReferences[0]
			dbPodIndo.OwnerName = resourceType.Name
			dbPodIndo.OwnerType = resourceType.Kind
			if consts.ResourceType(resourceType.Kind) == consts.DeploymentResource || consts.ResourceType(resourceType.Kind) == consts.ReplicaSetResource {
				isStable, err := replicasAreStable(ctx, kubeClient, namespace, resourceType.Name, resourceType.Kind, logger)
				if err != nil {
					logger.WithFields(logrus.Fields{"namespace": namespace, "resource": resourceType.Name, "pod": pod.Name}).Debug("Error in fetching Resource")
					return StalePodsInfo{}, err
				} else if isStable && (string(pod.Status.Phase) == string(consts.UnknownPodPhase) || string(pod.Status.Phase) == string(consts.FailedPodPhase)) {
					logger.WithFields(logrus.Fields{"namespace": namespace, "resource": resourceType.Name, "pod": pod.Name}).Info("Deleting pod related to resource")
					dbPodIndo.DeletedAt = time.Now()
					if err := kubeClient.DeletePod(ctx, namespace, pod.Name); err != nil {
						logger.Errorf("Failed to delete pod %s: %v", pod.Name, err)
						return StalePodsInfo{}, err
					}
					logger.WithFields(logrus.Fields{"namespace": namespace, "resource": resourceType.Name, "pod": pod.Name}).Info("Successfully deleted pod")
					dbPodIndo.ActionType = string(consts.DeleteActionType)
					if err := db.NewStalePurgerDB(stalePodsInfo.PostgresDB).SaveStalePodInfo(dbPodIndo); err != nil {
						return StalePodsInfo{}, err
					}
					// Initialize the slice only if it doesn't exist
					if _, exists := stalePodsInfo.Info[namespace]; !exists {
						stalePodsInfo.Info[namespace] = []PodInfo{}
					}
					// Store the pod info for metrics display
					stalePodsInfo.Info[namespace] = append(stalePodsInfo.Info[namespace], PodInfo{PodName: pod.Name, Phase: string(pod.Status.Phase), Reason: pod.Status.Reason})
				} else {
					logger.WithFields(logrus.Fields{"namespace": namespace, "resource": resourceType.Name, "pod": pod.Name}).Debug("Skipping deletion, replicas not stable")
					dbPodIndo.ActionType = string(consts.SkipOnNotStableActionType)
					if err := db.NewStalePurgerDB(stalePodsInfo.PostgresDB).SaveStalePodInfo(dbPodIndo); err != nil {
						return StalePodsInfo{}, err
					}
				}
			}
		}
	}
	return stalePodsInfo, nil
}

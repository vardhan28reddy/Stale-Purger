package db

import (
	"time"

	"github.com/google/uuid"
)

const (
	deletedPodTableName             = "deleted_pods"
	queryColumnsFromDeletedPodTable = "pod_name, namespace, owner_type, owner_name, deleted_at, deletion_reason, status"
	
)
const (
	QueryColumnsFromDeletedPodTable = "id, pod_name, namespace, node_name, owner_type, owner_name, deleted_at, action_type, deletion_reason, status"
)
type PodInfo struct {
	ID         uuid.UUID // Automatically generated UUID
	PodName    string
	Namespace  string
	NodeName   string
	OwnerType  string
	OwnerName  string
	ActionType string
	DeletedAt  time.Time
	Reason     string
	Status     string
	PodAgeW    int
}

type PodQueryInfo struct {
	PodName   string
	Namespace string
	OwnerType string
	OwnerName string
	Reason    string
	Status    string
}

package db

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

type StalePurgerDB struct {
	db *sql.DB
	sb sq.StatementBuilderType
}

func NewStalePurgerDB(DB *sql.DB) *StalePurgerDB {
	return &StalePurgerDB{
		db: DB,
		sb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(DB),
	}
}

func (s StalePurgerDB) SaveStalePodInfo(pod *PodInfo) error {
	query := `INSERT INTO deleted_pods (pod_name, namespace, node_name, owner_type, owner_name, deleted_at, action_type, deletion_reason, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := s.db.Exec(query, pod.PodName, pod.Namespace, pod.NodeName, pod.OwnerType, pod.OwnerName, pod.DeletedAt, pod.ActionType, pod.Reason, pod.Status)
	return err
}

func (s StalePurgerDB) GetStalePodInfoOnActionType(actionType string) ([]*PodQueryInfo, error) {
	var podActionTypeInfo []*PodQueryInfo
	rows, err := s.sb.Select(queryColumnsFromDeletedPodTable).From(deletedPodTableName).Where("action_type = ?", actionType).OrderBy("deleted_at DESC").Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := PodQueryInfo{}
		if err := rows.Scan(
			&item.PodName,
			&item.Namespace,
			&item.OwnerType,
			&item.OwnerName,
			&item.Reason,
			&item.Status); err != nil {
			return nil, err
		}
		podActionTypeInfo = append(podActionTypeInfo, &item)
	}
	return podActionTypeInfo, nil
}

func (s StalePurgerDB) GetStalePodInfoOnNamespace(namespace string) ([]*PodQueryInfo, error) {
	var podActionTypeInfo []*PodQueryInfo
	rows, err := s.sb.Select(queryColumnsFromDeletedPodTable).From(deletedPodTableName).Where("namespace = ?", namespace).OrderBy("deleted_at DESC").Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := PodQueryInfo{}
		if err := rows.Scan(
			&item.PodName,
			&item.Namespace,
			&item.OwnerType,
			&item.OwnerName,
			&item.Reason,
			&item.Status); err != nil {
			return nil, err
		}
		podActionTypeInfo = append(podActionTypeInfo, &item)
	}
	return podActionTypeInfo, nil
}

func (s StalePurgerDB) GetStalePodInfoOnNodeName(nodeName string) ([]*PodQueryInfo, error) {
	var podActionTypeInfo []*PodQueryInfo
	rows, err := s.sb.Select(queryColumnsFromDeletedPodTable).From(deletedPodTableName).Where("node_name = ?", nodeName).OrderBy("deleted_at DESC").Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := PodQueryInfo{}
		if err := rows.Scan(
			&item.PodName,
			&item.Namespace,
			&item.OwnerType,
			&item.OwnerName,
			&item.Reason,
			&item.Status); err != nil {
			return nil, err
		}
		podActionTypeInfo = append(podActionTypeInfo, &item)
	}
	return podActionTypeInfo, nil
}

func (s StalePurgerDB) GetStalePodInfoOnOwnerType(ownerType string) ([]*PodQueryInfo, error) {
	var podActionTypeInfo []*PodQueryInfo
	rows, err := s.sb.Select(queryColumnsFromDeletedPodTable).From(deletedPodTableName).Where("owner_type = ?", ownerType).OrderBy("deleted_at DESC").Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := PodQueryInfo{}
		if err := rows.Scan(
			&item.PodName,
			&item.Namespace,
			&item.OwnerType,
			&item.OwnerName,
			&item.Reason,
			&item.Status); err != nil {
			return nil, err
		}
		podActionTypeInfo = append(podActionTypeInfo, &item)
	}
	return podActionTypeInfo, nil
}

func (s StalePurgerDB) GetStalePodsInfoOnStatus(status string) ([]*PodQueryInfo, error) {
	var podActionTypeInfo []*PodQueryInfo
	rows, err := s.sb.Select(queryColumnsFromDeletedPodTable).From(deletedPodTableName).Where("status = ?", status).OrderBy("deleted_at DESC").Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := PodQueryInfo{}
		if err := rows.Scan(
			&item.PodName,
			&item.Namespace,
			&item.OwnerType,
			&item.OwnerName,
			&item.Reason,
			&item.Status); err != nil {
			return nil, err
		}
		podActionTypeInfo = append(podActionTypeInfo, &item)
	}
	return podActionTypeInfo, nil
}

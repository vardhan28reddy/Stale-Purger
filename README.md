# Stale-purger

**Stale-purger** is a Go-based tool designed to identify and clean up stale pods in Kubernetes clusters. It ensures that pods in undesired states (such as `Failed` or `Unknown`) are deleted only when their parent resources (Deployments or ReplicaSets) have stable replicas, helping maintain cluster state and resource efficiency.

## About

- Scans Kubernetes namespaces for pods in undesired states.
- Checks if parent resources (Deployment/ReplicaSet) have all desired replicas available before deleting pods.
- Logs actions and reasons for skipping or deleting pods.
- Saves pod actions and metadata to a PostgreSQL database for auditing.
- Provides metrics for deleted and skipped pods.

## How It Works

1. **Pod Scanning:**  
   The tool lists all pods in a namespace and checks their status against a set of undesired states (e.g., `Failed`, `Unknown`).

2. **Owner Verification:**  
   For each stale pod, it determines the owner resource (Deployment/ReplicaSet). If no owner is found, the pod is skipped.

3. **Replica Stability Check:**  
   Before deleting a pod, it verifies that the owner resource has all desired replicas available.

4. **Pod Deletion & Logging:**  
   If conditions are met, the pod is deleted and the action is logged in the database. Otherwise, the reason for skipping is recorded.


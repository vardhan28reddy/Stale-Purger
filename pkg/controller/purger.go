package controller

import (
	"Stale-purger/pkg/config"
	"Stale-purger/pkg/k8s"
	"Stale-purger/pkg/purger"
	"Stale-purger/pkg/utils"
	"context"
	"database/sql"
	"time"

	"github.com/sirupsen/logrus"
)

type PurgerComponent struct {
	db         *sql.DB
	logger     *logrus.Entry
	config     config.Config
	kubeClient k8s.KubeClient
}

func NewPurgerComponent(db *sql.DB, logger *logrus.Entry, kubeClient k8s.KubeClient, config config.Config) Component {
	return &PurgerComponent{db: db, logger: logger, config: config, kubeClient: kubeClient}
}

func (p *PurgerComponent) Name() string { return "Purger Handler" }

func (p *PurgerComponent) Start(ctx context.Context) {

	execPurgeStalePodsFunc := func() {
		p.logger.Info("Looking for stale pods")
		namespaces, err := p.kubeClient.ListNamespaces(ctx)
		utils.FatalFunc("Couldn't list namespaces", err, p.logger)
		stalePodsInfo := purger.StalePodsInfo{Info: make(map[string][]purger.PodInfo, 0), PostgresDB: p.db}
		for _, namespace := range namespaces.Items {
			p.logger.WithFields(logrus.Fields{"namespace": namespace.Name, "component": "state-keeper"}).Info("Scanning namespace")
			stalePodsInfo, err = purger.PurgeStalePods(ctx, stalePodsInfo, p.kubeClient, namespace.Name, p.logger)
			if err != nil {
				p.logger.Warnf("Couldn't purge the pods: %+v", err)
			}
		}
		p.logger.Infof("Purged info: %+v", stalePodsInfo)
	}
	execPurgeStalePodsFunc()
	ticker := time.NewTicker(24 * time.Hour)

loop:

	for {
		select {
		case <-ticker.C:
			execPurgeStalePodsFunc()
		case <-ctx.Done():
			ticker.Stop()
			break loop
		}
	}
}

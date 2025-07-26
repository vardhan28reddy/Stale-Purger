package config

import (
	"Stale-purger/pkg/k8s"
	"database/sql"
	"fmt"
	"net/url"
	"path/filepath"
	"time"

	_ "github.com/lib/pq"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Config struct {
	GenerateHTMLOutput bool   `default:"true"`
	LogLevel           string `default:"debug"`
	LogJsonFormat      bool   `default:"true"`
	KubernetesConfig   KubernetesConfig
	PostgresConfig     PostgresConfig
}

type KubernetesConfig struct {
	InCluster bool `default:"false"`
}

type PostgresConfig struct {
	Host     string `required:"true" default:"localhost"`
	Port     int    `required:"true,gte=1,lte=65535" default:"5435"`
	Username string `required:"true" default:"postgres"`
	Password string `required:"true" default:"postgres"`
	Database string `required:"true" default:"stale-pods-info"`
	Limit    struct {
		Conn struct {
			Idle        int `required:"true" default:"10"`
			Open        int `required:"true" default:"10"`
			MaxLifetime int `required:"true" default:"300" split_words:"true"`
		}
	}
	Sslmode     string `default:"disable"`
	SslRootCert string `default:"/opt/postgresql/ssl/postgres-root.pem"`
}

func InitializeConfig() (*Config, error) {
	s := &Config{}
	err := envconfig.Process("", s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (c Config) InitializeKubeClient(logger *log.Entry) (k8s.KubeClient, error) {
	var kubeConfig *rest.Config
	var err error
	if c.KubernetesConfig.InCluster {
		if kubeConfig, err = rest.InClusterConfig(); err != nil {
			return k8s.KubeClient{}, err
		}
	} else {
		if kubeConfig, err = clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config")); err != nil {
			return k8s.KubeClient{}, err
		}
	}

	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return k8s.KubeClient{}, err
	}
	return k8s.KubeClient{KubeClient: kubeClient, Logger: logger}, nil

}

func InitializeLogger(logLevel string, jsonFormat bool) *log.Entry {
	logger := log.New()

	var formatter log.Formatter
	if jsonFormat {
		formatter = &log.JSONFormatter{}
	} else {
		formatter = &log.TextFormatter{
			FullTimestamp:          true,
			DisableLevelTruncation: false,
		}
	}

	logger.SetFormatter(formatter)

	if level, err := log.ParseLevel(logLevel); err == nil {
		logger.SetLevel(level)
	} else {
		logger.SetLevel(log.InfoLevel)
	}

	return logger.WithFields(log.Fields{
		"service": "cluster-state-keeper",
	})
}

func (c Config) InitializeDB(logger *log.Entry) (*sql.DB, error) {
	dbConnectionURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s&sslrootcert=%s", c.PostgresConfig.Username, url.QueryEscape(c.PostgresConfig.Password), c.PostgresConfig.Host, c.PostgresConfig.Port, c.PostgresConfig.Database, c.PostgresConfig.Sslmode, c.PostgresConfig.SslRootCert)
	db, err := sql.Open("postgres", dbConnectionURL)
	if err != nil {
		return db, errors.Wrap(err, "failed to open postgres connection")
	}
	if err := db.Ping(); err != nil {
		return db, errors.Wrap(err, "failed to ping postgres")
	}

	db.SetMaxIdleConns(c.PostgresConfig.Limit.Conn.Idle)
	db.SetMaxOpenConns(c.PostgresConfig.Limit.Conn.Open)
	db.SetConnMaxLifetime(time.Duration(c.PostgresConfig.Limit.Conn.MaxLifetime) * time.Second)
	return db, nil
}

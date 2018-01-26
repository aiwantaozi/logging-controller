package v3

import (
	"github.com/rancher/norman/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Logging struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object’s metadata. More info:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// Specification of the desired behavior of the the cluster. More info:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#spec-and-status
	Spec LoggingSpec `json:"spec"`
	// Most recent observed status of the cluster. More info:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#spec-and-status
	Status      LoggingStatus `json:"status"`
	DisplayName string        `json:"displayName,omitempty"`
	ClusterName string        `json:"clusterName" norman:"type=reference[cluster]"`
}

type ProjectLogging struct {
	types.Namespaced
	metav1.TypeMeta `json:",inline"`
	// Standard object’s metadata. More info:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// Specification of the desired behavior of the the cluster. More info:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#spec-and-status
	Spec ProjectLoggingSpec `json:"spec"`
	// Most recent observed status of the cluster. More info:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#spec-and-status
	Status      LoggingStatus `json:"status"`
	DisplayName string        `json:"displayName,omitempty"`
	ProjectName string        `json:"projectName,omitempty" norman:"type=reference[project]"`
}

type LoggingSpec struct {
	CurrentTarget       string            `json:"currentTarget"`
	OutputFlushInterval int               `json:"outputFlushInterval"`
	OutputTags          map[string]string `json:"outputTags"`

	EmbeddedConfig      *EmbeddedConfig      `json:"embeddedConfig,omitempty"`
	ElasticsearchConfig *ElasticsearchConfig `json:"elasticsearchConfig,omitempty"`
	SplunkConfig        *SplunkConfig        `json:"splunkConfig,omitempty"`
	KafkaConfig         *KafkaConfig         `json:"kafkaConfig,omitempty"`
	SyslogConfig        *SyslogConfig        `json:"syslogConfig,omitempty`
}

type ProjectLoggingSpec struct {
	CurrentTarget       string            `json:"currentTarget"`
	OutputFlushInterval int               `json:"outputFlushInterval"`
	OutputTags          map[string]string `json:"outputTags"`

	ElasticsearchConfig *ElasticsearchConfig `json:"elasticsearchConfig,omitempty"`
	SplunkConfig        *SplunkConfig        `json:"splunkConfig,omitempty"`
	KafkaConfig         *KafkaConfig         `json:"kafkaConfig,omitempty"`
	SyslogConfig        *SyslogConfig        `json:"syslogConfig,omitempty`
}
type LoggingStatus struct {
	CurrentTarget string `json:"currentTarget"`
	//todo
}

type ElasticsearchConfig struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	IndexPrefix  string `json:"indexPrefix"`
	Dateformat   string `json:"dateformat"`
	AuthUser     string `json:"authUser"`     //secret
	AuthPassword string `json:"authPassword"` //secret
}

type SplunkConfig struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Protocol   string `json:"protocol"`
	Source     string `json:"source"`
	TimeFormat string `json:"timeFormat"`
	Token      string `json:"token"` //secret
}

type EmbeddedConfig struct {
	KibanaAccessURL        string `json:"kibanaAccessURL"`
	ElasticsearchAccessURL string `json:"elasticsearchAccessURL"`
}

type KafkaConfig struct {
	BrokerType     string   `json:"brokerType"`
	ZookeeperHost  string   `json:"zookeeperHost"`
	ZookeeperPort  int      `json:"zookeeperPort"`
	Brokers        []string `json:"brokers"`
	Topic          string   `json:"topic"`
	DataType       string   `json:"dataType"`
	MaxSendRetries int      `json:"maxSendRetries"`
}

type SyslogConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Severity string `json:"severity"`
	Program  string `json:"program"`
}

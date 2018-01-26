package schema

import (
	"github.com/rancher/norman/types"
	"github.com/rancher/types/apis/logging.cattle.io/v3"
	"github.com/rancher/types/factory"
	"github.com/rancher/types/mapper"
)

var (
	Version = types.APIVersion{
		Version: "v3",
		Group:   "logging.cattle.io",
		Path:    "/v3",
	}

	Schemas = factory.Schemas(&Version).
		Init(loggingTypes)
)

func loggingTypes(schemas *types.Schemas) *types.Schemas {
	return schemas.
		AddMapperForType(&Version, v3.ProjectLogging{},
			&mapper.NamespaceIDMapper{},
		).
		MustImportAndCustomize(&Version, v3.ElasticsearchConfig{}, func(schema *types.Schema) {
			schema.MustCustomizeField("dateformat", func(field types.Field) types.Field {
				field.Type = "enum"
				field.Options = []string{"YYYY.MM.DD", "YYYY.MM", "YYYY"}
				field.Default = "YYYY.MM.DD"
				return field
			})
			schema.MustCustomizeField("port", func(field types.Field) types.Field {
				field.Default = 9200
				return field
			})
		}).
		MustImportAndCustomize(&Version, v3.SplunkConfig{}, func(schema *types.Schema) {
			schema.MustCustomizeField("protocol", func(field types.Field) types.Field {
				field.Type = "enum"
				field.Options = []string{"http", "https"}
				field.Default = "http"
				return field
			})
			schema.MustCustomizeField("timeFormat", func(field types.Field) types.Field {
				field.Type = "enum"
				field.Options = []string{"unixtime", "localtime", "none"}
				field.Default = "unixtime"
				return field
			})
			schema.MustCustomizeField("port", func(field types.Field) types.Field {
				field.Default = 8088
				return field
			})
		}).
		MustImportAndCustomize(&Version, v3.KafkaConfig{}, func(schema *types.Schema) {
			schema.MustCustomizeField("brokerType", func(field types.Field) types.Field {
				field.Type = "enum"
				field.Options = []string{"broker", "zookeeper"}
				field.Default = "zookeeper"
				return field
			})
			schema.MustCustomizeField("zookeeperPort", func(field types.Field) types.Field {
				field.Default = 2181
				return field
			})
			schema.MustCustomizeField("maxSendRetries", func(field types.Field) types.Field {
				field.Default = 1
				return field
			})
			schema.MustCustomizeField("dataType", func(field types.Field) types.Field {
				field.Type = "enum"
				field.Options = []string{"json", "ltsv", "msgpack"}
				field.Default = "json"
				return field
			})
		}).
		MustImportAndCustomize(&Version, v3.SyslogConfig{}, func(schema *types.Schema) {
			schema.MustCustomizeField("port", func(field types.Field) types.Field {
				field.Default = 51400
				return field
			})
			schema.MustCustomizeField("severity", func(field types.Field) types.Field {
				field.Type = "enum"
				field.Options = []string{"emerg", "alert", "crit", "err", "warning", "notice", "info", "debug"}
				field.Default = "notice"
				return field
			})
		}).
		MustImport(&Version, v3.Logging{}).
		MustImport(&Version, v3.ProjectLogging{})
}

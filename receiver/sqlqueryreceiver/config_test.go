// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package sqlqueryreceiver

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap/confmaptest"
	"go.opentelemetry.io/collector/confmap/xconfmap"
	"go.opentelemetry.io/collector/scraper/scraperhelper"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/sqlquery"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sqlqueryreceiver/internal/metadata"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		fname        string
		id           component.ID
		expected     component.Config
		errorMessage string
	}{
		{
			id:    component.NewIDWithName(metadata.Type, ""),
			fname: "config.yaml",
			expected: &Config{
				Config: sqlquery.Config{
					ControllerConfig: scraperhelper.ControllerConfig{
						CollectionInterval: 10 * time.Second,
						InitialDelay:       time.Second,
					},
					Driver:     "postgres",
					DataSource: "host=localhost port=5432 user=me password=s3cr3t sslmode=disable",
					Queries: []sqlquery.Query{
						{
							SQL: "select count(*) as count, type from mytable group by type",
							Metrics: []sqlquery.MetricCfg{
								{
									MetricName:       "val.count",
									ValueColumn:      "count",
									AttributeColumns: []string{"type"},
									Monotonic:        false,
									ValueType:        sqlquery.MetricValueTypeInt,
									DataType:         sqlquery.MetricTypeSum,
									Aggregation:      sqlquery.MetricAggregationCumulative,
									StaticAttributes: map[string]string{"foo": "bar"},
								},
							},
						},
					},
				},
			},
		},
		{
			id:    component.NewIDWithName(metadata.Type, ""),
			fname: "config-datasource-config.yaml",
			expected: &Config{
				Config: sqlquery.Config{
					ControllerConfig: scraperhelper.ControllerConfig{
						CollectionInterval: 10 * time.Second,
						InitialDelay:       time.Second,
					},
					Driver:   "postgres",
					Host:     "localhost",
					Port:     5432,
					Database: "mydb",
					Username: "me",
					Password: "s3cr3t",
					Queries: []sqlquery.Query{
						{
							SQL: "select count(*) as count, type from mytable group by type",
							Metrics: []sqlquery.MetricCfg{
								{
									MetricName:       "val.count",
									ValueColumn:      "count",
									AttributeColumns: []string{"type"},
									Monotonic:        false,
									ValueType:        sqlquery.MetricValueTypeInt,
									DataType:         sqlquery.MetricTypeSum,
									Aggregation:      sqlquery.MetricAggregationCumulative,
									StaticAttributes: map[string]string{"foo": "bar"},
								},
							},
						},
					},
				},
			},
		},
		{
			fname:        "config-invalid-datatype.yaml",
			id:           component.NewIDWithName(metadata.Type, ""),
			errorMessage: "unsupported data_type: 'xyzgauge'",
		},
		{
			fname:        "config-invalid-valuetype.yaml",
			id:           component.NewIDWithName(metadata.Type, ""),
			errorMessage: "unsupported value_type: 'xyzint'",
		},
		{
			fname:        "config-invalid-aggregation.yaml",
			id:           component.NewIDWithName(metadata.Type, ""),
			errorMessage: "unsupported aggregation: 'xyzcumulative'",
		},
		{
			fname:        "config-invalid-missing-metricname.yaml",
			id:           component.NewIDWithName(metadata.Type, ""),
			errorMessage: "'metric_name' cannot be empty",
		},
		{
			fname:        "config-invalid-missing-valuecolumn.yaml",
			id:           component.NewIDWithName(metadata.Type, ""),
			errorMessage: "'value_column' cannot be empty",
		},
		{
			fname:        "config-invalid-missing-sql.yaml",
			id:           component.NewIDWithName(metadata.Type, ""),
			errorMessage: "'query.sql' cannot be empty",
		},
		{
			fname:        "config-invalid-missing-queries.yaml",
			id:           component.NewIDWithName(metadata.Type, ""),
			errorMessage: "'queries' cannot be empty",
		},
		{
			fname:        "config-invalid-missing-driver.yaml",
			id:           component.NewIDWithName(metadata.Type, ""),
			errorMessage: "'driver' cannot be empty",
		},
		{
			fname:        "config-invalid-missing-logs-metrics.yaml",
			id:           component.NewIDWithName(metadata.Type, ""),
			errorMessage: "at least one of 'query.logs' and 'query.metrics' must not be empty",
		},
		{
			fname:        "config-invalid-missing-datasource.yaml",
			id:           component.NewIDWithName(metadata.Type, ""),
			errorMessage: "'datasource' must be specified",
		},
		{
			fname: "config-invalid-missing-datasource-config-port-sql-server.yaml",
			id:    component.NewIDWithName(metadata.Type, ""),
			expected: &Config{
				Config: sqlquery.Config{
					ControllerConfig: scraperhelper.ControllerConfig{
						CollectionInterval: 10 * time.Second,
						InitialDelay:       time.Second,
					},
					Driver:   "sqlserver",
					Host:     "localhost",
					Database: "mydb",
					Username: "me",
					Password: "s3cr3t",
					Queries: []sqlquery.Query{
						{
							SQL: "select count(*) as count, type from mytable group by type",
							Metrics: []sqlquery.MetricCfg{
								{
									MetricName:       "val.count",
									ValueColumn:      "count",
									AttributeColumns: []string{"type"},
									Monotonic:        false,
									ValueType:        sqlquery.MetricValueTypeInt,
									DataType:         sqlquery.MetricTypeSum,
									Aggregation:      sqlquery.MetricAggregationCumulative,
									StaticAttributes: map[string]string{"foo": "bar"},
								},
							},
						},
					},
				},
			},
		},
		{
			fname:        "config-invalid-missing-datasource-config-port.yaml",
			id:           component.NewIDWithName(metadata.Type, ""),
			errorMessage: "'port' or 'datasource' must be specified",
		},
		{
			fname:        "config-invalid-missing-datasource-config-database.yaml",
			id:           component.NewIDWithName(metadata.Type, ""),
			errorMessage: "'database' or 'datasource' must be specified",
		},
		{
			fname: "config-logs.yaml",
			id:    component.NewIDWithName(metadata.Type, ""),
			expected: &Config{
				Config: sqlquery.Config{
					ControllerConfig: scraperhelper.ControllerConfig{
						CollectionInterval: 10 * time.Second,
						InitialDelay:       time.Second,
					},
					Driver:     "postgres",
					DataSource: "host=localhost port=5432 user=me password=s3cr3t sslmode=disable",
					Queries: []sqlquery.Query{
						{
							SQL:                "select * from test_logs where log_id > ?",
							TrackingColumn:     "log_id",
							TrackingStartValue: "10",
							Logs: []sqlquery.LogsCfg{
								{
									BodyColumn:       "log_body",
									AttributeColumns: []string{"log_attribute_1", "log_attribute_2"},
								},
							},
						},
					},
				},
			},
		},
		{
			fname: "config-logs-datasource-config.yaml",
			id:    component.NewIDWithName(metadata.Type, ""),
			expected: &Config{
				Config: sqlquery.Config{
					ControllerConfig: scraperhelper.ControllerConfig{
						CollectionInterval: 10 * time.Second,
						InitialDelay:       time.Second,
					},
					Driver:   "postgres",
					Host:     "localhost",
					Port:     5432,
					Database: "mydb",
					Username: "me",
					Password: "s3cr3t",
					Queries: []sqlquery.Query{
						{
							SQL:                "select * from test_logs where log_id > ?",
							TrackingColumn:     "log_id",
							TrackingStartValue: "10",
							Logs: []sqlquery.LogsCfg{
								{
									BodyColumn:       "log_body",
									AttributeColumns: []string{"log_attribute_1", "log_attribute_2"},
								},
							},
						},
					},
				},
			},
		},
		{
			fname:        "config-logs-missing-body-column.yaml",
			id:           component.NewIDWithName(metadata.Type, ""),
			errorMessage: "'body_column' must not be empty",
		},
		{
			fname:        "config-unnecessary-aggregation.yaml",
			id:           component.NewIDWithName(metadata.Type, ""),
			errorMessage: "aggregation=cumulative but data_type=gauge does not support aggregation",
		},
		{
			fname:        "config-both-datasource.yaml",
			id:           component.NewIDWithName(metadata.Type, ""),
			errorMessage: "'host' cannot be set when 'datasource' is specified",
		},
		{
			fname:        "config-invalid-driver.yaml",
			id:           component.NewIDWithName(metadata.Type, ""),
			errorMessage: "unsupported driver: invalid_driver",
		},
		{
			fname:        "config-invalid-driver-case-sensitive.yaml",
			id:           component.NewIDWithName(metadata.Type, ""),
			errorMessage: "unsupported driver: POSTGRES",
		},
		{
			fname:        "config-invalid-driver-partial-match.yaml",
			id:           component.NewIDWithName(metadata.Type, ""),
			errorMessage: "unsupported driver: postgre",
		},
	}

	for _, tt := range tests {
		t.Run(tt.fname, func(t *testing.T) {
			cm, err := confmaptest.LoadConf(filepath.Join("testdata", tt.fname))
			require.NoError(t, err)

			factory := NewFactory()
			cfg := factory.CreateDefaultConfig()

			sub, err := cm.Sub(tt.id.String())
			require.NoError(t, err)
			require.NoError(t, sub.Unmarshal(cfg))

			if tt.expected == nil {
				assert.ErrorContains(t, xconfmap.Validate(cfg), tt.errorMessage)
				return
			}
			assert.NoError(t, xconfmap.Validate(cfg))
			assert.Equal(t, tt.expected, cfg)
		})
	}
}

func TestCreateDefaultConfig(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	assert.Equal(t, 10*time.Second, cfg.CollectionInterval)
}

func TestConfig_Validate_Multierr(t *testing.T) {
	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config-invalid-multierr.yaml"))
	require.NoError(t, err)

	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()

	sub, err := cm.Sub(component.NewIDWithName(metadata.Type, "").String())
	require.NoError(t, err)
	require.NoError(t, sub.Unmarshal(cfg))

	err = xconfmap.Validate(cfg)

	assert.ErrorContains(t, err, "invalid metric config with metric_name 'my.metric'")
	assert.ErrorContains(t, err, "metric config has unsupported value_type: 'xint'")
	assert.ErrorContains(t, err, "metric config has unsupported data_type: 'xgauge'")
	assert.ErrorContains(t, err, "metric config has unsupported aggregation: 'xcumulative'")
}

func TestIsValidDriver(t *testing.T) {
	tests := []struct {
		name     string
		driver   string
		expected bool
	}{
		{"hdb driver", sqlquery.DriverHDB, true},
		{"mysql driver", sqlquery.DriverMySQL, true},
		{"oracle driver", sqlquery.DriverOracle, true},
		{"postgres driver", sqlquery.DriverPostgres, true},
		{"snowflake driver", sqlquery.DriverSnowflake, true},
		{"sqlserver driver", sqlquery.DriverSQLServer, true},
		{"tds driver", sqlquery.DriverTDS, true},
		{"invalid driver", "invalid", false},
		{"empty driver", "", false},
		{"case sensitive", "POSTGRES", false},
		{"partial match", "postgre", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sqlquery.IsValidDriver(tt.driver)
			assert.Equal(t, tt.expected, result)
		})
	}
}

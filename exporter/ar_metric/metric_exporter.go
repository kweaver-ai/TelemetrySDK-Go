package ar_metric

import (
	"bytes"
	"context"
	"encoding/json"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"

	"github.com/kweaver-ai/TelemetrySDK-Go/exporter/v2/common"
	"github.com/kweaver-ai/TelemetrySDK-Go/exporter/v2/public"
	"github.com/kweaver-ai/TelemetrySDK-Go/exporter/v2/resource"
	"github.com/kweaver-ai/TelemetrySDK-Go/exporter/v2/version"
)

// 跨包实现接口占位用。
var _ sdkmetric.Exporter = (*MetricExporter)(nil)

// MetricProvider 是一个全局变量，用于在业务代码中生产Meter。
var MetricProvider = (*sdkmetric.MeterProvider)(nil)

// Meter 是一个全局变量，用于在业务代码中生产Metric。
var Meter = otel.GetMeterProvider().Meter(
	version.MetricInstrumentationName,
	metric.WithInstrumentationVersion(version.TelemetrySDKVersion),
	metric.WithSchemaURL(version.MetricInstrumentationURL),
)

// MetricExporter 导出数据到AnyRobot Feed Ingester的 Metric 数据接收器。
type MetricExporter struct {
	*public.Exporter
}

// ExportMetrics 批量发送 AnyRobotMetrics 到AnyRobot Feed Ingester的 Metric 数据接收器。
func (e *MetricExporter) ExportMetrics(ctx context.Context, metrics []*metricdata.ResourceMetrics) error {
	if len(metrics) == 0 {
		return nil
	}
	arMetric := common.AnyRobotMetricsFromResourceMetrics(metrics)
	file := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "\t")
	if err := encoder.Encode(arMetric); err != nil {
		return err
	}
	return e.ExportData(ctx, file.Bytes())
}

// Temporality 计量方式，有累加值式和变化值式，默认使用累加值式。
func (e *MetricExporter) Temporality(k sdkmetric.InstrumentKind) metricdata.Temporality {
	return sdkmetric.DefaultTemporalitySelector(k)
}

// Aggregation 聚合类型，有7种，通过 metric.InstrumentKind 来区分。
func (e *MetricExporter) Aggregation(k sdkmetric.InstrumentKind) sdkmetric.Aggregation {
	return sdkmetric.DefaultAggregationSelector(k)
}

// Export 导出数据方法，没做缓存队列，每一条产生的Metric立即发送。
func (e *MetricExporter) Export(ctx context.Context, data *metricdata.ResourceMetrics) error {
	return e.ExportMetrics(ctx, []*metricdata.ResourceMetrics{data})
}

// ForceFlush 强制发送缓存队列中的数据，因为没做缓存队列，所以是空的。
func (e *MetricExporter) ForceFlush(ctx context.Context) error {
	return ctx.Err()
}

// NewExporter 创建已启动的 MetricExporter 。
func NewExporter(c public.Client) *MetricExporter {
	return &MetricExporter{
		public.NewExporter(c),
	}
}

// MetricResource 传入 Metric 的默认resource。
func MetricResource() *sdkresource.Resource {
	return resource.MetricResource()
}

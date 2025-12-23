package common

import (
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// Link 自定义 Link 统一Attribute。
type Link struct {
	SpanContext           trace.SpanContext `json:"SpanContext"`
	Attributes            []*Attribute      `json:"Attributes"`
	DroppedAttributeCount int               `json:"DroppedAttributeCount"`
}

// AnyRobotLinkFromLink 单条sdktrace.Link转换为*Link。
func AnyRobotLinkFromLink(link sdktrace.Link) *Link {
	return &Link{
		SpanContext:           link.SpanContext,
		Attributes:            AnyRobotAttributesFromKeyValues(link.Attributes),
		DroppedAttributeCount: link.DroppedAttributeCount,
	}
}

// AnyRobotLinksFromLinks 批量sdktrace.Link转换为[]*Link。
func AnyRobotLinksFromLinks(links []sdktrace.Link) []*Link {
	if links == nil {
		return make([]*Link, 0)
	}
	arLinks := make([]*Link, 0, len(links))
	for i := 0; i < len(links); i++ {
		arLinks = append(arLinks, AnyRobotLinkFromLink(links[i]))
	}
	return arLinks
}

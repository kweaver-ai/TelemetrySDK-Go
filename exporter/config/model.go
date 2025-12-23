package config

type ExportersTyp string
type ExportersSubTyp string

const (
	File    ExportersTyp = "file"
	Console ExportersTyp = "console"
	Http    ExportersTyp = "http"
	Mq      ExportersTyp = "protonmq"
)

const (
	ProtonMqNsq      ExportersSubTyp = "nsq"
	ProtonMqTonglink ExportersSubTyp = "tonglink"
	ProtonMqBmq      ExportersSubTyp = "bmq"
	ProtonMqKafka    ExportersSubTyp = "kafka"
	ProtonMqHtp20    ExportersSubTyp = "htp20"
)

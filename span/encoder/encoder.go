package encoder

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"strconv"
	"time"
	"unsafe"

	"github.com/kweaver-ai/TelemetrySDK-Go/span/v2/exporter"
	"github.com/kweaver-ai/TelemetrySDK-Go/span/v2/field"
)

var (
	maxWaitExporterTime = 20 * time.Second

	_quotationSafe   = []byte("\\\"")
	_reverseSafe     = []byte("\\\\")
	_backspaceSafe   = []byte("\\b")
	_formfeedSafe    = []byte("\\f")
	_horizontalSafe  = []byte("\\t")
	_lineFeedSafe    = []byte("\\n")
	_carriageSafe    = []byte("\\r")
	_quotation       = []byte("\"")
	_leftBracket     = []byte("[")
	_rightBracket    = []byte("]")
	_seperator       = []byte(",")
	_lineFeed        = []byte("\n")
	_leftBigBracket  = []byte("{")
	_rightBigBracket = []byte("}")
	_colon           = []byte(": ")
)

type Encoder interface {
	Write(field.Field) error
	Close() error
}

// JsonEncoder JSON编码器。
type JsonEncoder struct {
	w            io.Writer
	buf          io.Writer
	bufReal      *bytes.Buffer
	End          []byte
	logExporters map[string]exporter.LogExporter
	ctx          context.Context
	cancelFunc   context.CancelFunc
}

// NewJsonEncoder 创建2.5.0及之前的版本的JSON编码器。
func NewJsonEncoder(w io.Writer) Encoder {
	ctx, cancel := context.WithCancel(context.Background())
	encoder := &JsonEncoder{
		w:            w,
		buf:          nil,
		bufReal:      bytes.NewBuffer(make([]byte, 0, 4096)),
		End:          _lineFeed,
		logExporters: nil,
		ctx:          ctx,
		cancelFunc:   cancel,
	}
	encoder.buf = encoder.bufReal
	return encoder
}

// NewJsonEncoderBench 创建2.5.0及之前的版本的JSON编码器。
func NewJsonEncoderBench(w io.Writer) Encoder {
	ctx, cancel := context.WithCancel(context.Background())
	encoder := &JsonEncoder{
		w:            w,
		buf:          io.Discard,
		bufReal:      bytes.NewBuffer(make([]byte, 0, 4096)),
		End:          _lineFeed,
		logExporters: nil,
		ctx:          ctx,
		cancelFunc:   cancel,
	}
	return encoder
}

// NewJsonEncoderWithExporters 创建2.6.0后的异步发送模式JSON编码器。
func NewJsonEncoderWithExporters(exporters ...exporter.LogExporter) Encoder {
	ctx, cancel := context.WithCancel(context.Background())
	eps := make(map[string]exporter.LogExporter)
	for _, e := range exporters {
		eps[e.Name()] = e
	}
	encoder := &JsonEncoder{
		w:            nil,
		buf:          nil,
		bufReal:      bytes.NewBuffer(make([]byte, 0, 4096)),
		End:          _lineFeed,
		logExporters: eps,
		ctx:          ctx,
		cancelFunc:   cancel,
	}
	encoder.buf = encoder.bufReal
	return encoder
}

// Write 区分本地输出和上报AnyRobot，本地输出仅提取关键信息输出，上报AnyRobot为数组形式，有返回值error判断是否发送成功。
func (js *JsonEncoder) Write(f field.Field) error {
	//判断用户是否使用原来2.5.0之前的io.Writer实现输出，如果是的话由于批量发送，为保证与原来数据模型一致将数组遍历单个输出
	if js.w != nil {
		//将数组遍历使用io.Writer单个输出
		js.dealIoWriter(f)
		return nil
	}
	//判断用户是否使用2.6.0之后的exporter实现输出
	if js.logExporters != nil && len(js.logExporters) != 0 {
		//判断用户是否使用新增的RealTimeExporter实现标准输出，如果是的话由于批量发送，为保证与原来数据模型一致将数组遍历单个输出
		stdoutExporter, ok := js.logExporters["RealTimeExporter"]
		if ok {
			//将数组遍历使用RealTimeExporter单个输出，并提取关键信息输出。
			js.dealRealTimeExporter(stdoutExporter, f)
		}
		//不是RealTimeExporter的剩余exporter比如arExporter将输出整个数组，下面将整个数组转为byte
		err := js.write(f)
		if err != nil {
			log.Println(field.GenerateSpecificError(err))
		}
		_, writeBytesErr := js.WriteBytes(js.End)
		if writeBytesErr != nil {
			log.Println(field.GenerateSpecificError(writeBytesErr))
		}
		//调用不是RealTimeExporter的剩余exporter比如arExporter的输出方法将整个数组进行输出
		flushWithExportersErr := js.flushWithExporters()
		if flushWithExportersErr != nil {
			log.Println(field.GenerateSpecificError(flushWithExportersErr))
			return flushWithExportersErr
		}
	}
	return nil
}

// Close 关闭JSON编码器。
func (js *JsonEncoder) Close() error {
	if js.bufReal.Len() > 0 {
		_ = js.flush()
		_ = js.flushWithExporters()
		return nil
	}
	go func() {
		t := time.NewTimer(maxWaitExporterTime)
		defer t.Stop()
		<-t.C
		js.cancelFunc()
		if js.logExporters != nil && len(js.logExporters) != 0 {
			for _, exporter_ := range js.logExporters {
				_ = exporter_.Shutdown(js.ctx)
			}
		}
	}()
	return nil
}

func (js *JsonEncoder) dealIoWriter(f field.Field) {
	//断言为数组形式
	fieldArr, ok := f.(*field.ArrayField)
	if ok {
		for i := 0; i < fieldArr.Length(); i++ {
			oneField, errArr := fieldArr.At(i)
			if errArr != nil {
				log.Println(field.GenerateSpecificError(errArr))
			}
			//将单个log转为byte数组
			err := js.write(oneField)
			if err != nil {
				log.Println(field.GenerateSpecificError(err))
			}
			_, writeBytesErr := js.WriteBytes(js.End)
			if writeBytesErr != nil {
				log.Println(field.GenerateSpecificError(writeBytesErr))
			}
			//输出
			_ = js.flush()
		}
	}
}

// dealRealTimeExporter 处理控制台输出的特殊逻辑，控制台输出要求单条日志按时间顺序输出，并且不换行，只展示关键信息。
func (js *JsonEncoder) dealRealTimeExporter(realTimeExporter exporter.LogExporter, logs field.Field) {
	//断言为数组形式
	fieldArr, fieldArrOk := logs.(*field.ArrayField)
	if fieldArrOk {
		for i := 0; i < fieldArr.Length(); i++ {
			// 取出单个Log。
			oneField, errArr := fieldArr.At(i)
			if errArr != nil {
				log.Println(field.GenerateSpecificError(errArr))
			}
			//将单个log转为byte数组
			err := js.write(oneField)
			if err != nil {
				log.Println(field.GenerateSpecificError(err))
			}
			_, writeBytesErr := js.WriteBytes(js.End)
			if writeBytesErr != nil {
				log.Println(field.GenerateSpecificError(writeBytesErr))
			}
			//调用 realTimeExporter 的输出方法将单个log输出
			exportLogsErr := realTimeExporter.ExportLogs(js.ctx, js.bufReal.Bytes())
			if exportLogsErr != nil {
				log.Println(field.GenerateSpecificError(exportLogsErr))
			}
			js.bufReal.Reset()
		}
	}
}

func (js *JsonEncoder) flush() error {
	if js.w != nil {
		_, err := js.w.Write(js.bufReal.Bytes())
		if err != nil {
			log.Println(field.GenerateSpecificError(err))
		}
		js.bufReal.Reset()
	}
	return nil
}

func (js *JsonEncoder) flushWithExporters() error {
	var returnErr error = nil
	if js.logExporters != nil && len(js.logExporters) != 0 {
		for _, e := range js.logExporters {
			//过滤掉已经输出的RealTimeExporter，其他exporter正常输出
			if e.Name() == "RealTimeExporter" {
				continue
			}
			if err := e.ExportLogs(js.ctx, js.bufReal.Bytes()); err != nil {
				// 如果错误则记日志。
				log.Println(field.GenerateSpecificError(err))
				returnErr = err
			}
		}
		js.bufReal.Reset()
	}
	return returnErr
}

func (js *JsonEncoder) write(f field.Field) error {
	switch f.Type() {
	default:
		return nil
	case field.IntType:
		v := strconv.Itoa(int(f.(field.IntField)))
		bytes_ := js.string2Bytes(v)
		_, err := js.WriteBytes(bytes_)
		return err
	case field.Float64Type:
		v := strconv.FormatFloat(float64(f.(field.Float64Field)), 'f', -1, 64)
		bytes_ := js.string2Bytes(v)
		_, err := js.WriteBytes(bytes_)
		return err
	case field.StringType:
		v := string(f.(field.StringField))
		_, _ = js.WriteBytes(_quotation)
		_, _ = js.safeWriteString(v)
		_, err := js.WriteBytes(_quotation)
		return err
	case field.TimeType:
		v := strconv.FormatInt(time.Time(f.(field.TimeField)).UnixNano(), 10)
		bytes_ := js.string2Bytes(v)
		_, err := js.WriteBytes(bytes_)
		return err
	case field.ArrayType:
		v := f.(*field.ArrayField)
		_, _ = js.WriteBytes(_leftBracket)
		i := 0
		for ; i < len(*v)-1; i += 1 {
			_ = js.write((*v)[i])
			_, _ = js.WriteBytes(_seperator)
		}
		if i < len(*v) {
			_ = js.write((*v)[i])
		}
		_, err := js.WriteBytes(_rightBracket)
		return err
	case field.StructType:
		fs := f.(*field.StructField)
		_, _ = js.WriteBytes(_leftBigBracket)
		i := 0
		for ; i < fs.Length()-1; i += 1 {
			k, v, _ := fs.At(i)
			_, _ = js.WriteBytes(_quotation)
			_, _ = js.safeWriteString(k)
			_, _ = js.WriteBytes(_quotation)
			_, _ = js.WriteBytes(_colon)
			_ = js.write(v)
			_, _ = js.WriteBytes(_seperator)
		}
		if k, v, err := fs.At(i); err == nil {
			_, _ = js.WriteBytes(_quotation)
			_, _ = js.safeWriteString(k)
			_, _ = js.WriteBytes(_quotation)
			_, _ = js.WriteBytes(_colon)
			_ = js.write(v)
		}
		_, err := js.WriteBytes(_rightBigBracket)
		return err

	case field.JsonType:
		j := f.(*field.JsonFiled)
		b, err := json.Marshal(j.Data)

		if err != nil {
			return err
		}
		_, err = js.WriteBytes(b)
		return err

	case field.MapType:
		b, err := json.Marshal(f)

		if err != nil {
			return err
		}
		_, err = js.WriteBytes(b)
		return err
	}

}

// String2Bytes unsafe convert string to []byte, they point to the same memory
func (js *JsonEncoder) string2Bytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func (js *JsonEncoder) safeWriteString(s string) (int, error) {
	left := 0
	bytes_ := js.string2Bytes(s)
	var err error
	for k, i := range s {
		switch i {
		default:
		case '"':
			_, _ = js.WriteBytes(bytes_[left:k])
			left = k + 1
			_, err = js.WriteBytes(_quotationSafe)
		case '\\':
			_, _ = js.WriteBytes(bytes_[left:k])
			left = k + 1
			_, err = js.WriteBytes(_reverseSafe)
		case '\b':
			_, _ = js.WriteBytes(bytes_[left:k])
			left = k + 1
			_, err = js.WriteBytes(_backspaceSafe)
		case '\f':
			_, _ = js.WriteBytes(bytes_[left:k])
			left = k + 1
			_, err = js.WriteBytes(_formfeedSafe)
		case '\t':
			_, _ = js.WriteBytes(bytes_[left:k])
			left = k + 1
			_, err = js.WriteBytes(_horizontalSafe)
		case '\n':
			_, _ = js.WriteBytes(bytes_[left:k])
			left = k + 1
			_, err = js.WriteBytes(_lineFeedSafe)
		case '\r':
			_, _ = js.WriteBytes(bytes_[left:k])
			left = k + 1
			_, err = js.WriteBytes(_carriageSafe)
		}
		if err != nil {
			return left, err
		}
	}

	if left < len(bytes_) {
		_, err = js.WriteBytes(bytes_[left:])
	}
	return left, err
}

func (js *JsonEncoder) WriteBytes(c []byte) (int, error) {
	return js.buf.Write(c)
}

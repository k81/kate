package simple

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/modern-go/gls"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var bufferpool = buffer.NewPool()
var pid = os.Getpid()

// For escaping; see encoder.safeAddString below.
const _hex = "0123456789abcdef"

var encoderPool = sync.Pool{New: func() interface{} {
	return &encoder{}
}}

func getEncoder() *encoder {
	return encoderPool.Get().(*encoder)
}

func putJSONEncoder(enc *encoder) {
	if enc.reflectBuf != nil {
		enc.reflectBuf.Free()
	}

	enc.buf = nil
	enc.openNamespaces = 0
	enc.reflectBuf = nil
	enc.reflectEnc = nil
	encoderPool.Put(enc)
}

type encoder struct {
	buf            *buffer.Buffer
	openNamespaces int

	// for encoding generic values by reflection
	reflectBuf *buffer.Buffer
	reflectEnc *json.Encoder
}

// NewEncoder creates a fast, low-allocation encoder. The encoder
// appropriately escapes all field keys and values.
//
// Note that the encoder doesn't deduplicate keys, so it's possible to produce
// a message like
//   {"foo":"bar","foo":"baz"}
// This is permitted by the JSON specification, but not encouraged. Many
// libraries will ignore duplicate key-value pairs (typically keeping the last
// pair) when unmarshaling, but users should attempt to avoid adding duplicate
// keys.
func NewEncoder() zapcore.Encoder {
	return &encoder{
		buf: bufferpool.Get(),
	}
}

func (enc *encoder) AddArray(key string, arr zapcore.ArrayMarshaler) error {
	enc.addKey(key)
	return enc.AppendArray(arr)
}

func (enc *encoder) AddObject(key string, obj zapcore.ObjectMarshaler) error {
	enc.addKey(key)
	return enc.AppendObject(obj)
}

func (enc *encoder) AddBinary(key string, val []byte) {
	enc.AddString(key, base64.StdEncoding.EncodeToString(val))
}

func (enc *encoder) AddByteString(key string, val []byte) {
	enc.addKey(key)
	enc.AppendByteString(val)
}

func (enc *encoder) AddBool(key string, val bool) {
	enc.addKey(key)
	enc.AppendBool(val)
}

func (enc *encoder) AddComplex128(key string, val complex128) {
	enc.addKey(key)
	enc.AppendComplex128(val)
}

func (enc *encoder) AddDuration(key string, val time.Duration) {
	enc.addKey(key)
	enc.AppendDuration(val)
}

func (enc *encoder) AddFloat64(key string, val float64) {
	enc.addKey(key)
	enc.AppendFloat64(val)
}

func (enc *encoder) AddInt64(key string, val int64) {
	enc.addKey(key)
	enc.AppendInt64(val)
}

func (enc *encoder) resetReflectBuf() {
	if enc.reflectBuf == nil {
		enc.reflectBuf = bufferpool.Get()
		enc.reflectEnc = json.NewEncoder(enc.reflectBuf)

		// For consistency with our custom JSON encoder.
		enc.reflectEnc.SetEscapeHTML(false)
	} else {
		enc.reflectBuf.Reset()
	}
}

var nullLiteralBytes = []byte("null")

// Only invoke the standard JSON encoder if there is actually something to
// encode; otherwise write JSON null literal directly.
func (enc *encoder) encodeReflected(obj interface{}) ([]byte, error) {
	if obj == nil {
		return nullLiteralBytes, nil
	}

	enc.resetReflectBuf()

	if err := enc.reflectEnc.Encode(obj); err != nil {
		return nil, err
	}

	enc.reflectBuf.TrimNewline()

	return enc.reflectBuf.Bytes(), nil
}

func (enc *encoder) AddReflected(key string, obj interface{}) error {
	valueBytes, err := enc.encodeReflected(obj)
	if err != nil {
		return err
	}

	enc.addKey(key)
	_, err = enc.buf.Write(valueBytes)

	return err
}

func (enc *encoder) OpenNamespace(key string) {
	enc.addKey(key)
	enc.buf.AppendByte('{')
	enc.openNamespaces++
}

func (enc *encoder) AddString(key, val string) {
	enc.addKey(key)
	enc.AppendString(val)
}

func (enc *encoder) AddTime(key string, val time.Time) {
	enc.addKey(key)
	enc.AppendTime(val)
}

func (enc *encoder) AddUint64(key string, val uint64) {
	enc.addKey(key)
	enc.AppendUint64(val)
}

func (enc *encoder) AppendArray(arr zapcore.ArrayMarshaler) error {
	enc.addElementSeparator()
	enc.buf.AppendByte('[')
	err := arr.MarshalLogArray(enc)
	enc.buf.AppendByte(']')

	return err
}

func (enc *encoder) AppendObject(obj zapcore.ObjectMarshaler) error {
	enc.addElementSeparator()
	enc.buf.AppendByte('{')
	err := obj.MarshalLogObject(enc)
	enc.buf.AppendByte('}')

	return err
}

func (enc *encoder) AppendBool(val bool) {
	enc.addElementSeparator()
	enc.buf.AppendBool(val)
}

func (enc *encoder) AppendByteString(val []byte) {
	enc.addElementSeparator()
	enc.buf.AppendByte('"')
	enc.safeAddByteString(val)
	enc.buf.AppendByte('"')
}

func (enc *encoder) AppendComplex128(val complex128) {
	enc.addElementSeparator()
	// Cast to a platform-independent, fixed-size type.
	r, i := real(val), imag(val)

	enc.buf.AppendByte('"')
	// Because we're always in a quoted string, we can use strconv without
	// special-casing NaN and +/-Inf.
	enc.buf.AppendFloat(r, 64)
	enc.buf.AppendByte('+')
	enc.buf.AppendFloat(i, 64)
	enc.buf.AppendByte('i')
	enc.buf.AppendByte('"')
}

func (enc *encoder) AppendDuration(val time.Duration) {
	cur := enc.buf.Len()
	zapcore.NanosDurationEncoder(val, enc)

	if cur == enc.buf.Len() {
		// User-supplied EncodeDuration is a no-op. Fall back to nanoseconds to keep
		// JSON valid.
		enc.AppendInt64(int64(val))
	}
}

func (enc *encoder) AppendInt64(val int64) {
	enc.addElementSeparator()
	enc.buf.AppendInt(val)
}

func (enc *encoder) AppendReflected(val interface{}) error {
	valueBytes, err := enc.encodeReflected(val)
	if err != nil {
		return err
	}

	enc.addElementSeparator()
	_, err = enc.buf.Write(valueBytes)

	return err
}

func (enc *encoder) AppendString(val string) {
	enc.addElementSeparator()
	enc.buf.AppendByte('"')
	enc.safeAddString(val)
	enc.buf.AppendByte('"')
}

func (enc *encoder) AppendTime(val time.Time) {
	cur := enc.buf.Len()
	zapcore.ISO8601TimeEncoder(val, enc)

	if cur == enc.buf.Len() {
		// User-supplied EncodeTime is a no-op. Fall back to nanos since epoch to keep
		// output JSON valid.
		enc.AppendInt64(val.UnixNano())
	}
}

func (enc *encoder) AppendUint64(val uint64) {
	enc.addElementSeparator()
	enc.buf.AppendUint(val)
}

func (enc *encoder) AddComplex64(k string, v complex64) { enc.AddComplex128(k, complex128(v)) }
func (enc *encoder) AddFloat32(k string, v float32)     { enc.AddFloat64(k, float64(v)) }
func (enc *encoder) AddInt(k string, v int)             { enc.AddInt64(k, int64(v)) }
func (enc *encoder) AddInt32(k string, v int32)         { enc.AddInt64(k, int64(v)) }
func (enc *encoder) AddInt16(k string, v int16)         { enc.AddInt64(k, int64(v)) }
func (enc *encoder) AddInt8(k string, v int8)           { enc.AddInt64(k, int64(v)) }
func (enc *encoder) AddUint(k string, v uint)           { enc.AddUint64(k, uint64(v)) }
func (enc *encoder) AddUint32(k string, v uint32)       { enc.AddUint64(k, uint64(v)) }
func (enc *encoder) AddUint16(k string, v uint16)       { enc.AddUint64(k, uint64(v)) }
func (enc *encoder) AddUint8(k string, v uint8)         { enc.AddUint64(k, uint64(v)) }
func (enc *encoder) AddUintptr(k string, v uintptr)     { enc.AddUint64(k, uint64(v)) }
func (enc *encoder) AppendComplex64(v complex64)        { enc.AppendComplex128(complex128(v)) }
func (enc *encoder) AppendFloat64(v float64)            { enc.appendFloat(v, 64) }
func (enc *encoder) AppendFloat32(v float32)            { enc.appendFloat(float64(v), 32) }
func (enc *encoder) AppendInt(v int)                    { enc.AppendInt64(int64(v)) }
func (enc *encoder) AppendInt32(v int32)                { enc.AppendInt64(int64(v)) }
func (enc *encoder) AppendInt16(v int16)                { enc.AppendInt64(int64(v)) }
func (enc *encoder) AppendInt8(v int8)                  { enc.AppendInt64(int64(v)) }
func (enc *encoder) AppendUint(v uint)                  { enc.AppendUint64(uint64(v)) }
func (enc *encoder) AppendUint32(v uint32)              { enc.AppendUint64(uint64(v)) }
func (enc *encoder) AppendUint16(v uint16)              { enc.AppendUint64(uint64(v)) }
func (enc *encoder) AppendUint8(v uint8)                { enc.AppendUint64(uint64(v)) }
func (enc *encoder) AppendUintptr(v uintptr)            { enc.AppendUint64(uint64(v)) }

func (enc *encoder) Clone() zapcore.Encoder {
	clone := enc.clone()
	// nolint:errcheck
	clone.buf.Write(enc.buf.Bytes())

	return clone
}

func (enc *encoder) clone() *encoder {
	clone := getEncoder()
	clone.openNamespaces = enc.openNamespaces
	clone.buf = bufferpool.Get()

	return clone
}

func (enc *encoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	final := enc.clone()

	final.buf.AppendString(fmt.Sprintf("%-6v ", ent.Level.CapitalString()))
	final.buf.AppendString(ent.Time.Format("2006-01-02T15:04:05.000Z0700"))
	final.buf.AppendString(" [")
	final.buf.AppendInt(int64(pid))
	final.buf.AppendString("][")
	final.buf.AppendString(fmt.Sprintf("%05d", gls.GoID()))
	final.buf.AppendString("]\t")

	final.buf.AppendString(ent.Message)
	final.buf.AppendString("\t")

	final.buf.AppendByte('{')
	if enc.buf.Len() > 0 || len(fields) > 0 {
		if enc.buf.Len() > 0 {
			final.addElementSeparator()
			// nolint:errcheck
			final.buf.Write(enc.buf.Bytes())
		}

		final.addFields(fields)
		final.closeOpenNamespaces()
	}

	final.buf.AppendByte('}')

	final.buf.AppendString("\t[")
	final.buf.AppendString(ent.Caller.TrimmedPath())
	final.buf.AppendByte(']')
	final.buf.AppendString(zapcore.DefaultLineEnding)

	ret := final.buf
	putJSONEncoder(final)

	return ret, nil
}

func (enc *encoder) truncate() {
	enc.buf.Reset()
}

func (enc *encoder) closeOpenNamespaces() {
	for i := 0; i < enc.openNamespaces; i++ {
		enc.buf.AppendByte('}')
	}
}

func (enc *encoder) addKey(key string) {
	enc.addElementSeparator()
	enc.buf.AppendByte('"')
	enc.safeAddString(key)
	enc.buf.AppendByte('"')
	enc.buf.AppendByte(':')
}

func (enc *encoder) addElementSeparator() {
	last := enc.buf.Len() - 1
	if last < 0 {
		return
	}

	switch enc.buf.Bytes()[last] {
	case '{', '[', ':', ',', ' ':
		return
	default:
		enc.buf.AppendByte(',')
	}
}

func (enc *encoder) appendFloat(val float64, bitSize int) {
	enc.addElementSeparator()

	switch {
	case math.IsNaN(val):
		enc.buf.AppendString(`"NaN"`)
	case math.IsInf(val, 1):
		enc.buf.AppendString(`"+Inf"`)
	case math.IsInf(val, -1):
		enc.buf.AppendString(`"-Inf"`)
	default:
		enc.buf.AppendFloat(val, bitSize)
	}
}

// safeAddString JSON-escapes a string and appends it to the internal buffer.
// Unlike the standard library's encoder, it doesn't attempt to protect the
// user from browser vulnerabilities or JSONP-related problems.
func (enc *encoder) safeAddString(s string) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}

		r, size := utf8.DecodeRuneInString(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}

		enc.buf.AppendString(s[i : i+size])
		i += size
	}
}

// safeAddByteString is no-alloc equivalent of safeAddString(string(s)) for s []byte.
func (enc *encoder) safeAddByteString(s []byte) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}

		r, size := utf8.DecodeRune(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}

		// nolint:errcheck
		enc.buf.Write(s[i : i+size])
		i += size
	}
}

// tryAddRuneSelf appends b if it is valid UTF-8 character represented in a single byte.
func (enc *encoder) tryAddRuneSelf(b byte) bool {
	if b >= utf8.RuneSelf {
		return false
	}

	if 0x20 <= b && b != '\\' && b != '"' {
		enc.buf.AppendByte(b)
		return true
	}

	switch b {
	case '\\', '"':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte(b)
	case '\n':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('n')
	case '\r':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('r')
	case '\t':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('t')
	default:
		// Encode bytes < 0x20, except for the escape sequences above.
		enc.buf.AppendString(`\u00`)
		enc.buf.AppendByte(_hex[b>>4])
		enc.buf.AppendByte(_hex[b&0xF])
	}

	return true
}

func (enc *encoder) tryAddRuneError(r rune, size int) bool {
	if r == utf8.RuneError && size == 1 {
		enc.buf.AppendString(`\ufffd`)
		return true
	}

	return false
}

func (enc *encoder) addFields(fields []zapcore.Field) {
	for i := range fields {
		fields[i].AddTo(enc)
	}
}

func init() {
	// nolint: errcheck
	zap.RegisterEncoder("simple", func(config zapcore.EncoderConfig) (encoder zapcore.Encoder, e error) {
		return NewEncoder(), nil
	})
}

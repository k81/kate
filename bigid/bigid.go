package bigid

import (
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// 0-9 mod 10bits
// 10-17 vsid 8bits
// 18-57 timestamp 40bits
// 58-61 reserve 4bits
// 62-63 version 2bits
var (
	autoIncSeq uint64
	baseTime   time.Time
)

func init() {
	baseTime, _ = time.ParseInLocation("2006-01-02 15:04:05", "2015-06-06 00:00:00", time.Local)
}

type BigID int64

func New(vsId uint64) (bigId uint64) {
	seq := atomic.AddUint64(&autoIncSeq, 1)
	timestamp := uint64(time.Since(baseTime) / time.Millisecond)
	bigId = 1 << 4                                 // version + reserved
	bigId = bigId<<40 | (timestamp & 0xFFFFFFFFFF) // timestamp
	bigId = bigId<<8 | (vsId & 0xFF)               // vsId
	bigId = bigId<<10 | (seq & 0x3FF)              // auto increase seq
	return
}

// NewFromString parse bigid from 36B string
func NewFromString(str string) (BigID, error) {
	id, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return BigID(0), err
	}
	return BigID(id), nil
}

func Fake(vsId uint64) (bigId uint64) {
	bigId = (vsId & 0xFF) << 10
	return
}

func GetVSId(bigId uint64) uint64 {
	return (bigId >> 10) & 0xFF
}

// String return the string representation of BigID
func (bigID BigID) String() string {
	return strings.ToUpper(strconv.FormatInt(int64(bigID), 10))
}

// UnmarshalBind implements the BindUnmarshaller.
// see utils.BindUnmarshaller
func (bigID *BigID) UnmarshalBind(value string) error {
	id, err := strconv.ParseInt(value, 10, 64)
	*bigID = BigID(id)
	return err
}

// UnmarshalJSON implements json.Unmarshaller
func (bigID *BigID) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), "\"")
	id, err := strconv.ParseInt(str, 10, 64)
	*bigID = BigID(id)
	return err
}

// MarshalJSON implements json.Marshaller
func (bigID BigID) MarshalJSON() ([]byte, error) {
	data := []byte(fmt.Sprintf("\"%v\"", bigID.String()))
	return data, nil
}

type BigIDStruct struct {
	Version    uint64 `json:"version"`
	Reserved   uint64 `json:"reserved"`
	Timestamp  uint64 `json:"timestamp"`
	VSId       uint64 `json:"vsid"`
	AutoIncSeq uint64 `json:"auto_inc_seq"`
	CreateTime string `json:"create_time"`
}

func Parse(bigId uint64) *BigIDStruct {
	n := &BigIDStruct{}

	n.AutoIncSeq = bigId & 0x3FF
	bigId = bigId >> 10

	n.VSId = bigId & 0xFF
	bigId = bigId >> 8

	n.Timestamp = bigId & 0xFFFFFFFFFF
	bigId = bigId >> 40

	n.Reserved = bigId & 0xF
	bigId = bigId >> 4

	n.Version = bigId & 0x3

	d := time.Duration(n.Timestamp) * time.Millisecond
	n.CreateTime = baseTime.Add(d).Format(time.RFC3339Nano)

	return n
}

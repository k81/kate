package decimal

import (
	"os"
	"testing"

	"github.com/k81/kate/bigid"
	"github.com/k81/kate/orm"
)

type DecimalTest struct {
	Id    uint64 `orm:"pk"`
	Value Decimal
}

func (t *DecimalTest) TableName() string {
	return "test_decimal"
}

func TestDecimal(t *testing.T) {
	id := bigid.New(2)
	obj := &DecimalTest{
		Id:    id,
		Value: NewFromFloat(18.70),
	}

	if _, err := orm.NewOrm().Insert(obj); err != nil {
		t.Fatalf("orm insert failed: %v (bigid=%v, vsid=%v)", err, id, bigid.GetVSId(id))
	}

	obj2 := &DecimalTest{
		Id:    id,
		Value: Zero,
	}

	if err := orm.NewOrm().Read(obj2); err != nil {
		t.Fatalf("orm read failed: %v (bigid=%v, vsid=%v)", err, id, bigid.GetVSId(id))
	}
	t.Logf("DecimalTest.value=%v", obj2.Value)

	if !obj.Value.Equal(obj2.Value) {
		t.Fatalf("insert and read not equals: inserted=%v, got=%v", obj.Value, obj2.Value)
	}
}

func TestMain(m *testing.M) {
	orm.RegisterDataBase("default", "cdbpool", "tcp(192.168.0.132:9123,192.168.0.133:9123)/default?timeout=5s&readTimeout=15s&writeTimeout=15s", 20, 100)
	orm.RegisterModel("test", new(DecimalTest))
	os.Exit(m.Run())
}

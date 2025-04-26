package asana

import (
    "encoding/json"
    "testing"
)

func TestCustomFieldBase_Precision_ParseZero(t *testing.T) {
    cf := &CustomFieldBase{}
    if err := json.Unmarshal([]byte(`
{
	"precision": 0
}
`), cf); err != nil {
        t.Fatal(err)
    }

    if cf.Precision == nil || *cf.Precision != 0 {
        t.Errorf("Expected Precision to be a pointer to the integer zero, but saw %v", cf.Precision)
    }
}

func TestCustomFieldBase_Precision_ParseMissing(t *testing.T) {
    cf := &CustomFieldBase{}
    if err := json.Unmarshal([]byte(`
{
	"name": "name"
}
`), cf); err != nil {
        t.Fatal(err)
    }

    if cf.Precision != nil {
        t.Errorf("Expected Precision to be a nil, but saw %v", cf.Precision)
    }
}

func TestCustomFieldBase_Precision_SerializeZero(t *testing.T) {
    val := 0
    cf := &CustomFieldBase{Precision: &val}
    if bs, err := json.Marshal(cf); err != nil {
        t.Fatal(err)
    } else {
        if string(bs) != `{"precision":0,"resource_subtype":""}` {
            t.Errorf("Expected Precision to be a zero, but saw %v", string(bs))
        }
    }

}

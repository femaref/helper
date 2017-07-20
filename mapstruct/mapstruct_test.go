package mapstruct

import (
    "testing"
)

type t1 struct {
    A int
    B string

    Array []int
    StructArray []t2
    Nested t2
}

type t2 struct {
    A int
    B string
}

var maptests = []struct{
    in map[string]interface{}
    target t1
    expt t1
    out error
} {
    { map[string]interface{}{"a":2, "b":"foo"}, t1{}, t1{A:2, B:"foo"}, nil},
    { map[string]interface{}{"a":2, "b":"foo", "struct_array":[]map[string]interface{}{map[string]interface{}{"a":0}, map[string]interface{}{"a":1}}}, t1{}, t1{A:2, B:"foo", StructArray:[]t2{t2{A:0}, t2{A:1}}}, nil},
    { map[string]interface{}{"a":2, "b":"foo", "nested":map[string]interface{}{"a":3, "b":"baz"}}, t1{}, t1{A:2, B:"foo", Nested:t2{3, "baz"}}, nil},
    { map[string]interface{}{"a":2, "b":"foo", "array":[]int{1,2,3}}, t1{}, t1{A:2, B:"foo", Array:[]int{1,2,3}}, nil},
}

func TestMap(t *testing.T) {
	for _, tt := range maptests {
	    tt := tt
		err := MapToStructv2(&tt.target, tt.in)

		var is_error bool
		if tt.out != nil {
			if err == nil || err.Error() != tt.out.Error() {
				is_error = true
			}
		} else {
			if err != tt.out {
				is_error = true
			}
		}

		if !is_error {
		    if tt.target.A != tt.expt.A || tt.target.B != tt.expt.B || tt.target.Nested != tt.expt.Nested  {
		        is_error = true
		    }

		    if len(tt.target.Array) != len(tt.expt.Array) {
		        is_error = true
		    } else {
    		    for i := 0; i < len(tt.target.Array); i++ {
    		        if tt.target.Array[i] != tt.expt.Array[i] {
    		            is_error = true
    		            break
    		        }
    		    }
		    }

		    if len(tt.target.StructArray) != len(tt.expt.StructArray) {
		        is_error = true
		    } else {
    		    for i := 0; i < len(tt.target.StructArray); i++ {
    		        if tt.target.StructArray[i] != tt.expt.StructArray[i] {
    		            is_error = true
    		            break
    		        }
    		    }
		    }

		}

		if is_error {
			t.Errorf("MapStructv2(%v, %v) => (%v, %v) want (%v, %v)", t1{}, tt.in, tt.target, err, tt.expt, tt.out)
		}
	}
}

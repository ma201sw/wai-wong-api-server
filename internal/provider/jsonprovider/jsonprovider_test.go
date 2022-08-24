package jsonprovider

import "testing"

func Test_JSONMapToFloatSliceAs(t *testing.T) {
	t.Parallel()

	type args struct {
		data map[string]interface{}
	}

	tests := []struct {
		name        string
		c           jsonProviderImpl
		args        args
		expectedSum int
	}{
		{
			name: "JSONMapToFloatSliceAs-success",
			args: args{
				data: map[string]interface{}{
					"data1": []interface{}{float64(1), float64(2), float64(3), float64(4)},
					"data2": map[string]interface{}{"a": float64(6), "b": float64(4)},
					"data3": []interface{}{[]interface{}{[]interface{}{float64(2)}}},
					"data4": map[string]interface{}{"a": map[string]interface{}{"b": float64(4)}, "c": float64(-2)},
					"data5": map[string]interface{}{"a": []interface{}{float64(-1), float64(1), "dark"}},
					"data6": []interface{}{float64(-1), map[string]interface{}{"a": float64(1), "b": "light"}},
					"data7": []interface{}{},
					"data8": map[string]interface{}{},
					"data9": []interface{}{[]interface{}{map[string]interface{}{"a": float64(1)}}},
				},
			},
			c:           jsonProviderImpl{},
			expectedSum: 25,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			floatSlice := []float64{}
			tt.c.JSONMapToFloatSliceAs(tt.args.data, &floatSlice)

			sumResult := 0
			for _, v := range floatSlice {
				sumResult += int(v)
			}

			if tt.expectedSum != sumResult {
				t.Fatalf("expected sum: %v, got: %v", tt.expectedSum, sumResult)
			}
		})
	}
}

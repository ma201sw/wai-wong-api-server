package jsonprovider

type JSONProviderClientImplMock struct {
	JSONMapToFloatSliceAsFn func(data map[string]interface{}, out *[]float64)
}

func (c *JSONProviderClientImplMock) JSONMapToFloatSliceAs(data map[string]interface{}, out *[]float64) {
	if c != nil && c.JSONMapToFloatSliceAsFn != nil {
		c.JSONMapToFloatSliceAsFn(data, out)

		return
	}

	jsonProviderSrv := New()

	jsonProviderSrv.JSONMapToFloatSliceAs(data, out)
}

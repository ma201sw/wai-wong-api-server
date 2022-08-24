package jsonprovider

type Service interface {
	JSONMapToFloatSliceAs(data map[string]interface{}, out *[]float64)
}

type jsonProviderImpl struct{}

// verify interface compliance
var _ Service = (*jsonProviderImpl)(nil)

func New() jsonProviderImpl {
	return jsonProviderImpl{}
}

func (c jsonProviderImpl) JSONMapToFloatSliceAs(data map[string]interface{}, out *[]float64) {
	for _, rootElement := range data {
		switch rootElementTypeAsserted := rootElement.(type) {
		case map[string]interface{}: // if a map
			c.JSONMapToFloatSliceAs(rootElementTypeAsserted, out)
		case []interface{}: // if a slice
			for _, sliceElement := range rootElementTypeAsserted {
				switch sliceElementTypeAsserted := sliceElement.(type) {
				case map[string]interface{}:
					c.JSONMapToFloatSliceAs(sliceElementTypeAsserted, out)
				case float64:
					*out = append(*out, sliceElementTypeAsserted)
				case []interface{}:
					c.handleEmbeddedSlice(sliceElementTypeAsserted, out)
				}
			}
		case float64:
			*out = append(*out, rootElementTypeAsserted)
		}
	}
}

func (c jsonProviderImpl) handleEmbeddedSlice(data []interface{}, out *[]float64) {
	for _, rootElement := range data {
		switch rootElementTypeAsserted := rootElement.(type) {
		case map[string]interface{}:
			c.JSONMapToFloatSliceAs(rootElementTypeAsserted, out)
		case float64:
			*out = append(*out, rootElementTypeAsserted)
		case []interface{}:
			c.handleEmbeddedSlice(rootElementTypeAsserted, out)
		}
	}
}

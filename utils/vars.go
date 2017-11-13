package utils

var (
	globalVars map[string]interface{}
)

func init() {
	globalVars = make(map[string]interface{})
}

func SetVar(key string, variable interface{}, readOnly bool) {
	if _, ok := globalVars[key]; ok && readOnly {
		panic("key:" + key + " read only")
	}

	globalVars[key] = variable
}

func GetVar(key string) interface{}{
	variable, ok := globalVars[key]
	if !ok {
		return nil
	}

	return variable
}

package ontmobile

// ParameterType : the type of the parameter
type ParameterType int

const (
  Address     ParameterType = 0
  String      ParameterType = 1
  Integer     ParameterType = 2
  Fixed8      ParameterType = 3
  Array       ParameterType = 4
)

// Parameter : an invocation transaction parameter
type Parameter struct {
  Type  ParameterType
  Value interface{}
}

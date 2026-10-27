package golang

// The attributes model: a data type for every factory conversion kind
// (string, bool, int, int64, int32 and int8 through the instanceTypeName
// annotation, int16, byte, []byte, float64, float32, reflect.Type, time.Time,
// unsupported for java.lang.Object and an unknown class), data types with the
// defaultValue, testValue and instanceTypeName annotations, an EChar data
// type (no Java primitive default) next to the same type under another name,
// a non serializable data type; and a class with a string attribute with a
// default value literal, a boolean attribute, attributes named as Go keywords
// (type, func, range), an unsettable single attribute, an unsettable and a
// plain attribute list, a feature map list, an EJavaObject attribute.
// The quoted default literal and the attribute of an unknown Go type, which v1
// cannot generate, are in the v1bugs model.
func init() {
	register(fixture{Name: "attributes", Model: "attributes.ecore",
		V1Differences: map[string]string{
			"attributes/holder_impl_test.go": v1GoIntTestValue,
			"attributes/holder_mock_test.go": v1GoIntTestValue,
		}})
	// With isSetFlags, every single-valued changeable attribute keeps the ESet flag of the unsettable
	// one: the setter raises it, EIsSet reads it, EUnset sets the default and lowers it.
	register(fixture{Name: "attributes-issetflags", Model: "attributes.ecore",
		Props: map[string]string{"isSetFlags": "true"},
		V1Differences: map[string]string{
			"attributes/holder_impl.go":      v1IsSetFlags,
			"attributes/holder_impl_test.go": v1GoIntTestValue,
			"attributes/holder_mock_test.go": v1GoIntTestValue,
		}})
	// With floatFormat=xsd, the float data types are written by ecore.FormatFloatXSD.
	register(fixture{Name: "attributes-floatxsd", Model: "attributes.ecore",
		Props: map[string]string{"floatFormat": "xsd"},
		V1Differences: map[string]string{
			"attributes/attributesfactory_impl.go": v1FloatFormat,
			"attributes/holder_impl_test.go":       v1GoIntTestValue,
			"attributes/holder_mock_test.go":       v1GoIntTestValue,
		}})
}

// v1FloatFormat is why the factory differs from v1 with the property floatFormat=xsd.
const v1FloatFormat = "v1 writes floats with strconv.FormatFloat: Go's +Inf, -Inf, and 309 digits for 1.7976931348623157e+308"

// v1IsSetFlags is why the implementation differs from v1 with the property isSetFlags.
const v1IsSetFlags = "v1 has no isSetFlags: an attribute set to its default reads as unset"

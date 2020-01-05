package model

// Field field to add to each item
type Field struct {
	Operator  string
	Parameter string
	Selector  string
	Regexp    *RegexOperation
	Sprintf   *string
	Action    *string
}

// RegexOperation regexp to change field value
type RegexOperation struct {
	Expression string
	Group      int
}

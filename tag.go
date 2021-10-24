package seekpo

type Tag struct {
	Field Field
	Name  string
	Unit  string
	Type  Type
}

type Field = string

type Type int

package seekpo

type Tag struct {
	Code Code
	Name string
	Unit string
	Type Type
}

type Code string

type Type int

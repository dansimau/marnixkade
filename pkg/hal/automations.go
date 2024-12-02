package hal

type Automation struct {
	Name     string
	Entities Entities
	Action   func()
}

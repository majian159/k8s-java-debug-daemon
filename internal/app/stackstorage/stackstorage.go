package stackstorage

type ContainerStackModel struct {
	Namespace     string
	PodName       string
	ContainerName string
	Node          string
	Stack         string
}

type StackStorage interface {
	Store(model ContainerStackModel) error
}

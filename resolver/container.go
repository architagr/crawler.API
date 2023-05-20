package resolver

type IContainer interface {
	RegisterBinding(key string, obj interface{})
	GetBinding(key string) interface{}
}

var containerObj IContainer

type container struct {
	mapObj map[string]interface{}
}

func InitContainer() {
	containerObj = new(container)

}
func GetContainer() IContainer {
	if containerObj == nil {
		InitContainer()
	}
	return containerObj
}
func (cont *container) RegisterBinding(key string, obj interface{}) {
	cont.mapObj[key] = obj

}
func (cont *container) GetBinding(key string) interface{} {
	return cont.mapObj[key]
}

package protobuf

type modelCreator func() Model

var creatorList []modelCreator

func Register(method modelCreator) {

	creatorList = append(creatorList, method)
}

type modelList map[string]modelCreator

var _list modelList

func List() modelList {
	return _list
}

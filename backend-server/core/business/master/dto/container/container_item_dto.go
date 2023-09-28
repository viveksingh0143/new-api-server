package container

type ContainerItemDto struct {
	ID    int64  `json:"id"`
	Code  string `json:"code"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Count int64  `json:"count"`
}

package types

type ContainerData struct {
	Ip string
	Id string
	Port string
}

type VarnishConfiguration struct {
	Host	string
	EnableESI bool
	Backends []ContainerData
}

type ExecuteCommand struct {
	Url string
	Token string
}
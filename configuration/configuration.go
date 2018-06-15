package configuration

type Application struct {
	StackToWatch string
	Filepath string
	Rancher Rancher
}

type Rancher struct {
	URL string
	AccessKey string
	SecretKey string
}

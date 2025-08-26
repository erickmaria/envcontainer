package types

type BuildOptions struct {
	ImageName    string
	Dockerfile   string
	BuildContext string
	NetworkMode  string
}

package container

import "github.com/docker/docker/api/types/mount"

type ImagePullWay uint8

const (
	ImagePullFromBinary ImagePullWay = iota
	ImagePullFromUrl
	ImagePullFromRegistry
)

type Container struct {
	Source      Image             `json:"source"`
	Name        string            `json:"name"` // if name is nil, docker will generate a random strings as name
	ImageName   string            `json:"image_name"`
	ExposePort  map[string]int    `json:"expose_port,omitempty"` // map[containerPort]hostPort
	Network     string            `json:"network"`               // host, bridge, none(default)
	Env         map[string]string `json:"env,omitempty"`         // map[key]value
	Endpoint    string            `json:"endpoint"`              // optional
	Mount       []mount.Mount     `json:"mount,omitempty"`
	Compose     *[]byte           `json:"compose,omitempty"`
	ServicePort int               `json:"service_port"` // agent service port
}
type Image struct {
	PullWay ImagePullWay `json:"pull_way"` // 0: from binary, 1: from url, 2: from registry

	FromBinary   []byte `json:"from_binary"`   // a `tar` ball file
	FromUrl      string `json:"from_url"`      // such as https://example.com/mysql.tar
	FromRegistry string `json:"from_registry"` // such as `docker.io`(default)
	RegistryPath string `json:"registry_path"` // such as `library/mysql`
	Tag          string `json:"tag"`           // such as `latest`
	// FromFile   string `json:"from_file"` // at `Core` side path of the image

	RegistryAuth string `json:"registry_auth"` // base64 encoded string
}

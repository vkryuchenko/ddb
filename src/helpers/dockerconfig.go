package helpers

type DockerConfig struct {
	Target     string   `json:"target"`
	Apiversion string   `json:"apiversion"`
	Registry   string   `json:"registry"`
	Images     []string `json:"images"`
}

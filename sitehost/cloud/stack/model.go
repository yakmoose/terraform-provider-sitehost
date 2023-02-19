package stack

//	type DockerFileService struct {
//		ContainerName string   `yaml:"container_name"`
//		EnvFile       []string `yaml:"env_file"`
//		Environment   []string
//		Expose        []string
//		Image         string
//		Labels        []string
//		Restart       string
//		Volumes       []string
//	}
//
//	type Service struct {
//		Build       string            `yaml:"build,omitempty"`
//		Image       string            `yaml:"image,omitempty"`
//		Command     []string          `yaml:"command,omitempty"`
//		Ports       []string          `yaml:"ports,omitempty"`
//		Environment map[string]string `yaml:"environment,omitempty"`
//		Volumes     []string          `yaml:"volumes,omitempty"`
//	}
//
//	type DockerFile struct {
//		Services map[string]DockerFileService
//		// do we need the netowrks / version etc ... this stuff is for simple mapping and generating
//		//Version  string
//		//Networks a map of maps... and stuf...
//	}
type (
	DockerFileService struct {
		Build       string            `yaml:"build,omitempty"`
		Image       string            `yaml:"image,omitempty"`
		Command     []string          `yaml:"command,omitempty"`
		Ports       []string          `yaml:"ports,omitempty"`
		Environment map[string]string `yaml:"environment,omitempty"`
		EnvFile     []string          `yaml:"env_file,omitempty"`
		Restart     string            `yaml:"restart,omitempty"`
		Labels      map[string]string `yaml:"labels,omitempty"`
		Volumes     []string          `yaml:"volumes,omitempty"`
	}

	Compose struct {
		Version  string                       `yaml:"version"`
		Services map[string]DockerFileService `yaml:"services"`
		Networks map[string]struct {
			Driver string `yaml:"driver,omitempty"`
		} `yaml:"networks,omitempty"`
		Volumes map[string]struct {
			Driver string `yaml:"driver,omitempty"`
		} `yaml:"volumes,omitempty"`
	}
)

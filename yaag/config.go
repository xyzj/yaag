package yaag

import "github.com/xyzj/yaag/yaag/models"

type Config struct {
	On bool

	BaseUrls map[string]string

	DocTitle string
	DocPath  string
}

func (c *Config) ResetDoc() {
	spec = &models.Spec{}
}

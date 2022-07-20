package yaag

import "github.com/xyzj/yaag/yaag/models"

// Config 配置
type Config struct {
	On bool

	BaseUrls map[string]string

	DocTitle string
	DocPath  string
}

// ResetDoc 重置
func (c *Config) ResetDoc() {
	spec = &models.Spec{}
}

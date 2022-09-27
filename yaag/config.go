package yaag

import "github.com/xyzj/yaag/yaag/models"

// Config 配置
type Config struct {
	BaseUrls map[string]string
	DocTitle string
	DocPath  string
	DocDir   string
	On       bool
}

// ResetDoc 重置
func (c *Config) ResetDoc() {
	spec = &models.Spec{}
}

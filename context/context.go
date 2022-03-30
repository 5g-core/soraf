package context

import (
)

var ServiceCtx = ServiceContext{}

func init() {
	ServiceSelf().Name = "soraf"

}

type ServiceContext struct {
	Name string
	UriScheme   string //models.UriScheme
	BindingIPv4 string
	SBIPort     int
	HttpIPv6Address string
}

// Create new  context
func ServiceSelf() *ServiceContext {
	return &ServiceCtx
}

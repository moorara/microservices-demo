package component

import (
	"github.com/moorara/goto/config"
)

var Config = struct {
	ComponentTest bool
	ServiceURL    string
}{
	ServiceURL: "http://localhost:4040",
}

func init() {
	config.Pick(&Config)
}

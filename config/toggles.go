package config

import "fmt"

// ModuleToggles controls which modules are active at runtime.
var ModuleToggles = map[string]bool{
	"perps":     false,
	"options":   false,
	"lending":   false,
	"aiagent":   false,
	"predict":   false,
	"flashloan": false, // Disabled for regulatory compliance
}

// IsModuleEnabled returns whether a module is enabled.
func IsModuleEnabled(moduleName string) bool {
	enabled, exists := ModuleToggles[moduleName]
	if !exists {
		return true
	}
	return enabled
}

// ErrModuleDisabled returns an error for disabled modules.
func ErrModuleDisabled(moduleName string) error {
	return fmt.Errorf("module %s is currently disabled", moduleName)
}

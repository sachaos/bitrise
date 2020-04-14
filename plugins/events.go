package plugins

import (
	"encoding/json"
	"fmt"

	"github.com/bitrise-io/go-utils/colorstring"
	log2 "github.com/bitrise-io/go-utils/log"
)

// TriggerEventName ...
type TriggerEventName string

const (
	// DidFinishRun ...
	DidFinishRun TriggerEventName = "DidFinishRun"
)

// TriggerEvent ...
func TriggerEvent(name TriggerEventName, payload interface{}) error {
	// Create plugin input
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	pluginConfig := PluginConfig{
		PluginConfigTriggerEventKey: string(name),
	}

	// Load plugins
	plugins, err := LoadPlugins(string(name))
	if err != nil {
		return err
	}

	log2.Infof(colorstring.Red("Payload passed to plugin: %s"), string(payloadBytes))

	// Run plugins
	for _, plugin := range plugins {
		if err := RunPluginByEvent(plugin, pluginConfig, payloadBytes); err != nil {
			return err
		}
	}

	return nil
}

// LoadPlugins ...
func LoadPlugins(eventName string) ([]Plugin, error) {
	routing, err := readPluginRouting()
	if err != nil {
		return []Plugin{}, err
	}

	pluginNames := []string{}
	for name, route := range routing.RouteMap {
		if route.TriggerEvent == eventName {
			pluginNames = append(pluginNames, name)
		}
	}

	plugins := []Plugin{}
	for _, name := range pluginNames {
		plugin, found, err := LoadPlugin(name)
		if err != nil {
			return []Plugin{}, err
		}
		if !found {
			return []Plugin{}, fmt.Errorf("Plugin (%s) exist in routing, but not found", name)
		}
		plugins = append(plugins, plugin)
	}

	return plugins, nil
}

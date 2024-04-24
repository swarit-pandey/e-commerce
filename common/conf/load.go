package conf

import (
	"fmt"

	"github.com/spf13/viper"
	"k8s.io/klog/v2"
)

// LoadAndUnmarhsal will use viper to load the configuration
// and unmarshal onto the given structure
func LoadAndUnmarhsal(path string, config any) error {
	v := viper.New()

	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	err := v.ReadInConfig()
	if err != nil {
		return fmt.Errorf("common: config: failed: read: %v", err)
	}

	err = v.Unmarshal(config)
	if err != nil {
		return fmt.Errorf("common: config: failed: unmarhal: %v", err)
	}

	klog.InfoS("loaded and unmarshaled")
	return nil
}

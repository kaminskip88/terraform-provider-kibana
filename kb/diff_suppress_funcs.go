package kb

import (
	"fmt"

	"github.com/elastic/go-ucfg"
	"github.com/elastic/go-ucfg/diff"
	ucfgjson "github.com/elastic/go-ucfg/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	log "github.com/sirupsen/logrus"
)

// suppressEquivalentJSON permit to compare json string
func suppressEquivalentJSON(k, old, new string, d *schema.ResourceData) bool {
	if old == "" {
		old = `{}`
	}
	if new == "" {
		new = `{}`
	}
	confOld, err := ucfgjson.NewConfig([]byte(old), ucfg.PathSep("."))
	if err != nil {
		fmt.Printf("[ERR] Error when converting current Json: %s\ndata: %s", err.Error(), old)
		log.Errorf("Error when converting current Json: %s\ndata: %s", err.Error(), old)
		return false
	}
	confNew, err := ucfgjson.NewConfig([]byte(new), ucfg.PathSep("."))
	if err != nil {
		fmt.Printf("[ERR] Error when converting new Json: %s\ndata: %s", err.Error(), new)
		log.Errorf("Error when converting new Json: %s\ndata: %s", err.Error(), new)
		return false
	}

	currentDiff := diff.CompareConfigs(confOld, confNew)
	log.Debugf("Diff\n: %s", currentDiff.GoStringer())

	return !currentDiff.HasChanged()
}

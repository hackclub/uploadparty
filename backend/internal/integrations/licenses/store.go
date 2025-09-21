package licenses

import (
	"errors"
	"log"
	"strings"

	"github.com/uploadparty/app/backend/config"
)

// LicenseStore is an abstraction over the external license directory.
// Keep this generic to avoid leaking provider details in the public code.
// Implementations must be safe for concurrent use.
type LicenseStore interface {
	// Ping verifies the store connectivity with a very cheap request.
	Ping() error
}

// DefaultStore holds the initialized store if configured; may be a no-op.
var DefaultStore LicenseStore = &noopStore{}

// Init configures the DefaultStore based on environment configuration.
// To avoid disclosing provider details in logs for this open-source repo,
// we only log generic statuses.
func Init(cfg *config.Config) error {
	provider := strings.ToLower(strings.TrimSpace(cfg.LicensesProvider))
	if provider == "" || provider == "none" {
		DefaultStore = &noopStore{}
		log.Println("[licenses] external license store: disabled")
		return nil
	}

	switch provider {
	case "airtable":
		baseID, tableName := parseDSN(cfg.LicensesDSN)
		if cfg.LicensesToken == "" || baseID == "" || tableName == "" {
			return errors.New("license store is enabled but missing credentials or DSN parts")
		}
		store := newAirtableStore(cfg.LicensesToken, baseID, tableName)
		if err := store.Ping(); err != nil {
			return err
		}
		DefaultStore = store
		log.Println("[licenses] external license store: enabled")
		return nil
	default:
		return errors.New("unsupported license provider")
	}
}

// parseDSN expects semi-colon separated key=value pairs, but we only need base and table.
// Example: "base=appXXXXXXXXXXXX;table=Licenses"
func parseDSN(dsn string) (baseID, table string) {
	parts := strings.Split(dsn, ";")
	for _, p := range parts {
		kv := strings.SplitN(strings.TrimSpace(p), "=", 2)
		if len(kv) != 2 {
			continue
		}
		k := strings.ToLower(strings.TrimSpace(kv[0]))
		v := strings.TrimSpace(kv[1])
		switch k {
		case "base":
			baseID = v
		case "table":
			table = v
		}
	}
	return
}

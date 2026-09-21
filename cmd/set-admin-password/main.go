package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"cyberstrike-ai/internal/config"
	"cyberstrike-ai/internal/database"
	"cyberstrike-ai/internal/security"

	"go.uber.org/zap"
)

func main() {
	var configPath = flag.String("config", "config.yaml", "Path to the configuration file")
	var password = flag.String("password", "", "Plain-text admin password to set (alternative to env CYBERSTRIKE_ADMIN_PASSWORD)")
	flag.Parse()

	cp := strings.TrimSpace(*configPath)
	if cp == "" {
		cp = "config.yaml"
	}

	pw := strings.TrimSpace(*password)
	if pw == "" {
		pw = strings.TrimSpace(os.Getenv("CYBERSTRIKE_ADMIN_PASSWORD"))
	}
	if pw == "" {
		fmt.Fprintln(os.Stderr, "empty admin password; pass -password <pw> or set CYBERSTRIKE_ADMIN_PASSWORD")
		os.Exit(2)
	}
	if len(pw) < 8 {
		fmt.Fprintln(os.Stderr, "admin password must be at least 8 characters")
		os.Exit(2)
	}

	cfg, err := config.Load(cp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	dbPath := strings.TrimSpace(cfg.Database.Path)
	if dbPath == "" {
		dbPath = "data/conversations.db"
	}
	if dir := filepathDir(dbPath); dir != "" {
		if mkErr := os.MkdirAll(dir, 0o755); mkErr != nil {
			fmt.Fprintf(os.Stderr, "failed to create data dir: %v\n", mkErr)
			os.Exit(1)
		}
	}

	hash, err := security.HashPassword(pw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to hash password: %v\n", err)
		os.Exit(1)
	}

	db, err := database.NewDB(dbPath, zap.NewNop())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open database: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = db.Close() }()

	needs, err := db.RBACNeedsAdminPassword()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to inspect database: %v\n", err)
		os.Exit(1)
	}

	if needs {
		if err := db.BootstrapRBAC(hash, security.PermissionCatalog); err != nil {
			fmt.Fprintf(os.Stderr, "failed to bootstrap admin account: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("initialized built-in admin account with the provided password\n")
	} else {
		admin, err := db.GetRBACUserByUsername("admin")
		if err != nil {
			fmt.Fprintf(os.Stderr, "built-in admin account not found: %v\n", err)
			os.Exit(1)
		}
		if !admin.IsBuiltin {
			fmt.Fprintln(os.Stderr, "admin account is not built-in; refusing to overwrite it")
			os.Exit(1)
		}
		if err := db.UpdateRBACAdminPassword(hash); err != nil {
			fmt.Fprintf(os.Stderr, "failed to update admin password: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("reset built-in admin password to the provided value\n")
	}
}

func filepathDir(p string) string {
	i := strings.LastIndexAny(p, `/\`)
	if i < 0 {
		return ""
	}
	return p[:i]
}

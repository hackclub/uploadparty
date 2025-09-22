package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/uploadparty/app/config"
	uploadDB "github.com/uploadparty/app/pkg/db"
)

func main() {
	var dir string
	var list bool
	var dry bool
	var env string
	flag.StringVar(&dir, "dir", "migrations", "directory containing .sql migration files (relative to backend/)")
	flag.BoolVar(&list, "list", false, "list migration files without executing")
	flag.BoolVar(&dry, "dry-run", false, "print migration SQL without executing")
	flag.StringVar(&env, "env", os.Getenv("MIGRATIONS_ENV"), "environment selector: '' (default) or 'dev' to include *.dev.sql files")
	flag.Parse()

	cfg := config.Load()
	db, err := uploadDB.Connect(cfg)
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}

	// Resolve dir to absolute based on current working directory
	absDir, err := filepath.Abs(dir)
	if err != nil {
		log.Fatalf("resolve dir: %v", err)
	}

	includeDev := strings.EqualFold(strings.TrimSpace(env), "dev")

	entries := []string{}
	err = filepath.WalkDir(absDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(d.Name()) == ".sql" {
			name := d.Name()
			if strings.HasSuffix(name, ".dev.sql") && !includeDev {
				return nil
			}
			entries = append(entries, path)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("scan migrations: %v", err)
	}
	if len(entries) == 0 {
		log.Printf("no .sql files found in %s", absDir)
		return
	}
	// Ensure deterministic order
	sort.Strings(entries)

	fmt.Printf("Found %d migration file(s) in %s (env=%s)\n", len(entries), absDir, env)
	for _, f := range entries {
		fmt.Println(" -", filepath.Base(f))
	}
	if list {
		return
	}

	for _, file := range entries {
		fmt.Printf("\n==> Applying %s\n", filepath.Base(file))
		b, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("read %s: %v", file, err)
		}
		if dry {
			fmt.Printf("-- DRY RUN: not executing --\n%s\n", string(b))
			continue
		}
		if err := db.Exec(string(b)).Error; err != nil {
			log.Fatalf("apply %s: %v", file, err)
		}
		fmt.Println("OK")
	}

	fmt.Println("\nAll migrations applied successfully.")
}

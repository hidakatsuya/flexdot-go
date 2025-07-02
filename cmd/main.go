package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hidakatsuya/flexdot-go/internal/clearbackups"
	"github.com/hidakatsuya/flexdot-go/internal/config"
	initcmd "github.com/hidakatsuya/flexdot-go/internal/init"
	"github.com/hidakatsuya/flexdot-go/internal/install"
)

const version = "0.3.1"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	arg := os.Args[1]
	if arg == "--version" || arg == "-v" {
		fmt.Println(version)
		return
	}

	switch arg {
	case "install":
		runInstall(os.Args[2:])
	case "init":
		if err := initcmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to init: %v\n", err)
			os.Exit(1)
		}
	case "clear-backups":
		if err := clearbackups.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to clear backups: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", arg)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	usage := `
Usage: flexdot <command> [options]
Commands:
  install [-H|--home_dir path] <index.yml>
  init
  clear-backups`
	fmt.Println(usage)
}

func runInstall(args []string) {
	fs := flag.NewFlagSet("install", flag.ExitOnError)
	homeDirFlag := fs.String("home_dir", "", "Home directory")
	homeDirShortFlag := fs.String("H", "", "Home directory (shorthand)")
	fs.Usage = func() {
		printUsage()
	}
	fs.Parse(args)

	rest := fs.Args()
	var indexFile string
	if len(rest) > 1 {
		fmt.Fprintf(os.Stderr, "Too many arguments for install command\n")
		printUsage()
		os.Exit(1)
	}
	if len(rest) == 1 {
		indexFile = rest[0]
	}

	// dotfilesDir = current directory
	dotfilesDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get current directory: %v\n", err)
		os.Exit(1)
	}

	// Load config.yml if present
	cfg, err := config.LoadConfig(dotfilesDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config.yml: %v\n", err)
		os.Exit(1)
	}

	// Determine homeDir and indexFile (priority: CLI > config.yml > error)
	homeDir := ""
	if *homeDirFlag != "" {
		homeDir = *homeDirFlag
	} else if *homeDirShortFlag != "" {
		homeDir = *homeDirShortFlag
	} else if cfg != nil && cfg.HomeDir != "" {
		homeDir = cfg.HomeDir
	}
	if homeDir == "" {
		fmt.Fprintf(os.Stderr, "home_dir must be specified by --home_dir/-h or config.yml\n")
		os.Exit(1)
	}

	if indexFile == "" && cfg != nil && cfg.IndexYml != "" {
		indexFile = cfg.IndexYml
	}
	if indexFile == "" {
		fmt.Fprintf(os.Stderr, "<index.yml> must be specified as argument or config.yml\n")
		os.Exit(1)
	}

	keepMaxBackupCount := cfg.GetKeepMaxCount()

	// indexFile may be relative to dotfilesDir
	if !filepath.IsAbs(indexFile) {
		indexFile = filepath.Join(dotfilesDir, indexFile)
	}

	if err := install.Run(indexFile, homeDir, dotfilesDir, keepMaxBackupCount); err != nil {
		fmt.Fprintf(os.Stderr, "Install failed: %v\n", err)
		os.Exit(1)
	}
}

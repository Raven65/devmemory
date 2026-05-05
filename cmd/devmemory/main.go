package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"text/tabwriter"
	"time"

	"devmemory/internal/config"
	"devmemory/internal/core"
	"devmemory/internal/server"
	"devmemory/internal/service"
	"devmemory/internal/store"

	fyneApp "devmemory/internal/ui/fyne"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

// waitForExit pauses on Windows so the console doesn't close immediately
// when double-clicking the exe. No-op on other platforms.
func waitForExit() {
	if runtime.GOOS == "windows" {
		fmt.Fprintln(os.Stderr, "\nPress Enter to exit...")
		bufio.NewReader(os.Stdin).ReadString('\n')
	}
}

func main() {
	if len(os.Args) < 2 {
		// No args: launch Fyne GUI
		cmdGUI()
		return
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "version", "--version", "-v":
		fmt.Printf("devmemory %s (commit: %s, built: %s)\n", Version, Commit, Date)
	case "gui":
		cmdGUI()
	case "add":
		cmdAdd(args)
	case "list", "ls":
		cmdList(args)
	case "show":
		cmdShow(args)
	case "delete", "rm":
		cmdDelete(args)
	case "edit":
		cmdEdit(args)
	case "search", "s":
		cmdSearch(args)
	case "today":
		cmdToday(args)
	case "export":
		cmdExport(args)
	case "import":
		cmdImport(args)
	case "serve":
		cmdServe(args)
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", cmd)
		printUsage()
		waitForExit()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "devmemory — personal memory & action hub for developers")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "  devmemory <command> [options] [args]")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Commands:")
	fmt.Fprintln(os.Stderr, "  version              Show version")
	fmt.Fprintln(os.Stderr, "  gui                  Launch desktop GUI (default)")
	fmt.Fprintln(os.Stderr, "  add <content>        Add a new entry")
	fmt.Fprintln(os.Stderr, "  list                 List entries")
	fmt.Fprintln(os.Stderr, "  show <id>            Show entry details")
	fmt.Fprintln(os.Stderr, "  delete <id>          Delete an entry")
	fmt.Fprintln(os.Stderr, "  edit <id>            Edit an entry")
	fmt.Fprintln(os.Stderr, "  search <query>       Search entries")
	fmt.Fprintln(os.Stderr, "  today                Show today's entries")
	fmt.Fprintln(os.Stderr, "  export today|json    Export entries")
	fmt.Fprintln(os.Stderr, "  import <file>        Import from JSON")
	fmt.Fprintln(os.Stderr, "  serve                Start web UI server")
	fmt.Fprintln(os.Stderr, "  help                 Show this help")
}

// openService initializes config, ensures data directory, opens store, and creates a MemoryService.
func openService() (*service.MemoryService, func()) {
	cfg := config.DefaultConfig()
	if err := config.EnsureDataDir(cfg.DataDir); err != nil {
		fmt.Fprintf(os.Stderr, "error creating data directory: %v\n", err)
		os.Exit(1)
	}
	s := store.NewBBoltStore(cfg.DBPath)
	if err := s.Open(); err != nil {
		fmt.Fprintf(os.Stderr, "error opening database: %v\n", err)
		os.Exit(1)
	}
	return service.NewMemoryService(s), func() { s.Close() }
}

// cmdAdd handles: devmemory add <content> [--type] [--title] [--project] [--tags]
func cmdAdd(args []string) {
	content, entryType, title, project, tags := parseAddArgs(args)

	if content == "" {
		fmt.Fprintln(os.Stderr, "error: content is required")
		fmt.Fprintln(os.Stderr, "Usage: devmemory add <content> [--type T] [--title T] [--project P] [--tags a,b]")
		os.Exit(1)
	}

	svc, cleanup := openService()
	defer cleanup()

	entry, err := svc.CreateEntry(service.CreateEntryInput{
		Content: content,
		Type:    core.EntryType(entryType),
		Title:   title,
		Project: project,
		Tags:    splitTags(tags),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created %s entry: %s\n", entry.Type, entry.ID[:8])
	fmt.Printf("  %s\n", entry.TitleOrContent())
}

// cmdList handles: devmemory list [--type] [--project] [--tag]
func cmdList(args []string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	entryType := fs.String("type", "", "filter by type")
	project := fs.String("project", "", "filter by project")
	tag := fs.String("tag", "", "filter by tag")
	fs.Parse(args)

	svc, cleanup := openService()
	defer cleanup()

	entries, err := svc.ListEntries(store.ListOptions{
		Type:    core.EntryType(*entryType),
		Project: *project,
		Tag:     *tag,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error listing entries: %v\n", err)
		os.Exit(1)
	}

	if len(entries) == 0 {
		fmt.Println("No entries found.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "ID\tTYPE\tPROJECT\tTAGS\tTITLE\n")
	for _, e := range entries {
		id := e.ID[:8]
		tagsStr := strings.Join(e.Tags, ",")
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", id, e.Type, e.Project, tagsStr, e.TitleOrContent())
	}
	w.Flush()

	fmt.Printf("\n%d entries\n", len(entries))
}

// cmdSearch handles: devmemory search <query> [--type]
func cmdSearch(args []string) {
	query, entryType := parseSearchArgs(args)

	if query == "" {
		fmt.Fprintln(os.Stderr, "error: search query is required")
		os.Exit(1)
	}

	svc, cleanup := openService()
	defer cleanup()

	results, err := svc.SearchEntries(query, core.EntryType(entryType))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if len(results) == 0 {
		fmt.Printf("No results for '%s'\n", query)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "ID\tTYPE\tSCORE\tTITLE\n")
	for _, r := range results {
		id := r.Entry.ID[:8]
		fmt.Fprintf(w, "%s\t%s\t%.1f\t%s\n", id, r.Entry.Type, r.Score, r.Entry.TitleOrContent())
	}
	w.Flush()

	fmt.Printf("\n%d results\n", len(results))
}

// cmdToday handles: devmemory today
func cmdToday(args []string) {
	svc, cleanup := openService()
	defer cleanup()

	entries, err := svc.GetToday()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	now := time.Now()
	fmt.Printf("Today's entries (%s)\n", now.Format("2006-01-02"))
	fmt.Println(strings.Repeat("-", 50))

	if len(entries) == 0 {
		fmt.Println("  (no entries yet)")
		return
	}

	for _, e := range entries {
		timeStr := e.CreatedAt.Format("15:04")
		fmt.Printf("  %s  [%s]  %s", timeStr, e.Type, e.TitleOrContent())
		if e.Project != "" {
			fmt.Printf("  (%s)", e.Project)
		}
		if len(e.Tags) > 0 {
			fmt.Printf("  #%s", strings.Join(e.Tags, " #"))
		}
		fmt.Println()
	}

	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("%d entries\n", len(entries))
}

// cmdShow handles: devmemory show <id>
func cmdShow(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "error: entry ID is required")
		fmt.Fprintln(os.Stderr, "Usage: devmemory show <id>")
		os.Exit(1)
	}

	svc, cleanup := openService()
	defer cleanup()

	entry, err := svc.GetEntry(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("ID:        %s\n", entry.ID)
	fmt.Printf("Type:      %s\n", entry.Type)
	if entry.Title != "" {
		fmt.Printf("Title:     %s\n", entry.Title)
	}
	fmt.Printf("Content:   %s\n", entry.Content)
	if entry.Summary != "" {
		fmt.Printf("Summary:   %s\n", entry.Summary)
	}
	if entry.Project != "" {
		fmt.Printf("Project:   %s\n", entry.Project)
	}
	if len(entry.Tags) > 0 {
		fmt.Printf("Tags:      %s\n", strings.Join(entry.Tags, ", "))
	}
	fmt.Printf("Favorite:  %v\n", entry.Favorite)
	fmt.Printf("Dangerous: %v\n", entry.Dangerous)
	fmt.Printf("Archived:  %v\n", entry.Archived)
	fmt.Printf("UseCount:  %d\n", entry.UseCount)
	fmt.Printf("Created:   %s\n", entry.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated:   %s\n", entry.UpdatedAt.Format("2006-01-02 15:04:05"))
	if entry.LastUsedAt != nil {
		fmt.Printf("LastUsed:  %s\n", entry.LastUsedAt.Format("2006-01-02 15:04:05"))
	}
}

// cmdDelete handles: devmemory delete <id> [--force]
func cmdDelete(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "error: entry ID is required")
		fmt.Fprintln(os.Stderr, "Usage: devmemory delete <id> [--force]")
		os.Exit(1)
	}

	force := false
	for _, a := range args[1:] {
		if a == "--force" || a == "-f" {
			force = true
		}
	}

	svc, cleanup := openService()
	defer cleanup()

	if !force {
		entry, err := svc.GetEntry(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Delete this entry?\n  [%s] %s\n\n  (y/N): ", entry.Type, entry.TitleOrContent())
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
			fmt.Println("Cancelled.")
			return
		}
	}

	if err := svc.DeleteEntry(args[0]); err != nil {
		fmt.Fprintf(os.Stderr, "error deleting entry: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Deleted entry %s\n", args[0][:8])
}

// cmdEdit handles: devmemory edit <id> [--title] [--content] [--type] [--project] [--tags] [--favorite] [--archive]
func cmdEdit(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "error: entry ID is required")
		fmt.Fprintln(os.Stderr, "Usage: devmemory edit <id> [--title T] [--content C] [--type T] [--project P] [--tags a,b] [--favorite] [--archive]")
		os.Exit(1)
	}

	var title, content, entryType, project, tags string
	var favorite, archive *bool
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--title":
			i++
			if i < len(args) {
				title = args[i]
			}
		case "--content":
			i++
			if i < len(args) {
				content = args[i]
			}
		case "--type", "-t":
			i++
			if i < len(args) {
				entryType = args[i]
			}
		case "--project", "-p":
			i++
			if i < len(args) {
				project = args[i]
			}
		case "--tags":
			i++
			if i < len(args) {
				tags = args[i]
			}
		case "--favorite":
			b := true
			favorite = &b
		case "--archive":
			b := true
			archive = &b
		}
	}

	svc, cleanup := openService()
	defer cleanup()

	fields := service.UpdateEntryFields{}
	changed := false
	if title != "" {
		fields.Title = &title
		changed = true
	}
	if content != "" {
		fields.Content = &content
		changed = true
	}
	if entryType != "" {
		t := core.EntryType(entryType)
		fields.Type = &t
		changed = true
	}
	if project != "" {
		fields.Project = &project
		changed = true
	}
	if tags != "" {
		t := splitTags(tags)
		fields.Tags = &t
		changed = true
	}
	if favorite != nil {
		fields.Favorite = favorite
		changed = true
	}
	if archive != nil {
		fields.Archived = archive
		changed = true
	}

	if !changed {
		fmt.Fprintln(os.Stderr, "error: no fields specified to update")
		os.Exit(1)
	}

	entry, err := svc.UpdateEntry(args[0], fields)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error updating entry: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Updated entry %s\n", entry.ID[:8])
}

// cmdExport handles: devmemory export today|json [-o file]
func cmdExport(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "error: export subcommand required (today, json)")
		fmt.Fprintln(os.Stderr, "Usage: devmemory export today|json [-o file]")
		os.Exit(1)
	}

	subcmd := args[0]
	outputFile := ""
	for i := 1; i < len(args); i++ {
		if args[i] == "-o" && i+1 < len(args) {
			outputFile = args[i+1]
			i++
		}
	}

	svc, cleanup := openService()
	defer cleanup()

	switch subcmd {
	case "today":
		md, err := svc.ExportTodayMarkdown()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		if outputFile != "" {
			if err := os.WriteFile(outputFile, []byte(md), 0644); err != nil {
				fmt.Fprintf(os.Stderr, "error writing file: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Exported today's entries to %s\n", outputFile)
		} else {
			fmt.Print(md)
		}

	case "json":
		data, err := svc.ExportJSON()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		if outputFile != "" {
			if err := os.WriteFile(outputFile, data, 0644); err != nil {
				fmt.Fprintf(os.Stderr, "error writing file: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Exported entries to %s\n", outputFile)
		} else {
			fmt.Print(string(data))
		}

	default:
		fmt.Fprintf(os.Stderr, "unknown export subcommand: %s\n", subcmd)
		fmt.Fprintln(os.Stderr, "Usage: devmemory export today|json [-o file]")
		os.Exit(1)
	}
}

// cmdImport handles: devmemory import <file>
func cmdImport(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "error: file path is required")
		fmt.Fprintln(os.Stderr, "Usage: devmemory import <file>")
		os.Exit(1)
	}

	data, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
		os.Exit(1)
	}

	svc, cleanup := openService()
	defer cleanup()

	imported, err := svc.ImportJSON(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Imported %d entries\n", imported)
}

// cmdServe handles: devmemory serve [--port] [--no-open]
func cmdServe(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	port := fs.Int("port", 8420, "server port")
	noOpen := fs.Bool("no-open", false, "don't open browser")
	fs.Parse(args)

	cfg := config.DefaultConfig()
	if err := config.EnsureDataDir(cfg.DataDir); err != nil {
		fmt.Fprintf(os.Stderr, "error creating data directory: %v\n", err)
		os.Exit(1)
	}

	s := store.NewBBoltStore(cfg.DBPath)
	if err := s.Open(); err != nil {
		fmt.Fprintf(os.Stderr, "error opening database: %v\n", err)
		os.Exit(1)
	}
	defer s.Close()

	svc := service.NewMemoryService(s)
	srv := server.New(svc, *port)
	fmt.Printf("DevMemory web UI: http://127.0.0.1:%d\n", *port)
	if err := srv.Start(!*noOpen); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// cmdGUI launches the Fyne desktop GUI.
func cmdGUI() {
	a := fyneApp.NewApp()
	a.Run()
}

// parseAddArgs extracts content and flags from mixed-position arguments.
func parseAddArgs(args []string) (content, entryType, title, project, tags string) {
	var contentParts []string
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--type", "-t":
			i++
			if i < len(args) {
				entryType = args[i]
			}
		case "--title":
			i++
			if i < len(args) {
				title = args[i]
			}
		case "--project", "-p":
			i++
			if i < len(args) {
				project = args[i]
			}
		case "--tags":
			i++
			if i < len(args) {
				tags = args[i]
			}
		default:
			contentParts = append(contentParts, args[i])
		}
	}
	content = strings.Join(contentParts, " ")
	return
}

// parseSearchArgs extracts query and optional type from mixed-position arguments.
func parseSearchArgs(args []string) (query, entryType string) {
	var queryParts []string
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--type", "-t":
			i++
			if i < len(args) {
				entryType = args[i]
			}
		default:
			queryParts = append(queryParts, args[i])
		}
	}
	query = strings.Join(queryParts, " ")
	return
}

func splitTags(s string) []string {
	parts := strings.Split(s, ",")
	tags := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			tags = append(tags, p)
		}
	}
	return tags
}

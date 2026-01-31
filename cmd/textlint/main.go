package textlint

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"textlint/internal/fix"
	"textlint/internal/log"
	"textlint/internal/report"
	"textlint/internal/rules"
	"textlint/internal/scanner"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	FormatCompact = "compact"
	FormatTAP     = "tap"
	FormatJSON    = "json"
)

var (
	checkFlag      bool
	fixFlag        bool
	formatStr      string
	quietFlag      bool
	verboseFlag    bool
	debugFlag      bool
	checkHadIssues bool
)

var rootCmd = &cobra.Command{
	Use:   "textlint [--check|--fix] [--format=compact|tap|json] [--quiet|--verbose|--debug] <path> [path ...]",
	Short: "Check or fix text files",
	Long:  "Check or fix text files. Pass one or more files or directories (recursive). Only text files (valid UTF-8, no null bytes) are processed.",
	Args:  cobra.ArbitraryArgs,
	RunE:  runE,
}

func init() {
	rootCmd.Flags().BoolVar(&checkFlag, "check", false, "check files and report issues (default)")
	rootCmd.Flags().BoolVar(&fixFlag, "fix", false, "fix issues in place")
	rootCmd.Flags().StringVar(&formatStr, "format", FormatCompact, "output format: compact, tap, or json")
	rootCmd.PersistentFlags().BoolVarP(&quietFlag, "quiet", "q", false, "quiet: only fatal errors")
	rootCmd.PersistentFlags().BoolVar(&verboseFlag, "verbose", false, "verbose: steps, skipped files, timing")
	rootCmd.PersistentFlags().BoolVar(&debugFlag, "debug", false, "debug: internal state, non-text files skipped with reason")
	rootCmd.SetHelpFunc(helpFunc)
}

var verbosityFlagNames = map[string]bool{"quiet": true, "q": true, "verbose": true, "debug": true}

func helpFunc(cmd *cobra.Command, args []string) {
	out := cmd.OutOrStderr()
	fmt.Fprintf(out, "%s\n\n", cmd.Short)
	if cmd.Long != "" {
		fmt.Fprintf(out, "%s\n\n", cmd.Long)
	}
	fmt.Fprintf(out, "Usage:\n  %s\n\n", cmd.UseLine())
	fmt.Fprintln(out, "Options:")
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if !verbosityFlagNames[f.Name] {
			printFlagUsage(out, f)
		}
	})
	cmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		if !verbosityFlagNames[f.Name] {
			printFlagUsage(out, f)
		}
	})
	fmt.Fprintln(out, "Verbosity:")
	cmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		if verbosityFlagNames[f.Name] {
			printFlagUsage(out, f)
		}
	})
}

func printFlagUsage(out io.Writer, f *pflag.Flag) {
	if f.Shorthand != "" && f.Name != f.Shorthand {
		fmt.Fprintf(out, "  -%s, --%s\t%s\n", f.Shorthand, f.Name, f.Usage)
	} else {
		fmt.Fprintf(out, "      --%s\t%s\n", f.Name, f.Usage)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if checkHadIssues {
		os.Exit(1)
	}
}

func verbosityLevel() log.Level {
	if debugFlag {
		return log.Debug
	}
	if verboseFlag {
		return log.Verbose
	}
	if quietFlag {
		return log.Quiet
	}
	return log.Normal
}

func runE(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		helpFunc(cmd, nil)
		return nil
	}
	if checkFlag && fixFlag {
		return fmt.Errorf("cannot use both --check and --fix")
	}
	if !checkFlag && !fixFlag {
		checkFlag = true
	}
	validFormats := map[string]bool{FormatCompact: true, FormatTAP: true, FormatJSON: true}
	if !validFormats[formatStr] {
		return fmt.Errorf("invalid format %q (use compact, tap, or json)", formatStr)
	}
	log.SetLevel(verbosityLevel())
	hadIssues, err := run(checkFlag, fixFlag, formatStr, args)
	if err != nil {
		return err
	}
	checkHadIssues = checkFlag && hadIssues
	return nil
}

func run(check, doFix bool, format string, paths []string) (hadIssues bool, err error) {
	start := time.Now()
	lvl := log.GetLevel()
	if lvl >= log.Debug {
		log.Logf(log.Debug, "debug: check=%v format=%s paths=%v\n", check, format, paths)
	}
	if lvl >= log.Verbose {
		log.Logf(log.Verbose, "Scanning %d path(s): %s\n", len(paths), strings.Join(paths, ", "))
	}
	files, skipped, err := scanner.Scan(paths)
	if err != nil {
		return false, err
	}
	elapsedScan := time.Since(start)
	if lvl >= log.Verbose {
		if len(files) == 0 {
			log.Logf(log.Verbose, "No text files found. Scanned 0 text file(s), skipped %d path(s).\n", len(skipped))
		} else {
			log.Logf(log.Verbose, "Scanned %d text file(s), skipped %d path(s).\n", len(files), len(skipped))
		}
	}
	if lvl >= log.Debug {
		for _, p := range sortedKeys(skipped) {
			log.Logf(log.Debug, "scanner: rejected %s (reason: %s)\n", p, skipped[p])
		}
		for _, p := range files {
			log.Logf(log.Debug, "scanner: accepted %s\n", p)
		}
	}
	if lvl >= log.Verbose && lvl < log.Debug {
		for p := range skipped {
			log.Logf(log.Verbose, "Skipped (not text): %s\n", p)
		}
	}
	if len(files) == 0 {
		if lvl >= log.Normal {
			fmt.Fprintln(os.Stdout, "No text files found.")
			if check {
				report.Write(os.Stdout, report.Format(format), nil, 0)
			}
		}
		return false, nil
	}
	var allIssues []rules.Issue
	fileIssues := make(map[string][]rules.Issue)
	for _, path := range files {
		if lvl >= log.Verbose {
			if check {
				log.Logf(log.Verbose, "Checking %s\n", path)
			} else {
				log.Logf(log.Verbose, "Fixing %s\n", path)
			}
		}
		issues, err := rules.CheckFile(path)
		if err != nil {
			return false, err
		}
		if len(issues) > 0 {
			fileIssues[path] = issues
			allIssues = append(allIssues, issues...)
			if lvl >= log.Debug {
				ruleIDs := make(map[string]bool)
				for _, i := range issues {
					ruleIDs[i.RuleID] = true
				}
				var ids []string
				for id := range ruleIDs {
					ids = append(ids, id)
				}
				sort.Strings(ids)
				log.Logf(log.Debug, "rules: %s -> %d issue(s): %s\n", path, len(issues), strings.Join(ids, ", "))
			}
		}
	}
	if check {
		if lvl >= log.Normal {
			if err := report.Write(os.Stdout, report.Format(format), allIssues, len(files)); err != nil {
				return false, err
			}
		}
		elapsed := time.Since(start)
		if lvl >= log.Verbose {
			log.Logf(log.Verbose, "Completed in %s\n", elapsed.Round(time.Millisecond))
		}
		_ = elapsedScan
		return len(allIssues) > 0, nil
	}
	for path := range fileIssues {
		if err := fix.Apply(path); err != nil {
			return false, err
		}
		if lvl >= log.Debug {
			log.Logf(log.Debug, "fix: applied to %s\n", path)
		}
	}
	if lvl >= log.Normal && len(fileIssues) > 0 {
		fmt.Fprintf(os.Stdout, "Fixed %d file(s).\n", len(fileIssues))
	}
	elapsed := time.Since(start)
	if lvl >= log.Verbose {
		log.Logf(log.Verbose, "Completed in %s\n", elapsed.Round(time.Millisecond))
	}
	return false, nil
}

func sortedKeys(m map[string]string) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

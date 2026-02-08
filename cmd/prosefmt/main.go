package prosefmt

import (
	"fmt"
	"io"
	"os"
	"prosefmt/internal/fix"
	"prosefmt/internal/log"
	"prosefmt/internal/report"
	"prosefmt/internal/rules"
	"prosefmt/internal/scanner"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	version        = "dev"
	checkHadIssues bool
)

const rootDescription = "The simplest text formatter for making your files look correct."

var rootCmd = &cobra.Command{
	Use:   "prosefmt [command]",
	Short: rootDescription,
	Long:  rootDescription,
	Args:  cobra.ArbitraryArgs,
	RunE:  rootRunE,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}

var checkCmd = &cobra.Command{
	Use:   "check [flags] paths...",
	Short: "Review the given paths for format issues (default)",
	Long:  "Recursively scan the given paths and check any non-binary files for format issues. Exit with code 1 if any problems are detected; otherwise, exit with code 0. This is the default behavior when no command is specified.",
	Args:  cobra.ArbitraryArgs,
	RunE:  checkRunE,
}

var writeCmd = &cobra.Command{
	Use:   "write [flags] paths...",
	Short: "Apply fixes in place to the given paths",
	Long:  "Recursively scan the given paths and fix format issues in any non-binary files.",
	Args:  cobra.ArbitraryArgs,
	RunE:  writeRunE,
}

func addOutputFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("silent", false, "No output printed")
	cmd.Flags().Bool("compact", false, "Show formatted or errored files (default)")
	cmd.Flags().Bool("verbose", false, "Print debug output (steps, scanner, rules, timing)")
}

func outputLevelFromCmd(cmd *cobra.Command) log.Level {
	silent, _ := cmd.Flags().GetBool("silent")
	compact, _ := cmd.Flags().GetBool("compact")
	verbose, _ := cmd.Flags().GetBool("verbose")
	if verbose {
		return log.Verbose
	}
	if compact {
		return log.Normal
	}
	if silent {
		return log.Silent
	}
	return log.Normal
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(writeCmd)
	addOutputFlags(checkCmd)
	addOutputFlags(writeCmd)
	rootCmd.SetHelpFunc(rootHelpFunc)
	checkCmd.SetHelpFunc(commandHelpFunc)
	writeCmd.SetHelpFunc(commandHelpFunc)
}

var outputFlagOrder = []string{"silent", "compact", "verbose"}

func rootHelpFunc(cmd *cobra.Command, args []string) {
	out := cmd.OutOrStderr()
	if cmd.Short != "" {
		fmt.Fprintf(out, "%s\n\n", cmd.Short)
	}
	fmt.Fprintf(out, "Usage:\n  %s\n\n", cmd.UseLine())
	if len(cmd.Commands()) > 0 {
		fmt.Fprintln(out, "Commands:")
		maxLen := 0
		for _, c := range cmd.Commands() {
			if c.IsAvailableCommand() && len(c.Name()) > maxLen {
				maxLen = len(c.Name())
			}
		}
		for _, c := range cmd.Commands() {
			if c.IsAvailableCommand() {
				fmt.Fprintf(out, "  %-*s  %s\n", maxLen, c.Name(), c.Short)
			}
		}
		fmt.Fprintln(out, "")
	}
	fmt.Fprintln(out, "With no command, runs 'check' by default. Use 'check' or 'write' for output options (--silent, --compact, --verbose).")
	if version != "" {
		fmt.Fprintf(out, "\nVersion: %s\n", version)
	}
}

func wrapWords(s string, width int) string {
	if width <= 0 {
		return s
	}
	var b strings.Builder
	for _, para := range strings.Split(s, "\n") {
		para = strings.TrimSpace(para)
		if para == "" {
			b.WriteString("\n")
			continue
		}
		for {
			if len(para) <= width {
				b.WriteString(para)
				b.WriteString("\n")
				break
			}
			i := strings.LastIndex(para[:width+1], " ")
			if i <= 0 {
				i = width
			}
			b.WriteString(para[:i])
			b.WriteString("\n")
			para = strings.TrimSpace(para[i:])
		}
	}
	return strings.TrimSuffix(b.String(), "\n")
}

func commandHelpFunc(cmd *cobra.Command, args []string) {
	out := cmd.OutOrStderr()
	if cmd.Long != "" {
		fmt.Fprintf(out, "%s\n\n", wrapWords(cmd.Long, 72))
	}
	fmt.Fprintf(out, "Usage:\n  %s\n\n", cmd.UseLine())
	fmt.Fprintln(out, "Output:")
	for _, name := range outputFlagOrder {
		if f := cmd.Flags().Lookup(name); f != nil {
			printFlagUsage(out, f)
		}
	}
	if version != "" {
		fmt.Fprintf(out, "\nVersion: %s\n", version)
	}
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

func rootRunE(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		rootHelpFunc(cmd, nil)
		return nil
	}
	log.SetLevel(log.Normal)
	hadIssues, err := run(true, false, args)
	if err != nil {
		return err
	}
	checkHadIssues = hadIssues
	return nil
}

func checkRunE(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		commandHelpFunc(cmd, nil)
		return nil
	}
	log.SetLevel(outputLevelFromCmd(cmd))
	hadIssues, err := run(true, false, args)
	if err != nil {
		return err
	}
	checkHadIssues = hadIssues
	return nil
}

func writeRunE(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		commandHelpFunc(cmd, nil)
		return nil
	}
	log.SetLevel(outputLevelFromCmd(cmd))
	_, err := run(false, true, args)
	return err
}

func run(check, doWrite bool, paths []string) (hadIssues bool, err error) {
	start := time.Now()
	lvl := log.GetLevel()
	if lvl >= log.Verbose {
		log.Logf(log.Verbose, "Configuration: check=%v paths=%v\n", check, paths)
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
		for _, p := range sortedKeys(skipped) {
			log.Logf(log.Verbose, "scanner: rejected %s (reason: %s)\n", p, skipped[p])
		}
		for _, p := range files {
			log.Logf(log.Verbose, "scanner: accepted %s\n", p)
		}
	}
	if len(files) == 0 {
		if lvl >= log.Normal {
			fmt.Fprintln(os.Stdout, "No text files found.")
			if check {
				report.Write(os.Stdout, report.FormatCompact, nil, 0, nil)
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
				log.Logf(log.Verbose, "Writing %s\n", path)
			}
		}
		issues, err := rules.CheckFile(path)
		if err != nil {
			return false, err
		}
		if len(issues) > 0 {
			fileIssues[path] = issues
			allIssues = append(allIssues, issues...)
			if lvl >= log.Verbose {
				ruleIDs := make(map[string]bool)
				for _, i := range issues {
					ruleIDs[i.RuleID] = true
				}
				var ids []string
				for id := range ruleIDs {
					ids = append(ids, id)
				}
				sort.Strings(ids)
				log.Logf(log.Verbose, "rules: %s -> %d issue(s): %s\n", path, len(issues), strings.Join(ids, ", "))
			}
		}
	}
	if check {
		if lvl >= log.Normal {
			if err := report.Write(os.Stdout, report.FormatCompact, allIssues, len(files), files); err != nil {
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
		if lvl >= log.Verbose {
			log.Logf(log.Verbose, "write: applied to %s\n", path)
		}
	}
	if lvl >= log.Normal && len(fileIssues) > 0 {
		paths := make([]string, 0, len(fileIssues))
		for p := range fileIssues {
			paths = append(paths, p)
		}
		sort.Strings(paths)
		fmt.Fprintf(os.Stdout, "Wrote %d file(s):\n", len(paths))
		for _, p := range paths {
			fmt.Fprintln(os.Stdout, p)
		}
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

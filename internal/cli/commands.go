package cli

import (
	"fmt"
	"strings"

	"github.com/axia/axia-cli/internal/actions"
	"github.com/spf13/cobra"
)

// GetMapCmd returns the map subcommand
func GetMapCmd() *cobra.Command {
	var format string
	
	cmd := &cobra.Command{
		Use:   "map",
		Short: "Generate or display trust map",
		RunE: func(cmd *cobra.Command, args []string) error {
			graph := actions.NewTrustGraph()
			
			var result string
			var err error
			
			switch format {
			case "dot":
				result, err = graph.MapDOT()
			case "ascii":
				result, err = graph.MapASCII()
			default:
				result, err = graph.Map()
			}
			
			if err != nil {
				return fmt.Errorf("failed to generate trust map: %w", err)
			}
			
			fmt.Println(result)
			return nil
		},
	}
	
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, dot, ascii)")
	return cmd
}

// GetGetCmd returns the get subcommand
func GetGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [id]",
		Short: "Get trust information for an identity",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			graph := actions.NewTrustGraph()
			
			claims, err := graph.Get(id)
			if err != nil {
				return fmt.Errorf("failed to get trust information: %w", err)
			}

			if len(claims) == 0 {
				fmt.Printf("No trust claims found for '%s'\n", id)
				return nil
			}

			fmt.Printf("Trust claims for '%s':\n", id)
			for _, claim := range claims {
				if claim.Subject == id {
					fmt.Printf("  → Trusts '%s' as '%s'\n", claim.Object, claim.Predicate)
				} else {
					fmt.Printf("  ← Trusted by '%s' as '%s'\n", claim.Subject, claim.Predicate)
				}
			}
			return nil
		},
	}
	return cmd
}

// GetClaimCmd returns the claim subcommand
func GetClaimCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim [subject] [predicate] [object]",
		Short: "Make a trust claim",
		Long:  `Create a new trust claim stating that [subject] trusts [object] as [predicate]`,
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			subject := strings.TrimSpace(args[0])
			predicate := strings.TrimSpace(args[1])
			object := strings.TrimSpace(args[2])

			if subject == "" || predicate == "" || object == "" {
				return fmt.Errorf("subject, predicate, and object cannot be empty")
			}

			graph := actions.NewTrustGraph()
			err := graph.Claim(subject, predicate, object)
			if err != nil {
				return fmt.Errorf("failed to create trust claim: %w", err)
			}

			fmt.Printf("Created trust claim: %s -[%s]-> %s\n", subject, predicate, object)
			return nil
		},
	}
	return cmd
}

// GetSearchCmd returns the search subcommand
func GetSearchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search [query]",
		Short: "Search trust claims",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := args[0]
			graph := actions.NewTrustGraph()
			
			results, err := graph.Search(query)
			if err != nil {
				return fmt.Errorf("search failed: %w", err)
			}
			
			if len(results) == 0 {
				fmt.Printf("No results found for query: %s\n", query)
				return nil
			}
			
			fmt.Printf("Found %d results for query: %s\n", len(results), query)
			for _, claim := range results {
				fmt.Printf("  %s -[%s]-> %s\n", claim.Subject, claim.Predicate, claim.Object)
			}
			return nil
		},
	}
	return cmd
}

// GetStatsCmd returns the stats subcommand
func GetStatsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Show trust graph statistics",
		RunE: func(cmd *cobra.Command, args []string) error {
			graph := actions.NewTrustGraph()
			stats := graph.Stats()
			
			fmt.Printf("Trust Graph Statistics:\n")
			fmt.Printf("  Total Claims: %d\n", stats["total_claims"])
			fmt.Printf("  Unique Entities: %d\n", stats["unique_entities"])
			fmt.Printf("  Unique Predicates: %d\n", stats["unique_predicates"])
			
			fmt.Printf("\nTop Predicates:\n")
			for _, pc := range stats["top_predicates"].([]actions.predCount) {
				fmt.Printf("  %s: %d claims\n", pc.Predicate, pc.Count)
			}
			
			return nil
		},
	}
	return cmd
} 
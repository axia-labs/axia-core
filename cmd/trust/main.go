package main

import (
	"github.com/spf13/cobra"
	"axia/internal/axiom"
	"axia/internal/logging"
	"axia/internal/trust"
	"context"
	"os"
	"os/signal"
	"syscall"
	"axia/internal/server"
	"axia/internal/database"
	"fmt"
	"time"
	"github.com/google/uuid"
	"axia/internal/storage/ipfs"
	"axia/internal/auth"
)

func main() {
	logger := logging.NewLogger()

	// Check for secret key
	if os.Getenv("AXIA_SECRET_KEY") == "" {
		logger.Fatal("AXIA_SECRET_KEY environment variable must be set")
	}

	auth, err := auth.NewAuthenticator(logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize authenticator")
	}

	// Add authentication context to all operations
	ctx := context.WithValue(context.Background(), "secret_key", os.Getenv("AXIA_SECRET_KEY"))

	// Initialize database with auth context
	db, err := database.New(database.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     5432,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
	}, logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to database")
	}
	defer db.Close()

	// Pass authenticated context to components
	network := trust.NewNetwork(logger, db, auth)
	manager := axiom.NewManager(logger, db, auth)

	var rootCmd = &cobra.Command{
		Use:   "axios",
		Short: "Axiomatic Trust Graph CLI",
	}

	var claimCmd = &cobra.Command{
		Use:   "claim",
		Short: "Create an axiomatic claim",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := auth.ValidateContext(ctx); err != nil {
				return fmt.Errorf("authentication failed: %w", err)
			}
			agent, _ := cmd.Flags().GetString("agent")
			subject, _ := cmd.Flags().GetString("subject")
			axiomText, _ := cmd.Flags().GetString("axiom")
			confidence, _ := cmd.Flags().GetFloat64("confidence")
			tags, _ := cmd.Flags().GetStringSlice("tags")

			claim, err := manager.CreateClaim(agent, subject, axiomText, confidence, tags)
			if err != nil {
				return err
			}

			return network.AddClaim(claim)
		},
	}

	var truthCmd = &cobra.Command{
		Use:   "truth",
		Short: "Query the trust network",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := trust.QueryOptions{
				Observer:      cmd.Flag("observer").Value.String(),
				Agent:        cmd.Flag("agent").Value.String(),
				Subject:      cmd.Flag("subject").Value.String(),
				Tags:         cmd.Flag("tags").Value.String(),
				Depth:        cmd.Flag("depth").Value.String(),
				MinConfidence: cmd.Flag("min-confidence").Value.String(),
				MaxConfidence: cmd.Flag("max-confidence").Value.String(),
				UseConsensus: cmd.Flag("consensus").Value.String() == "true",
				UseTrustDecay: cmd.Flag("decay").Value.String() == "true",
			}

			results, err := network.Query(opts)
			if err != nil {
				return err
			}

			// Format and display results
			return nil
		},
	}

	var serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Start the Axia webhook server",
		RunE: func(cmd *cobra.Command, args []string) error {
			port, _ := cmd.Flags().GetInt("port")
			srv := server.NewServer(port, manager, network, logger)

			// Handle graceful shutdown
			done := make(chan os.Signal, 1)
			signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

			go func() {
				if err := srv.Start(); err != nil && err != http.ErrServerClosed {
					logger.WithError(err).Fatal("Server failed")
				}
			}()

			logger.Info("Server started")

			<-done
			logger.Info("Shutting down server")

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := srv.Shutdown(ctx); err != nil {
				logger.WithError(err).Error("Server shutdown failed")
				return err
			}

			logger.Info("Server stopped")
			return nil
		},
	}

	var migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			schemaFile := "internal/database/schema.sql"
			schema, err := os.ReadFile(schemaFile)
			if err != nil {
				return fmt.Errorf("failed to read schema file: %w", err)
			}

			_, err = db.pool.Exec(context.Background(), string(schema))
			if err != nil {
				return fmt.Errorf("failed to execute schema: %w", err)
			}

			logger.Info("Database migration completed successfully")
			return nil
		},
	}

	var ipfsCmd = &cobra.Command{
		Use:   "ipfs",
		Short: "Manage IPFS storage of trust graphs",
	}

	var uploadCmd = &cobra.Command{
		Use:   "upload",
		Short: "Upload trust graph to IPFS",
		RunE: func(cmd *cobra.Command, args []string) error {
			filters, _ := cmd.Flags().GetStringToString("filter")
			
			// Get claims based on filters
			claims, err := db.QueryClaims(context.Background(), filters)
			if err != nil {
				return fmt.Errorf("failed to query claims: %w", err)
			}

			// Create graph representation
			graphData := map[string]interface{}{
				"claims": claims,
				"metadata": map[string]interface{}{
					"timestamp": time.Now(),
					"version":   "1.0",
				},
			}

			// Upload to IPFS
			ipfsClient := ipfs.NewTatumClient(os.Getenv("TATUM_API_KEY"), logger)
			ipfsID, err := ipfsClient.UploadGraph(context.Background(), graphData)
			if err != nil {
				return fmt.Errorf("failed to upload to IPFS: %w", err)
			}

			// Store IPFS record
			record := &database.IPFSRecord{
				ID:        uuid.New(),
				IPFSID:    ipfsID,
				Type:      "trust_graph",
				Metadata:  filters,
				CreatedAt: time.Now(),
			}

			if err := db.StoreIPFSRecord(context.Background(), record); err != nil {
				return fmt.Errorf("failed to store IPFS record: %w", err)
			}

			fmt.Printf("Successfully uploaded trust graph to IPFS: %s\n", ipfsID)
			return nil
		},
	}

	var getCmd = &cobra.Command{
		Use:   "get [ipfs-id]",
		Short: "Get trust graph from IPFS",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ipfsID := args[0]
			
			ipfsClient := ipfs.NewTatumClient(os.Getenv("TATUM_API_KEY"), logger)
			data, err := ipfsClient.GetGraph(context.Background(), ipfsID)
			if err != nil {
				return fmt.Errorf("failed to get graph from IPFS: %w", err)
			}

			fmt.Println(string(data))
			return nil
		},
	}

	// Add flags
	claimCmd.Flags().String("agent", "", "DID or URL of AI agent making the claim")
	claimCmd.Flags().String("subject", "", "DID or URL of claim subject")
	claimCmd.Flags().String("axiom", "", "Axiomatic statement being claimed")
	claimCmd.Flags().Float64("confidence", 0.0, "Confidence score in range 0..1")
	claimCmd.Flags().StringSlice("tags", []string{}, "Categorical tags for the claim")

	truthCmd.Flags().String("observer", "", "Observer agent's perspective")
	truthCmd.Flags().String("agent", "", "Filter by claim-making agent")
	truthCmd.Flags().String("subject", "", "Filter by claim subject")
	truthCmd.Flags().StringSlice("tags", []string{}, "Filter by categorical tags")
	truthCmd.Flags().Int("depth", 3, "Search depth in trust network")
	truthCmd.Flags().Float64("min-confidence", 0.0, "Minimum confidence threshold")
	truthCmd.Flags().Float64("max-confidence", 1.0, "Maximum confidence threshold")
	truthCmd.Flags().Bool("consensus", false, "Generate consensus analysis")
	truthCmd.Flags().Bool("decay", false, "Trust decay with network distance")

	serverCmd.Flags().Int("port", 8080, "Port to run the server on")

	uploadCmd.Flags().StringToString("filter", nil, "Filters for claims to include in graph")
	ipfsCmd.AddCommand(uploadCmd, getCmd)

	rootCmd.AddCommand(claimCmd, truthCmd, serverCmd, migrateCmd, ipfsCmd)
	rootCmd.Execute()
} 
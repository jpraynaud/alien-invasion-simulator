package cmd

import (
	"context"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/jpraynaud/alien-invasion-simulator/pkg/simulator"
)

// rootCmd represents the base command when called without any subcommands
var (
	// Flags
	totalAliens uint
	maxSteps    uint
	mapFilepath string

	// Commands
	rootCmd = &cobra.Command{
		Use:   "alien-invasion",
		Short: "An Alien Invasion Simulator",
		Long: `An Alien Invasion Simulator.
More informations available at: https://github.com/jpraynaud/alien-invasion-simulator`,
		RunE: func(cmd *cobra.Command, args []string) error {
			in, err := os.Open(mapFilepath)
			defer func() { _ = in.Close() }()
			if err != nil {
				return err
			}
			c := &config{
				totalAliens: totalAliens,
				maxSteps:    maxSteps,
				in:          in,
				out:         cmd.OutOrStdout(),
			}
			return runSimulator(cmd.Context(), c)
		},
	}
)

// Execute executes the command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Flag setup
	rootCmd.Flags().UintVarP(&totalAliens, "aliens", "n", 5, "total number of aliens")
	rootCmd.Flags().UintVarP(&maxSteps, "steps", "s", 10000, "maximum number of steps")
	rootCmd.Flags().StringVarP(&mapFilepath, "file", "m", "map.txt", "world map file path")
}

type dependencies struct {
	simulator simulator.Simulator
	world     simulator.WorldStorer
	random    simulator.Randomer
}

type config struct {
	totalAliens, maxSteps uint
	in                    io.ReadCloser
	out                   io.Writer
}

func initDependencies(c *config) (*dependencies, error) {
	deps := &dependencies{}
	deps.random = simulator.NewRandomSimple()
	deps.world = simulator.NewWorld()
	deps.simulator = simulator.NewSimulationEngine(
		c.totalAliens,
		c.maxSteps,
		deps.world,
		deps.random,
		c.in,
		c.out)
	return deps, nil
}

func runSimulator(ctx context.Context, c *config) error {
	//Init dependencies
	deps, err := initDependencies(c)
	if err != nil {
		log.WithError(err).Fatal("an error occurred on init")
	}

	// Run simulator
	return deps.simulator.Run(ctx)
}

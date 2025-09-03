/*
Copyright © 2025 EmilyxFox 48589793+EmilyxFox@users.noreply.github.com
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/emilyxfox/untear/util"
	cp "github.com/otiai10/copy"
	"github.com/spf13/cobra"
)

type uncoupledWorlds struct {
	world  []os.DirEntry
	nether []os.DirEntry
	end    []os.DirEntry
}

type worlds struct {
	world  os.DirEntry
	nether os.DirEntry
	end    os.DirEntry
}

var worldPrefix string
var verbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "untear [path]",
	Example: "untear server-files/ --prefix world",
	Short:   "Reassemble your Minecraft world which Paper tore apart",
	Long:    `Untear lets you rejoin the three Minecraft dimensions into a single world file so it usable with vanilla servers and in singleplayer.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			log.SetLevel(log.DebugLevel)
		}
		log.Debug("Setting world prefix.", "value", worldPrefix)

		if len(args) > 1 {
			log.Fatal("More than one argument provided. Please provide a (relative) path to a directory containing your torn world folders, or none for current directory.")
		}

		var currentDir string
		if len(args) == 0 || args[0] == "." || args[0] == "./" {
			var err error
			currentDir, err = os.Getwd()
			if err != nil {
				log.Debug(err)
				log.Fatal("Could not get current directory.")
			}
		} else {
			var err error
			currentDir, err = util.ResolvePath(args[0])
			if err != nil {
				log.Debug(err)
				log.Fatal("Could not get specified directory: %v", args[0])
			}
		}

		log.Info("Searching for world folders.", "path", currentDir, "world_prefix", worldPrefix)

		worldName := worldPrefix
		netherName := fmt.Sprintf("%v_nether", worldPrefix)
		theEndName := fmt.Sprintf("%v_the_end", worldPrefix)

		discoveredWorlds := new(uncoupledWorlds)

		entries, err := os.ReadDir(currentDir)
		if err != nil {
			log.Debug(err)
			log.Fatal("Could not read files in current directory.")
		}

		for _, e := range entries {
			log.Debug("processing dir entry", "entry", e.Name())
			if !e.IsDir() {
				log.Debug("entry not a directory, skipping.", "file", e.Name())
				continue
			}
			switch e.Name() {
			case worldName:
				discoveredWorlds.world = append(discoveredWorlds.world, e)
			case netherName:
				discoveredWorlds.nether = append(discoveredWorlds.nether, e)
			case theEndName:
				discoveredWorlds.end = append(discoveredWorlds.end, e)
			}
		}

		if len(discoveredWorlds.world) > 1 || len(discoveredWorlds.nether) > 1 || len(discoveredWorlds.end) > 1 {
			log.Fatal("Too many world directories found... How did this happen?")
		}
		switch 0 {
		case len(discoveredWorlds.world):
			log.Fatalf("No world directory found. (looking for `%v`)", worldName)
		case len(discoveredWorlds.nether):
			log.Fatalf("No nether directory found. (looking for `%v`)", netherName)
		case len(discoveredWorlds.end):
			log.Fatalf("No the_end directory found. (looking for `%v`)", theEndName)
		}

		paperWorlds := worlds{
			world:  discoveredWorlds.world[0],
			nether: discoveredWorlds.nether[0],
			end:    discoveredWorlds.end[0],
		}

		log.Info("Found world directories.",
			"world", filepath.Join(currentDir, paperWorlds.world.Name()),
			"nether", filepath.Join(currentDir, paperWorlds.nether.Name()),
			"the_end", filepath.Join(currentDir, paperWorlds.end.Name()),
		)

		vanillaWorldDirPath := filepath.Join(currentDir, fmt.Sprintf("vanilla_%v", worldPrefix))

		log.Info("Creating vanilla world directory", "dir", vanillaWorldDirPath)
		if err = os.Mkdir(vanillaWorldDirPath, os.ModePerm); err != nil {
			log.Fatalf("failed to create directory 'vanilla_%s': %v", worldPrefix, err)
		}

		log.Info("Copying overworld into vanilla world directory...", "from", filepath.Join(currentDir, paperWorlds.world.Name()), "to", vanillaWorldDirPath)
		err = cp.Copy(filepath.Join(currentDir, paperWorlds.world.Name()), vanillaWorldDirPath)
		if err != nil {
			log.Fatal("failed to copy paper overworld to vanilla world directory")
		}

		log.Info("Copying nether into vanilla world directory...", "from", fmt.Sprintf("%v/DIM-1", paperWorlds.nether.Name()), "to", fmt.Sprintf("%v/DIM-1", fmt.Sprintf("vanilla_%v", worldPrefix)))
		err = cp.Copy(filepath.Join(currentDir, paperWorlds.nether.Name(), "DIM-1"), filepath.Join(vanillaWorldDirPath, "DIM-1"))

		log.Info("Copying the_end into vanilla world directory...", "from", fmt.Sprintf("%v/DIM1", paperWorlds.end.Name()), "to", fmt.Sprintf("%v/DIM1", fmt.Sprintf("vanilla_%v", worldPrefix)))
		err = cp.Copy(filepath.Join(currentDir, paperWorlds.end.Name(), "DIM1"), filepath.Join(vanillaWorldDirPath, "DIM1"))

		log.Info("Removing `paper-world` file...", "from", fmt.Sprintf("%v/paper-world.yml", fmt.Sprintf("vanilla_%v", worldPrefix)))
		err = os.Remove(filepath.Join(vanillaWorldDirPath, "paper-world.yml"))
		if err != nil {
			log.Warnf("failed to remove paper-world file. %v", err)
		}

		log.Info("🎉 Success!")
		log.Info("🟩 Worlds have been merged")
		log.Warn("⚠️ DO NOT DELETE PAPER WORLDS BEFORE YOU HAVE CHECKED ALL DIMENSIONS IN THE MERGED WORLD!")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable debug logging (default is `false`)")
	rootCmd.PersistentFlags().StringVarP(&worldPrefix, "prefix", "p", "world", "world prefix (default is `world`)")
}

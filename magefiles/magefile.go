// ⚡ Core Mage Tasks.
package main

import (
	"os"

	"github.com/DelineaXPM/terraform-provider-dsv/v2/magefiles/constants"

	"github.com/magefile/mage/mg"
	"github.com/pterm/pterm"
	"github.com/sheldonhull/magetools/ci"
	"github.com/sheldonhull/magetools/tooling"

	// mage:import
	"github.com/sheldonhull/magetools/gotools"
)

// createDirectories creates the local working directories for build artifacts and tooling.
func createDirectories() error {
	for _, dir := range []string{constants.ArtifactDirectory, constants.CacheDirectory} {
		if err := os.MkdirAll(dir, constants.PermissionUserReadWriteExecute); err != nil {
			pterm.Error.Printf("failed to create dir: [%s] with error: %v\n", dir, err)

			return err
		}
		pterm.Success.Printf("✅ [%s] dir created\n", dir)
	}

	return nil
}

// Init runs multiple tasks to initialize all the requirements for running a project for a new contributor.
func Init() error { //nolint:deadcode // Not dead, it's alive.
	pterm.DefaultHeader.Println("running Init()")

	pterm.Info.Println("Mod Download")
	mg.Deps(
		gotools.Go{}.Tidy,
	)

	pterm.Info.Println("Installing Core CI Dependencies")
	if err := tooling.SilentInstallTools([]string{
		"github.com/goreleaser/goreleaser@latest",
		"github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest",
		// "github.com/release-lab/whatchanged/cmd/whatchanged@latest",
	}); err != nil {
		return err
	}
	if ci.IsCI() {
		pterm.DefaultHeader.Println("CI detected, minimal init was applied, finished")
		return nil
	}
	mg.SerialDeps(
		Clean,
		createDirectories,
		(gotools.Go{}.Tidy),
		(gotools.Go{}.Init),
	)

	if ci.IsCI() {
		pterm.Debug.Println("CI detected, done with init")
		return nil
	}

	pterm.DefaultSection.Println("Setup Project Specific Tools")
	if err := tooling.SilentInstallTools(ToolList); err != nil {
		return err
	}
	// These can run in parallel as different toolchains.
	mg.Deps()
	return nil
}

// Clean up after yourself.
func Clean() {
	pterm.Success.Println("Cleaning...")
	for _, dir := range []string{constants.ArtifactDirectory, constants.CacheDirectory} {
		err := os.RemoveAll(dir)
		if err != nil {
			pterm.Error.Printf("failed to removeall: [%s] with error: %v\n", dir, err)
		}
		pterm.Success.Printf("🧹 [%s] dir removed\n", dir)
	}
	mg.Deps(createDirectories)
}

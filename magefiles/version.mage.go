package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"
	mtu "github.com/sheldonhull/magetools/pkg/magetoolsutils"
	"github.com/sheldonhull/magetools/pkg/req"
)

// Svu uses svu cli for tasks related to semver versioning.
type Svu mg.Namespace

// ⬆️ Bump increases the tag by a patched increment and pushes to the origin repo (no v prefix).
func (Svu) Bump() error {
	mtu.CheckPtermDebug()
	var err error
	binary, err := req.ResolveBinaryByInstall("svu", "github.com/caarlos0/svu")
	if err != nil {
		return err
	}
	_, tagPrefix, err := Svu{}.GetVersion()
	if err != nil {
		return err
	}

	if tagPrefix == "" {
		pterm.Warning.Println("no tag prefix, so checking for APP environment variable to use as a prefix")
		envTagPrefix := os.Getenv("TAG_PREFIX")
		if envTagPrefix != "" {
			pterm.Success.Printfln("TAG_PREFIX env variable found and setting to tagPrefix: %s", envTagPrefix)
			if !strings.Contains(envTagPrefix, "/") {
				return fmt.Errorf(
					"TAG_PREFIX env variable doesn't have trailing slash. This is key for parsing tag prefixes, ensure TAG_PREFIX has trailing slash in inputs",
				)
			}
			tagPrefix = envTagPrefix
		}
	}

	svuArgs := []string{"patch"}
	if tagPrefix != "" {
		pterm.Debug.Printfln("tagPrefix not empty, so setting pattern and prefix args for svu")
		svuArgs = append(svuArgs, "--pattern", tagPrefix+"*")
		svuArgs = append(svuArgs, "--prefix", tagPrefix)
	} else {
		pterm.Debug.Printfln("tagPrefix was empty, so stripping default 'v' prefix")
		svuArgs = append(svuArgs, "--strip-prefix", "--tag-mode", "current-branch")
	}

	newVersion, err := sh.Output(binary, svuArgs...)
	if err != nil {
		return err
	}
	pterm.Info.Printfln("out: %s", newVersion)

	newTag, _, err := Svu{}.GetVersion()
	if err != nil {
		return err
	}

	// Extra safety net to validate that the calculated semver can be correctly parsed as a Semver.
	versionReturnedByTool, err := semver.NewVersion(newTag)
	if err != nil {
		pterm.Warning.Printfln(
			"(nonterminating)  semver.NewVersion(newVersion) can't parse into semver version type: %v from: %v",
			err,
			newVersion,
		)
		return fmt.Errorf(
			"semver.NewVersion(newVersion) can't parse into semver version type: %w from: %v",
			err,
			newVersion,
		)
	}

	// Azure DevOps Build Variables
	// TODO: add GitHub Action ones
	pterm.Info.Printfln("versionReturnedByTool: %s", versionReturnedByTool.String())
	metadata := []string{}
	reason, isSet := os.LookupEnv("BUILD_REASON")
	if isSet {
		pterm.Debug.Printfln("os.LookupEnv(\"BUILD_REASON\"): %s", reason)
		metadata = append(metadata, "BuildReason: "+reason)
	}
	// Azure DevOps Build Variables
	// TODO: add GitHub Action ones
	requestedFor, isSet := os.LookupEnv("BUILD_REQUESTEDFOR")
	if isSet {
		pterm.Debug.Printfln("os.LookupEnv(\"BUILD_REQUESTEDFOR\"): %s", requestedFor)
		metadata = append(metadata, "BuildRequestedFor: "+requestedFor)
	}
	// Azure DevOps Build Variables
	// TODO: add GitHub Action ones
	// If a commitsha is provided by CI system, let's tag this specifically to be very specific on our tagging.
	currentSHA := os.Getenv("BUILD_SOURCEVERSION")
	var commitSHAToTag string
	if currentSHA != "" {
		commitSHAToTag = currentSHA
	}
	tagArgs := []string{}
	tagArgs = append(tagArgs, []string{
		"tag",
		"--force",
		"-a",
		newVersion,
	}...)
	if currentSHA != "" {
		tagArgs = append(tagArgs, commitSHAToTag)
	}
	tagArgs = append(tagArgs, []string{
		"-m",
		fmt.Sprintf("[ci] %s", strings.Join(metadata, " ")),
	}...)
	// --force is used to ensure no fatal error due to tag existing already. Just forcibly replace for now.
	_, err = sh.Output("git", tagArgs...)

	if err != nil {
		pterm.Warning.Println("failed to forciby tag with new version")
		return err
	}

	pterm.Info.Println("attempting to push to origin now")
	if err := sh.RunV("git", "push", "origin", newVersion); err != nil {
		return fmt.Errorf("unable to push tag for commit: %w", err)
	}

	pterm.Success.Println("(Svu) Bump()")
	return nil
}

// GetVersion returns the current version found by Svu for usage in other commands (no v prefix).
func (Svu) GetVersion() (semver, prefix string, err error) {
	mtu.CheckPtermDebug()

	pterm.Debug.Println("attempt to evaluate based TAG_PREFIX in environment")
	prefix = os.Getenv("TAG_PREFIX")
	if prefix == "" {
		return "", "", fmt.Errorf("TAG_PREFIX is required for versioning in this monorepo")
	}

	pterm.Success.Printfln("TAG_PREFIX env variable found and setting to prefix: %s", prefix)
	if !strings.Contains(prefix, "/") {
		return "", "", fmt.Errorf(
			"TAG_PREFIX env variable doesn't have trailing slash. This is key for parsing tag prefixes, ensure TAG_PREFIX has trailing slash in inputs",
		)
	}

	binary, err := req.ResolveBinaryByInstall("svu", "github.com/caarlos0/svu")
	if err != nil {
		return "", "", err
	}
	svuArgs := []string{"current"}
	if prefix != "" {
		pterm.Debug.Printfln("prefix not empty, so setting pattern and prefix args for svu")
		svuArgs = append(svuArgs, "--pattern", prefix+"*")
		svuArgs = append(svuArgs, "--prefix", prefix)
		svuArgs = append(svuArgs, "--strip-prefix") // To allow us to preturn just the version for goreleaser and tools.
	} else {
		return "", "", fmt.Errorf("prefix was empty, so unable to proceed: %w", err)
	}
	var trimmedSemverOnly string

	semverTag, err := sh.Output(binary, svuArgs...)
	trimmedSemverOnly = strings.TrimSpace(semverTag)
	pterm.Debug.Printfln("GetVersion svu %+v: %s", svuArgs, semverTag)
	if err != nil {
		return "", "", err
	}
	pterm.Info.Printfln("svu current: %s", trimmedSemverOnly)
	return trimmedSemverOnly, prefix, nil
}

// Predict returns the best estimated semver expected next.
func (Svu) Predict() error {
	mtu.CheckPtermDebug()
	binary, err := req.ResolveBinaryByInstall("svu", "github.com/caarlos0/svu")
	if err != nil {
		return fmt.Errorf("unable to setup svu: %w", err)
	}

	out, err := sh.Output(binary, "next", "--strip-prefix", "--tag-mode=current-branch", "--force-patch-increment")
	if err != nil {
		return err
	}
	cleaner := strings.TrimSpace(out)
	pterm.Info.Printfln("svu predicts: %s", cleaner)
	return nil
}

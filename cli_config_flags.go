package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type configFlagOverrides struct {
	downloadDir             string
	camHost                 string
	camUser                 string
	camPass                 string
	camInsecureSkipVerify   bool
	scanLastDays            int
	scanFromLocalLatestDays int
	httpRetryCnt            int
	httpTimeout             int
	noProxy                 bool
}

func addConfigFlags(cmd *cobra.Command, overrides *configFlagOverrides) {
	flags := cmd.Flags()

	flags.StringVarP(&overrides.downloadDir, "download-dir", "d", "", "Override DOWNLOAD_DIR")
	flags.StringVarP(&overrides.camHost, "cam-host", "H", "", "Override CAM_HOST")
	flags.StringVarP(&overrides.camUser, "cam-user", "u", "", "Override CAM_USER")
	flags.StringVarP(&overrides.camPass, "cam-pass", "p", "", "Override CAM_PASS")
	flags.BoolVarP(&overrides.camInsecureSkipVerify, "cam-insecure-skip-verify", "k", false, "Override CAM_INSECURE_SKIP_VERIFY")
	flags.IntVarP(&overrides.scanLastDays, "scan-last-days", "s", 0, "Override SCAN_LAST_DAYS")
	flags.IntVarP(&overrides.scanFromLocalLatestDays, "scan-from-local-latest", "l", 0, "Override SCAN_FROM_LOCAL_LATEST")
	flags.IntVarP(&overrides.httpRetryCnt, "http-retry-cnt", "r", 0, "Override HTTP_RETRY_CNT")
	flags.IntVarP(&overrides.httpTimeout, "http-timeout", "t", 0, "Override HTTP_TIMEOUT")
	flags.BoolVarP(&overrides.noProxy, "no-proxy", "P", false, "Override NO_PROXY")
}

func applyConfigFlagOverrides(flagSet *pflag.FlagSet, overrides *configFlagOverrides) error {
	if err := setStringEnvIfChanged(flagSet, "download-dir", "DOWNLOAD_DIR", overrides.downloadDir); err != nil {
		return err
	}
	if err := setStringEnvIfChanged(flagSet, "cam-host", "CAM_HOST", overrides.camHost); err != nil {
		return err
	}
	if err := setStringEnvIfChanged(flagSet, "cam-user", "CAM_USER", overrides.camUser); err != nil {
		return err
	}
	if err := setStringEnvIfChanged(flagSet, "cam-pass", "CAM_PASS", overrides.camPass); err != nil {
		return err
	}
	if err := setBoolEnvIfChanged(flagSet, "cam-insecure-skip-verify", "CAM_INSECURE_SKIP_VERIFY", overrides.camInsecureSkipVerify); err != nil {
		return err
	}
	if err := setIntEnvIfChanged(flagSet, "scan-last-days", "SCAN_LAST_DAYS", overrides.scanLastDays); err != nil {
		return err
	}
	if err := setIntEnvIfChanged(flagSet, "scan-from-local-latest", "SCAN_FROM_LOCAL_LATEST", overrides.scanFromLocalLatestDays); err != nil {
		return err
	}
	if err := setIntEnvIfChanged(flagSet, "http-retry-cnt", "HTTP_RETRY_CNT", overrides.httpRetryCnt); err != nil {
		return err
	}
	if err := setIntEnvIfChanged(flagSet, "http-timeout", "HTTP_TIMEOUT", overrides.httpTimeout); err != nil {
		return err
	}
	if err := setBoolEnvIfChanged(flagSet, "no-proxy", "NO_PROXY", overrides.noProxy); err != nil {
		return err
	}

	return nil
}

func setStringEnvIfChanged(flagSet *pflag.FlagSet, flagName, envKey, value string) error {
	if !flagSet.Changed(flagName) {
		return nil
	}

	return setEnvOverride(flagName, envKey, value)
}

func setBoolEnvIfChanged(flagSet *pflag.FlagSet, flagName, envKey string, value bool) error {
	if !flagSet.Changed(flagName) {
		return nil
	}

	return setEnvOverride(flagName, envKey, strconv.FormatBool(value))
}

func setIntEnvIfChanged(flagSet *pflag.FlagSet, flagName, envKey string, value int) error {
	if !flagSet.Changed(flagName) {
		return nil
	}

	return setEnvOverride(flagName, envKey, strconv.Itoa(value))
}

func setEnvOverride(flagName, envKey, value string) error {
	if err := os.Setenv(envKey, value); err != nil {
		return fmt.Errorf("override %s via --%s: %w", envKey, flagName, err)
	}

	return nil
}

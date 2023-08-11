package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration management utility",
	Long:  `Configuration management utility.`,
	Args:  cobra.NoArgs,
}

func addConfigFlags(flags *pflag.FlagSet) {
	addServerFlags(flags)
	flags.String("shell", "", "shell command to which other commands should be appended")

	flags.String("auth.header", "", "HTTP header for auth.method=proxy")
	flags.String("auth.command", "", "command for auth.method=hook")

	flags.String("recaptcha.host", "https://www.google.com", "use another host for ReCAPTCHA. recaptcha.net might be useful in China")
	flags.String("recaptcha.key", "", "ReCaptcha site key")
	flags.String("recaptcha.secret", "", "ReCaptcha secret")

	flags.String("branding.name", "", "replace 'File Browser' by this name")
	flags.String("branding.color", "", "set the theme color")
	flags.String("branding.files", "", "path to directory with images and custom styles")
	flags.Bool("branding.disableExternal", false, "disable external links such as GitHub links")
	flags.Bool("branding.disableUsedPercentage", false, "disable used disk percentage graph")
}

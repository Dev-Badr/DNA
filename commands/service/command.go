/*
 * Copyright (C) 2019 Skytells, Inc.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package service

import (
	"fmt"
	"os"
	"strings"

	"github.com/skytells-research/DNA/network/node/cmd"
	"github.com/skytells-research/DNA/network/node/cmd/commands/license"
	"github.com/skytells-research/DNA/network/node/core/service"
	"github.com/skytells-research/DNA/network/node/identity"
	identity_selector "github.com/skytells-research/DNA/network/node/identity/selector"
	"github.com/skytells-research/DNA/network/node/metadata"
	openvpn_service "github.com/skytells-research/DNA/network/node/services/openvpn/service"
	wireguard_service "github.com/skytells-research/DNA/network/node/services/wireguard/service"
	"github.com/skytells-research/DNA/network/node/tequilapi/client"
	"github.com/urfave/cli"
)

const serviceCommandName = "service"

var (
	identityFlag = cli.StringFlag{
		Name:  "identity",
		Usage: "Keystore's identity used to provide service. If not given identity will be created automatically",
		Value: "",
	}
	identityPassphraseFlag = cli.StringFlag{
		Name:  "identity.passphrase",
		Usage: "Used to unlock keystore's identity",
		Value: "",
	}

	agreedTermsConditionsFlag = cli.BoolFlag{
		Name:  "agreed-terms-and-conditions",
		Usage: "Agree with terms & conditions",
	}
)

// NewCommand function creates service command
func NewCommand(licenseCommandName string) *cli.Command {
	var di cmd.Dependencies
	command := &cli.Command{
		Name:      serviceCommandName,
		Usage:     "Starts and publishes services on sdna Network",
		ArgsUsage: "comma separated list of services to start",
		Action: func(ctx *cli.Context) error {
			if !ctx.Bool(agreedTermsConditionsFlag.Name) {
				printTermWarning(licenseCommandName)
				os.Exit(2)
			}

			errorChannel := make(chan error)
			nodeOptions := cmd.ParseFlagsNode(ctx)
			if err := di.Bootstrap(nodeOptions); err != nil {
				return err
			}
			go func() { errorChannel <- di.Node.Wait() }()

			cmd.RegisterSignalCallback(func() { errorChannel <- nil })

			cmdService := &serviceCommand{
				tequilapi:    client.NewClient(nodeOptions.TequilapiAddress, nodeOptions.TequilapiPort),
				errorChannel: errorChannel,
				identityHandler: identity_selector.NewHandler(
					di.IdentityManager,
					di.sdnaAPI,
					identity.NewIdentityCache(nodeOptions.Directories.Keystore, "remember.json"),
					di.SignerFactory),
			}

			go func() {
				errorChannel <- cmdService.Run(ctx)
			}()

			return <-errorChannel
		},
		After: func(ctx *cli.Context) error {
			return di.Shutdown()
		},
	}

	registerFlags(&command.Flags)

	return command
}

// serviceCommand represent entrypoint for service command with top level components
type serviceCommand struct {
	identityHandler identity_selector.Handler
	tequilapi       *client.Client
	errorChannel    chan error
}

// Run runs a command
func (sc *serviceCommand) Run(ctx *cli.Context) (err error) {
	arg := ctx.Args().Get(0)
	if arg != "" {
		serviceTypes = strings.Split(arg, ",")
	}

	identity, err := sc.unlockIdentity(parseFlags(ctx))
	if err != nil {
		return err
	}

	if err := sc.runServices(ctx, identity.Address, serviceTypes); err != nil {
		return err
	}

	return <-sc.errorChannel
}

func (sc *serviceCommand) unlockIdentity(identityOptions service.OptionsIdentity) (identity.Identity, error) {
	loadIdentity := identity_selector.NewLoader(sc.identityHandler, identityOptions.Identity, identityOptions.Passphrase)
	return loadIdentity()
}

func (sc *serviceCommand) runServices(ctx *cli.Context, providerID string, serviceTypes []string) error {
	for _, serviceType := range serviceTypes {
		options, err := parseFlagsByServiceType(ctx, serviceType)
		if err != nil {
			return err
		}
		go sc.runService(providerID, serviceType, options)
	}

	return nil
}

func (sc *serviceCommand) runService(providerID, serviceType string, options service.Options) {
	_, err := sc.tequilapi.ServiceStart(providerID, serviceType, options)
	if err != nil {
		sc.errorChannel <- err
	}
}

// registerFlags function register service flags to flag list
func registerFlags(flags *[]cli.Flag) {
	*flags = append(*flags,
		agreedTermsConditionsFlag,
		identityFlag, identityPassphraseFlag,
	)
	openvpn_service.RegisterFlags(flags)
	wireguard_service.RegisterFlags(flags)
}

// parseFlags function fills in service command options from CLI context
func parseFlags(ctx *cli.Context) service.OptionsIdentity {
	return service.OptionsIdentity{
		Identity:   ctx.String(identityFlag.Name),
		Passphrase: ctx.String(identityPassphraseFlag.Name),
	}
}

func parseFlagsByServiceType(ctx *cli.Context, serviceType string) (service.Options, error) {
	if f, ok := serviceTypesFlagsParser[serviceType]; ok {
		return f(ctx), nil
	}
	return service.OptionsIdentity{}, fmt.Errorf("unknown service type: %q", serviceType)
}

func printTermWarning(licenseCommandName string) {
	fmt.Println(metadata.VersionAsSummary(metadata.LicenseCopyright(
		"run program with 'myst "+licenseCommandName+" --"+license.LicenseWarrantyFlag.Name+"' option",
		"run program with 'myst "+licenseCommandName+" --"+license.LicenseConditionsFlag.Name+"' option",
	)))
	fmt.Println()

	fmt.Println("If you agree with these Terms & Conditions, run program again with '--agreed-terms-and-conditions' flag")
}

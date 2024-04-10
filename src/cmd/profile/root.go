package profile

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ephex2/go-gpt-cli/config"
	"github.com/ephex2/go-gpt-cli/config/profile"
	"github.com/ephex2/go-gpt-cli/log"

	"github.com/spf13/cobra"
)

var ProfileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Allows the creation, viewing, and modification of user profiles for each endpoint",
}

var readCmd = &cobra.Command{
	Use:               "read",
	Short:             "Reads the contents of a profile out to the terminal",
	Long:              "Reads the contents of a profile out to the terminal. When writing this to a file, it should be a valid configuration file to be used when performing update commands",
	Run:               profileReadCommandRun,
	Example:           "go-gpt-cli profile read endpointName profileName",
	Args:              cobra.ExactArgs(2),
	Aliases:           []string{"get"},
	ValidArgsFunction: validEndpointAndProfileArgs,
}

var createCmd = &cobra.Command{
	Use:               "create",
	Short:             "For a given endpoint, create a new profile",
	Run:               profileCreateCommandRun,
	Example:           "go-gpt-cli profile create endpointName profileName",
	Args:              cobra.ExactArgs(2),
	Aliases:           []string{"new"},
	ValidArgsFunction: validEndpointArgs, // don't autosuggest existing profiles
}

var updateCmd = &cobra.Command{
	Use:               "update",
	Short:             "On a given endpoint, update a profile from a provided configuration file",
	Long:              "On a given endpoint, update a profile from a provided configuration file. Note that given the way that profiles are configured the name of the profile should already be present from the config file itself",
	Run:               profileUpdateCommandRun,
	Example:           "go-gpt-cli profile update endpointName configFilePath",
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: validEndpointArgs,
}

var deleteCmd = &cobra.Command{
	Use:               "delete",
	Short:             "Within an endpoint, delete a specified profile",
	Run:               profileDeleteCommandRun,
	Example:           "go-gpt-cli profile delete endpointName profileName",
	Args:              cobra.ExactArgs(2),
	Aliases:           []string{"remove"},
	ValidArgsFunction: validEndpointAndProfileArgs,
}

var getAllCmd = &cobra.Command{
	Use:               "getall",
	Short:             "Get all profiles defined for an endpoint",
	Run:               profileGetAllCommandRun,
	Example:           "go-gpt-cli profile getall endpointName",
	Args:              cobra.ExactArgs(1),
	Aliases:           []string{"list"},
	ValidArgsFunction: validEndpointAndProfileArgs,
}

var endpointsCmd = &cobra.Command{
	Use:               "endpoints",
	Short:             "Get all endpoints that can use profiles",
	Run:               endpointsCommandRun,
	Example:           "go-gpt-cli profile endpoints",
	Args:              cobra.ExactArgs(0),
}

var defaultCmd = &cobra.Command{
	Use:               "default",
	Short:             "Sets the default profile for an endpoint",
	Run:               profileDefaultCommandRun,
	Example:           "go-gpt-cli profile default endpointName profileName",
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: validEndpointAndProfileArgs,
}

func init() {
	ProfileCmd.AddCommand(readCmd)
	ProfileCmd.AddCommand(createCmd)
	ProfileCmd.AddCommand(updateCmd)
	ProfileCmd.AddCommand(deleteCmd)
	ProfileCmd.AddCommand(getAllCmd)
	ProfileCmd.AddCommand(endpointsCmd)
	ProfileCmd.AddCommand(defaultCmd)

}

func Execute(cmd *cobra.Command, args []string) (err error) {

	if len(args) < 1 {
		cmd.Help()
		// err = errors.New(config.CmdHelpString)
	}

	err = cmd.Execute()

	return
}

func profileCreateCommandRun(cmd *cobra.Command, args []string) {
	e := getEndpoint(args[0])

	err := profile.RuntimeRepository.Create(e, args[1])
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}
}

func profileReadCommandRun(cmd *cobra.Command, args []string) {
	endpoint := getEndpoint(args[0])
	profileName := args[1]

	dummyP := endpoint.DefaultProfile() //.SetName(profileName)
	repo := dummyP.ProfileRepository()

	p, err := repo.Read(profileName, endpoint.Name())
	if err != nil {
		log.Critical("Error while trying to read profile: %s\n", err.Error())
		os.Exit(1)
	}

	formattedProfile, err := endpoint.ProfileFromJsonBuf(p)
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	formattedProfileBuf, err := json.MarshalIndent(formattedProfile, "", "    ")
	if err != nil {
		log.Critical("Error while trying to format profile: %s", err.Error())
		os.Exit(1)
	}

	fmt.Println(string(formattedProfileBuf))
}

func profileUpdateCommandRun(cmd *cobra.Command, args []string) {
	e := getEndpoint(args[0])

	newProfilePath := args[1]
	newProfileBytes, err := os.ReadFile(newProfilePath)
	if err != nil {
		log.Critical("Error while trying to read from the new profile path provided: %s\nError: %s\n", newProfilePath, err.Error())
		os.Exit(1)
	}

	p, err := e.ProfileFromJsonBuf(newProfileBytes)
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	repo := p.ProfileRepository()
	err = repo.Update(p)
	if err != nil {
		log.Critical("Error while updating profile config file: %s\n", err.Error())
		os.Exit(1)
	}
}

func profileDeleteCommandRun(cmd *cobra.Command, args []string) {
	endpointName := checkEndpoint(args[0])
	profileName := args[1]

	err := profile.RuntimeRepository.Delete(endpointName, profileName)
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}
}

func profileGetAllCommandRun(cmd *cobra.Command, args []string) {
	endpointName := checkEndpoint(args[0])

	names, err := profile.RuntimeRepository.GetAll(endpointName)
	if err != nil {
		log.Debug(err.Error() + "\n")
		return
	}

	for _, name := range names {
		fmt.Println(name)
	}
}

func endpointsCommandRun(cmd *cobra.Command, args []string) {
	names, err := profile.EndpointRegistry.List()
	if err != nil {
		log.Critical(err.Error() + "\n")
		return
	}

	for _, name := range names {
		fmt.Println(name)
	}
}

func profileDefaultCommandRun(cmd *cobra.Command, args []string) {
	endpointName := checkEndpoint(args[0])

	profileName := args[1]
	err := config.SetDefaultProfile(endpointName, profileName, true)
	if err != nil {
		log.Critical(err.Error() + "\n")
	}
}

// utility
func checkEndpoint(name string) string {
	endpointName := strings.ToLower(name)
	_, err := profile.EndpointRegistry.Get(endpointName)
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	return endpointName
}

func getEndpoint(name string) profile.Endpoint {
	endpointName := strings.ToLower(name)
	e, err := profile.EndpointRegistry.Get(endpointName)
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	return e
}

func validEndpointArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		endpoints, err := profile.EndpointRegistry.List()
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		return endpoints, cobra.ShellCompDirectiveNoFileComp
	}

	// "Else"
	return nil, cobra.ShellCompDirectiveDefault
}

func validEndpointAndProfileArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		endpoints, err := profile.EndpointRegistry.List()
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		return endpoints, cobra.ShellCompDirectiveNoFileComp
	} else if len(args) == 1 {
		names, err := profile.RuntimeRepository.GetAll(args[0])
		if err != nil {
			log.Debug(err.Error() + "\n")
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		return names, cobra.ShellCompDirectiveNoFileComp
	}

	// "Else"
	return nil, cobra.ShellCompDirectiveNoFileComp
}

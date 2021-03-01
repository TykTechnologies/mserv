package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/TykTechnologies/mserv/mservclient/client/mw"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:     "push",
	Short:   "Pushes a middleware to mserv",
	Long:    `Uploads a bundle file created with tyk CLI to mserv`,
	Example: `mservctl push /path/to/bundle.zip`,
	Args:    cobra.ExactArgs(1),
	Run:     pushMiddleware,
}

func init() {
	rootCmd.AddCommand(pushCmd)

	pushCmd.Flags().BoolP("storeonly", "s", false, "Don't process, just store it")
	pushCmd.Flags().StringP("apiid", "a", "", "Optional API ID")
}

func pushMiddleware(cmd *cobra.Command, args []string) {
	file, err := os.Open(args[0])
	if err != nil {
		log.WithError(err).Error("Couldn't open the bundle file")
		return
	}
	defer file.Close()

	apiID := cmd.Flag("apiid").Value.String()
	storeOnly := cmd.Flag("storeonly").Value.String() == "true"

	params := mw.NewMwAddParams().WithUploadFile(file).WithAPIID(&apiID).WithStoreOnly(&storeOnly)

	resp, err := mservapi.Mw.MwAdd(params, defaultAuth())
	if err != nil {
		log.WithError(err).Error("Couldn't push middleware")
		return
	}

	cmd.Printf("Middleware uploaded successfully, ID: %s\n", resp.GetPayload().Payload.BundleID)
}

package cmd

import (
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/TykTechnologies/mserv/mservclient/client/mw"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Updates a middleware on mserv",
	Long:    "Updates a middleware by ID",
	Example: `mservctl update 13b0eb10-419f-40ef-838d-6d26bb2eeaa8 /path/to/bundle.zip`,
	Args:    cobra.ExactArgs(2),
	Run:     updateMiddleware,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func updateMiddleware(cmd *cobra.Command, args []string) {
	file, err := os.Open(args[1])
	if err != nil {
		log.WithError(err).Error("Couldn't open the bundle file")
		return
	}
	defer file.Close()

	params := mw.NewMwUpdateParams().WithID(args[0]).WithUploadFile(file).WithTimeout(120 * time.Second)

	resp, err := mservapi.Mw.MwUpdate(params, defaultAuth())
	if err != nil {
		log.WithError(err).Error("Couldn't update middleware")
		return
	}

	cmd.Printf("Middleware uploaded successfully, ID: %s\n", resp.GetPayload().Payload.BundleID)
}

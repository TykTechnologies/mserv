package cmd

import (
	"github.com/TykTechnologies/mserv/mservclient/client/mw"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes a middleware from mserv",
	Long: `Deletes a middleware record by ID, e.g.:

$ mservctl delete 13b0eb10-419f-40ef-838d-6d26bb2eeaa8`,
	Args: cobra.ExactArgs(1),
	Run:  deleteMiddleware,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func deleteMiddleware(cmd *cobra.Command, args []string) {
	params := mw.NewMwDeleteParams().WithID(args[0])
	resp, err := mservapi.Mw.MwDelete(params, defaultAuth())
	if err != nil {
		log.WithError(err).Error("Couldn't delete middleware")
		return
	}

	cmd.Printf("Middleware deleted successfully, ID: %s\n", resp.GetPayload().Payload.BundleID)
}

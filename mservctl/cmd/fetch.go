package cmd

import (
	"fmt"

	coprocess "github.com/TykTechnologies/mserv/coprocess/bindings/go"
	"github.com/TykTechnologies/mserv/mservclient/client/mw"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetches a middleware record from mserv",
	Long: `Fetches a middleware record by ID, e.g.:

$ mservctl fetch 13b0eb10-419f-40ef-838d-6d26bb2eeaa8`,
	Args: cobra.ExactArgs(1),
	Run:  fetchMiddleware,
}

func init() {
	rootCmd.AddCommand(fetchCmd)

	fetchCmd.Flags().BoolP("functions", "f", false, "Show plugin functions")
}

func fetchMiddleware(cmd *cobra.Command, args []string) {
	params := mw.NewMwFetchParams().WithID(args[0])
	resp, err := mservapi.Mw.MwFetch(params, defaultAuth())
	if err != nil {
		log.WithError(err).Error("Couldn't fetch middleware")
		return
	}

	showFuncs := cmd.Flag("functions").Value.String() == "true"

	headers := []string{"ID", "Active", "Serve Only", "Last Update"}
	if showFuncs {
		headers = append(headers, "Function", "Type")
	}

	table := tablewriter.NewWriter(cmd.OutOrStdout())
	table.SetHeader(headers)
	table.SetBorder(false)
	table.SetRowLine(false)
	table.SetCenterSeparator("")
	table.SetAutoMergeCells(false)
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderAlignment(3)

	mw := resp.GetPayload().Payload

	row := []string{
		mw.UID,
		fmt.Sprintf("%v", mw.Active),
		fmt.Sprintf("%v", mw.DownloadOnly),
		mw.Added.String(),
	}
	table.Append(row)

	if showFuncs {
		for _, pl := range mw.Plugins {
			plrow := []string{
				"",
				"",
				"",
				"",
				pl.Name,
				coprocess.HookType(pl.Type).String(),
			}
			table.Append(plrow)
		}
	}

	table.Render()
}

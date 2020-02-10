package cmd

import (
	"fmt"

	coprocess "github.com/TykTechnologies/mserv/coprocess/bindings/go"
	"github.com/TykTechnologies/mserv/mservclient/client/mw"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List middleware in mserv",
	Long:  `Lists middleware known to mserv. By default everything is listed, use filters to limit output.`,
	Run:   listMiddleware,
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolP("functions", "f", false, "Show plugin functions")
	listCmd.Flags().BoolP("downloadonly", "d", false, "Show only download middleware")
	listCmd.Flags().BoolP("pluginsonly", "p", false, "Show only plugin middleware")
}

func listMiddleware(cmd *cobra.Command, args []string) {
	resp, err := mservapi.Mw.MwListAll(mw.NewMwListAllParams(), defaultAuth())
	if err != nil {
		log.WithError(err).Error("Couldn't fetch middleware")
		return
	}

	showFuncs := cmd.Flag("functions").Value.String() == "true"
	downloadOnly := cmd.Flag("downloadonly").Value.String() == "true"
	pluginsOnly := cmd.Flag("pluginsonly").Value.String() == "true"

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

	mws := resp.GetPayload().Payload
	for _, mw := range mws {
		if downloadOnly && !mw.DownloadOnly {
			continue
		}
		if pluginsOnly && mw.DownloadOnly {
			continue
		}

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
	}
	table.Render()
}

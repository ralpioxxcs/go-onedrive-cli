package main

import (
	"log"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/ralpioxxcs/go-onedrive-cli/graph"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "ls",
	Short: "List items such as the unix command, ls",
	Run: func(cmd *cobra.Command, args []string) {
		accessToken := credential.GetString("access_token")

		res, err := graph.List(accessToken)
		if err != nil {
			log.Fatalf("ls error: (err: %v)", err)
		}

		t := table.NewWriter()
		// t.SetColumnConfigs([]table.ColumnConfig{
		// 	{
		// 		Align:    text.AlignRight,
		// 		WidthMax: 64,
		// 	},
		// })
		t.SetTitle("Path: /drive/root")
		t.SetStyle(table.StyleLight)
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Index", "File Name", "Create Date"})

		for i, v := range res {
			t.AppendRows([]table.Row{{i, v.Name, v.CreateDateTime}})
			//t.AppendSeparator()
		}
		t.Render()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	var longListing bool
	listCmd.Flags().BoolVarP(&longListing, "detail", "d", false, "show more informtaion")
}

// func display(files []string) {
// 	const column = 5

// 	tw := new(tabwriter.Writer)
// 	tw.Init(os.Stdout, 8, 8, 0, '\t', 0) // minWidth, tabWidth, padding

// 	// get maximum width of strings
// 	maxWidth := 0
// 	for _, s := range files {
// 		if len(s) > maxWidth {
// 			maxWidth = len(s)
// 		}
// 	}
// 	format := fmt.Sprintf("%%-%ds%%s", maxWidth)

// 	rows := (len(files) + column - 1) / column
// 	for row := 0; row < rows; row++ {
// 		for col := 0; col < column; col++ {
// 			i := col*rows + row
// 			if i >= len(files) {
// 				break
// 			}
// 			padding := ""
// 			if i < 9 {
// 				padding = " "
// 			}
// 			fmt.Printf(format, files[i], padding)
// 		}
// 		fmt.Printf("\n")
// 	}
// 	// //numFiles := len(files)
// 	// for _, v := range files {
// 	// 	fmt.Fprintf(tw, "%v\t", v)
// 	// }
// 	// tw.Flush()
// }

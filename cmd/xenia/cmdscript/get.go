package cmdscript

import (
	"encoding/json"

	"github.com/coralproject/xenia/cmd/xenia/web"
	"github.com/coralproject/xenia/pkg/script"

	"github.com/spf13/cobra"
)

var getLong = `Retrieves a Script record from the system with the supplied name.

Example:
	script get -n pre_script
`

// get contains the state for this command.
var get struct {
	name string
}

// addGet handles the retrival Script records, displayed in json formatted response.
func addGet() {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Retrieves a Script record by name.",
		Long:  getLong,
		Run:   runGet,
	}

	cmd.Flags().StringVarP(&get.name, "name", "n", "", "Name of the Script.")

	scriptCmd.AddCommand(cmd)
}

// runGet is the code that implements the get command.
func runGet(cmd *cobra.Command, args []string) {
	if conn == nil {
		runGetWeb(cmd)
		return
	}

	runGetDB(cmd)
}

// runListWeb issues the command talking to the web service.
func runGetWeb(cmd *cobra.Command) {
	verb := "GET"
	url := "/1.0/script/" + get.name

	resp, err := web.Request(cmd, verb, url, nil)
	if err != nil {
		cmd.Println("Getting Script : ", err)
	}

	cmd.Printf("\n%s\n\n", resp)
}

// runGetDB issues the command talking to the DB.
func runGetDB(cmd *cobra.Command) {
	cmd.Printf("Getting Script : Name[%s]\n", get.name)

	if get.name == "" {
		cmd.Help()
		return
	}

	set, err := script.GetByName("", conn, get.name)
	if err != nil {
		cmd.Println("Getting Script : ", err)
		return
	}

	data, err := json.MarshalIndent(set, "", "    ")
	if err != nil {
		cmd.Println("Getting Script : ", err)
		return
	}

	cmd.Printf("\n%s\n\n", string(data))
	return
}

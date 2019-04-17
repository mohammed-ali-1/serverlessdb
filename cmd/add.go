// Copyright © 2019 Mohammed Al-Ameen <mohammed.alameen@protonmail.ch>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var (
	hostname string
	trigger  string
	table    string
)

func configureDb(cfg configuration) {
	var connStr = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.DbName)
	db, err := sql.Open("postgres", connStr)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	var notifyEventFunc = `CREATE OR REPLACE FUNCTION notify_event() RETURNS TRIGGER AS $$

    DECLARE 
        data json;
        notification json;
    
    BEGIN
    
        -- Convert the old or new row to JSON, based on the kind of action.
        -- Action = DELETE?             -> OLD row
        -- Action = INSERT or UPDATE?   -> NEW row
        IF (TG_OP = 'DELETE') THEN
            data = row_to_json(OLD);
        ELSE
            data = row_to_json(NEW);
        END IF;
        
        -- Contruct the notification as a JSON string.
        notification = json_build_object(
                          'table',TG_TABLE_NAME,
                          'action', TG_OP,
                          'data', data);
        
                        
        -- Execute pg_notify(channel, notification)
        PERFORM pg_notify('events',notification::text);
        
        -- Result is ignored since this is an AFTER trigger
        RETURN NULL; 
    END;
    
$$ LANGUAGE plpgsql;`

	var triggerCreate = `CREATE TRIGGER products_notify_event
AFTER INSERT OR UPDATE OR DELETE ON products
	FOR EACH ROW EXECUTE PROCEDURE notify_event();`

	//Creates the notify_event function
	rows, err := db.Exec(notifyEventFunc)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(rows)
	}

	//Create trigger and connect it to table
	rows, err = db.Exec(triggerCreate)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(rows)
	}

}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a trigger to a Postgresql table",
	Long:  `Adds a trigger to a Postgresql table`,
	Run: func(cmd *cobra.Command, args []string) {
		configureDb(cfg)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	addCmd.PersistentFlags().StringVar(&trigger, "trigger", "", "Trigger name (required)")
	addCmd.PersistentFlags().StringVar(&table, "table", "", "Table name (required)")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

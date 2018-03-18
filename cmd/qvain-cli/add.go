
package main

import (
	"fmt"
	"flag"
	"io/ioutil"

	"github.com/wvh/uuid/flag"
	"github.com/NatLibFi/qvain-api/psql"
	"github.com/NatLibFi/qvain-api/models"
)


func runAddRecord(psql *psql.PsqlService, args []string) error {
	flags := flag.NewFlagSet("add", flag.ExitOnError)
	var (
		creator uuidflag.Uuid
		owner   uuidflag.Uuid // = uuidflag.DefaultFromString("053bffbcc41edad4853bea91fc42ea18") // 053bffbcc41edad4853bea91fc42ea18
		schema  string
	)
	flags.Var(&creator,      "creator","creator `uuid`")
	flags.Var(&owner,        "owner",  "owner `uuid`")
	flags.StringVar(&schema, "schema", "metax", "schema identifier for given metadata record")
	
	flags.Usage = usageFor(flags, "add [flags] <json file>")
	if err := flags.Parse(args); err != nil {
		return err
	}
	
	if flags.NArg() < 1 {
		flags.Usage()
		return fmt.Errorf("error: missing some required arguments")
	}
	
	if !creator.IsSet() {
		return fmt.Errorf("error: flag `creator` must be set")
	}
	
	if schema == "" {
		return fmt.Errorf("error: flag `schema` must be set")
	}
	
	blob, err := ioutil.ReadFile(flags.Arg(0))
	if err != nil {
		return fmt.Errorf("error: can't read record: %s", err)
	}
	
	record, err := models.NewRecord(creator.Get())
	if err != nil {
		return err
	}
	record.SetMetadata(schema, string(blob))
	fmt.Printf("%+v\n", record)
	
	err = psql.Store(record)
	if err != nil {
		return err
	}
	
	return nil
}

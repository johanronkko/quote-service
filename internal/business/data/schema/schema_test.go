package schema

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	migs := parseMigrations(schemaDoc)
	var buf bytes.Buffer
	for _, mig := range migs {
		buf.WriteString(fmt.Sprintf("-- version: %.1f\n", mig.Version))
		buf.WriteString(fmt.Sprintf("-- description: %s\n", mig.Description))
		buf.WriteString(mig.Script)
	}

	sql := strings.ToLower(schemaDoc)
	if sql != buf.String() {
		// Don't forget to add an extra line at the end of the .sql file in case
		// got and exp looks identical.
		t.Logf("got: %v", string(buf.Bytes()))
		t.Logf("exp: %v", string([]byte(sql)))
	}
}

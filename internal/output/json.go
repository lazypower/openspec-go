package output

import (
	"encoding/json"
	"fmt"
	"io"
)

// WriteJSON marshals v as indented JSON and writes to w.
func WriteJSON(w io.Writer, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, string(data))
	return err
}

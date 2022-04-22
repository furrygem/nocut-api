package links

import (
	"encoding/json"
	"fmt"
)

// URLCheckException describes a result of a failed URL check
type URLCheckException struct {
	HostIsUp   bool  `json:"host_is_up"`
	URLIsValid bool  `json:"url_is_valid"`
	Err        error `json:"-"`
}

func (e *URLCheckException) Error() string {
	r, _ := json.Marshal(e)
	return fmt.Sprintf("URL checks failed. '%s'", r)
}

package routes_base

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Default Handler/
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Default route\n")
}

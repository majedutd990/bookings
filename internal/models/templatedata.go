package models

import "github.com/majedutd990/bookings/internal/forms"

//TemplateData contains our data to be imported in our template(sent from handlers to templats).
// Can have every data in map formats
// We use interface{} to map a keyed data which we have not created yet
// Cross Site Request Forgery Token CSRFToken it is related to security in forms
// We deal with it later
// Also we add some messages to users
type TemplateData struct {
	StrMap    map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
	Form      *forms.Form
}

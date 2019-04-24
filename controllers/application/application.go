package application

// import (
// 	"html/template"
// 	"net/http"
// )
//
// func Application(rw http.ResponseWriter, r *http.Request) {
// 	data := make(map[string]interface{})
//
// 	tname := "main"
// 	t, err := template.New(tname).ParseFiles("templates/main.tmpl", "templates/index.tmpl")
// 	if err != nil {
// 		http.Error(rw, "Error parsing templates.", 400)
// 	}
// 	err = t.ExecuteTemplate(rw, tname, data)
// 	if err != nil {
// 		http.Error(rw, "Error executing templates.", 400)
// 	}
// }

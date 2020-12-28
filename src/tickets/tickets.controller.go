package tickets

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type TicketsController struct {
	ticketsTemplate *template.Template
	DAL             DataHandler
}

func (t *TicketsController) RegisterRoutes() {
	http.HandleFunc("/api/tickets", t.handleTickets)
	http.HandleFunc("/tickets", t.showTickets)
}

func (t *TicketsController) RegisterTemplate(layout *template.Template) {
	f, err := os.Open("tickets/tickets.html")
	defer f.Close()
	if err != nil {
		log.Fatalf("Failed to open template: %v", err)
	}
	content, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("Failed to read content from '%v': %v", f.Name(), err)
	}
	tmpl := template.Must(layout.Clone())
	_, err = tmpl.Parse(string(content))
	if err != nil {
		log.Fatalf("Failed to parse contents of '%v' as template: %v", f.Name(), err)
	}
	t.ticketsTemplate = tmpl
}

func (t *TicketsController) showTickets(w http.ResponseWriter, r *http.Request){
	tickets := t.DAL.GetTickets()
	w.Header().Set("Content-Type", "text/html")
	t.ticketsTemplate.Execute(w, tickets)
}

func (t *TicketsController) handleTickets(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tickets := t.DAL.GetTickets()
		ticketsJSON, err := json.Marshal(tickets)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(ticketsJSON)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

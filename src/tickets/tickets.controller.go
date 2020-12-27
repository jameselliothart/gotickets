package tickets

import (
	"encoding/json"
	"html/template"
	"net/http"
)

type TicketsController struct {
	TicketsTemplate *template.Template
	DAL             DataHandler
}

func (t TicketsController) RegisterRoutes() {
	http.HandleFunc("/api/tickets", t.handleTickets)
}

func (t TicketsController) handleTickets(w http.ResponseWriter, r *http.Request) {
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

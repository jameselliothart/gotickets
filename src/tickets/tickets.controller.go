package tickets

import (
	"encoding/json"
	"html/template"
	"io"
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

func (t *TicketsController) RegisterTemplates(layout *template.Template) {
	const basePath = "tickets/templates/"
	t.ticketsTemplate = t.registerTemplate(layout, basePath+"tickets.html")
	t.newTicketTemplate = t.registerTemplate(layout, basePath+"new_ticket.html")
}

func (t *TicketsController) registerTemplate(layout *template.Template, fileName string) *template.Template {
	f, err := os.Open(fileName)
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
	return tmpl
}

func (t *TicketsController) showTickets(w http.ResponseWriter, r *http.Request) {
	tickets := t.DAL.GetTickets()
	w.Header().Set("Content-Type", "text/html")
	t.ticketsTemplate.Execute(w, tickets)
}

func (t *TicketsController) handleTickets(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tickets := t.DAL.GetTickets()
		w.Header().Set("Content-Type", "application/json")
		encodeAsJSON(tickets, w)
	case http.MethodPost:
		var dto CreateTicketDto
		err := json.NewDecoder(r.Body).Decode(&dto)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			log.Printf("Could not parse request: %v", err)
			return
		}
		id, err := t.DAL.CreateTicket(dto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			log.Printf("Could not create ticket: %v", err)
			return
		}
		w.WriteHeader(http.StatusCreated)
		encodeAsJSON(id, w)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func encodeAsJSON(data interface{}, w io.Writer) {
	enc := json.NewEncoder(w)
	enc.Encode(data)
}

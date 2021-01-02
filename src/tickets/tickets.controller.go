package tickets

import (
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/jameselliothart/gotickets/cqrs"
)

type TicketsController struct {
	ticketsTemplate   *template.Template
	newTicketTemplate *template.Template
	QueryHandler      TicketQueryHandler
	CommandHandler    cqrs.CommandHandler
}

func (t *TicketsController) RegisterRoutes() {
	http.HandleFunc("/tickets", t.showTickets)
	http.HandleFunc("/tickets/new", t.newTicket)
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

func (t *TicketsController) newTicket(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "text/html")
		t.newTicketTemplate.Execute(w, CreateTicketDto{})
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			log.Printf("Cannot parse form: %v", err)
			return
		}
		ticketToCreate := NewCreateTicketCmd(r.Form.Get("summary"))
		cqrs.LogWithCorrelation(ticketToCreate, "Create Ticket Command created:", ticketToCreate)
		t.CommandHandler.HandleCommand(ticketToCreate)
		http.Redirect(w, r, "/tickets", http.StatusFound)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (t *TicketsController) showTickets(w http.ResponseWriter, r *http.Request) {
	query := ActiveTicketsQuery{}
	tickets := t.QueryHandler.HandleQuery(query)
	w.Header().Set("Content-Type", "text/html")
	t.ticketsTemplate.Execute(w, tickets)
}

func encodeAsJSON(data interface{}, w io.Writer) {
	enc := json.NewEncoder(w)
	enc.Encode(data)
}

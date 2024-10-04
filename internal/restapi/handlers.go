package restapi

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sunwild/domain-checker_api/internal/grpcapi"
	service "github.com/sunwild/domain-checker_api/internal/service/domains"
	"github.com/sunwild/domain-checker_api/pkg/domains"
	"log"
	"net/http"
	"strconv"
)

type DomainHandler struct {
	Service *service.Service
}

func NewDomainHandler(service *service.Service) *DomainHandler {
	return &DomainHandler{Service: service}
}

// handler for get all domain
func (h *DomainHandler) GetAllDomains(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling request: GetAllDomains")
	domainList, err := h.Service.GetAllDomains()
	if err != nil {
		log.Printf("Error getting all domains: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	log.Println("Successfully fetched all domains")
	json.NewEncoder(w).Encode(domainList)
}

// handler for get domain on ID
func (h *DomainHandler) GetDomainById(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling request: GetDomainById")
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid domain ID: %s", idStr)
		http.Error(w, "Invalid domain ID", http.StatusBadRequest)
		return
	}
	domain, err := h.Service.GetDomainById(id)
	if err != nil {
		log.Printf("Error getting domain by ID %d: %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	log.Printf("Successfully fetched domain with ID: %d", id)
	json.NewEncoder(w).Encode(domain)
}

func (h *DomainHandler) AddDomain(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling request: AddDomain")
	var domain domains.Domain

	if err := json.NewDecoder(r.Body).Decode(&domain); err != nil {
		log.Printf("Error decoding domain payload: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := h.Service.AddDomain(&domain)
	if err != nil {
		log.Printf("Error adding domain: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	log.Printf("Successfully added domain: %s", domain.Name)
	json.NewEncoder(w).Encode(domain)

}

func (h *DomainHandler) DeleteDomainById(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling request: DeleteDomainById")
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid domain ID: %s", idStr)
		http.Error(w, "Invalid domain ID", http.StatusBadRequest)
	}
	err = h.Service.DeleteDomainById(id)
	if err != nil {
		log.Printf("Error deleting domain by ID %d: %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	log.Printf("Successfully deleted domain with ID: %d", id)
	json.NewEncoder(w).Encode(id)
}

func (h *DomainHandler) UpdateDomainById(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling request: UpdateDomainById")
	var domain domains.Domain
	if err := json.NewDecoder(r.Body).Decode(&domain); err != nil {
		log.Printf("Error decoding update payload: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
	}
	err := h.Service.UpdateDomainById(&domain)
	if err != nil {
		log.Printf("Error updating domain: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	log.Printf("Successfully updated domain with ID: %d", domain.ID)
	json.NewEncoder(w).Encode(domain)
}

// Обработка ручной проверки доменов через gRPC
func (h *DomainHandler) HandleManualDomainCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling request: HandleManualDomainCheck")
	var domainsList []string

	// Декодируем список доменов из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&domainsList); err != nil {
		log.Printf("Error decoding domain list: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Отправляем запрос через gRPC на проверку доменов
	log.Printf("Sending gRPC request to check domains: %v", domainsList)
	response, err := grpcapi.CheckDomainsWithGRPC(domainsList, true)
	if err != nil {
		log.Printf("gRPC request failed: %v", err)
		if response != nil && len(response.Statuses) > 0 {
			// Если есть частичный ответ, возвращаем его
			log.Println("Returning partial results despite gRPC error")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
			return
		} else {
			// Если ответа нет совсем, возвращаем ошибку
			http.Error(w, "gRPC request failed", http.StatusGatewayTimeout)
			return
		}
	}

	log.Println("Successfully checked domains, returning response")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

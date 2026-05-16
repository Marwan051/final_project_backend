package handlers

import (
	"log"
	"net/http"
	"strings"

	agent_service "github.com/Marwan051/final_project_backend/internal/service/agent"
	"github.com/Marwan051/final_project_backend/internal/utils"
)

type AgentHandler struct {
	client agent_service.Agent
}

func NewAgentHandler(client agent_service.Agent) *AgentHandler {
	return &AgentHandler{client: client}
}

// Query processes a message exchange with the agent.
// @Summary      Agent Query
// @Description  Sends a message to the agent and returns the reply plus session id
// @Tags         agent
// @Accept       json
// @Produce      json
// @Param        request body agent_service.Request true "Agent Query Request"
// @Success      200  {object}  agent_service.Response
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /api/v1/agent/query [post]
func (h *AgentHandler) Query(w http.ResponseWriter, r *http.Request) {
	var req agent_service.Request
	if err := utils.DecodeJSONBody(r, &req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if strings.TrimSpace(req.Message) == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "Missing required message")
		return
	}

	resp, err := h.client.Query(r.Context(), req)
	if err != nil {
		log.Printf("Error processing agent query: %v", err)
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to process agent request")
		return
	}

	if err := utils.WriteJSONResponse(w, http.StatusOK, resp); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
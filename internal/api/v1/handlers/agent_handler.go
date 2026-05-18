package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"

	agent_service "github.com/Marwan051/final_project_backend/internal/service/agent"
	"github.com/Marwan051/final_project_backend/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
// @Param        request body agent_service.AgentRequest true "Agent Query Request"
// @Success      200  {object}  agent_service.AgentResponse
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /api/v1/agent/query [post]
func (h *AgentHandler) Query(w http.ResponseWriter, r *http.Request) {
	var req agent_service.AgentRequest
	if err := utils.DecodeJSONBody(r, &req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if strings.TrimSpace(req.UserQuery) == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "Missing required message")
		return
	}

	resp, err := h.client.ProcessQuery(r.Context(), req)
	if err != nil {
		log.Printf("Error processing agent query: %v", err)

		var grpcStatus interface{ GRPCStatus() *status.Status }
		if errors.As(err, &grpcStatus) {
			switch grpcStatus.GRPCStatus().Code() {
			case codes.Unimplemented:
				utils.WriteJSONError(w, http.StatusBadGateway, "Agent service does not implement query")
				return
			case codes.Unavailable:
				utils.WriteJSONError(w, http.StatusServiceUnavailable, "Agent service is unavailable")
				return
			}
		}

		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to process agent request")
		return
	}

	if err := utils.WriteJSONResponse(w, http.StatusOK, resp); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

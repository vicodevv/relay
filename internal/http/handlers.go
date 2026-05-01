package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicodevv/relay/pkg/workflow"
)

type CreateDefinitionRequest struct {
	Name          string                  `json:"name" binding:"required"`
	Description   string                  `json:"description"`
	Steps         []workflow.Step         `json:"steps" binding:"required"`
	Compensations []workflow.Compensation `json:"compensations"`
}

type CreateWorkflowRequest struct {
	DefinitionName string                 `json:"definition_name" binding:"required"`
	Input          map[string]interface{} `json:"input"`
}

func (s *Server) createDefinition(c *gin.Context) {
	var req CreateDefinitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	def := &workflow.WorkflowDefinition{
		Name:          req.Name,
		Description:   req.Description,
		Version:       1,
		Steps:         req.Steps,
		Compensations: req.Compensations,
	}

	id, err := s.repo.CreateDefinition(c.Request.Context(), def)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      id,
		"message": "Workflow definition created successfully",
	})
}

func (s *Server) createWorkflow(c *gin.Context) {
	var req CreateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	def, err := s.repo.GetDefinitionByName(c.Request.Context(), req.DefinitionName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow definition not found"})
		return
	}

	instanceID, err := s.repo.CreateInstance(c.Request.Context(), def.ID, req.Input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = s.repo.CreateEvent(c.Request.Context(), instanceID, "workflow_created", nil, req.Input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      instanceID,
		"status":  workflow.StatusPending,
		"message": "Workflow instance created successfully",
	})
}

func (s *Server) getWorkflow(c *gin.Context) {
	id := c.Param("id")

	instance, err := s.repo.GetInstance(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
		return
	}

	c.JSON(http.StatusOK, instance)
}

func (s *Server) listWorkflows(c *gin.Context) {
	statusParam := c.Query("status")

	var status *workflow.WorkflowStatus
	if statusParam != "" {
		s := workflow.WorkflowStatus(statusParam)
		status = &s
	}

	instances, err := s.repo.ListInstances(c.Request.Context(), status, 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"workflows": instances,
		"count":     len(instances),
	})
}

func (s *Server) getWorkflowEvents(c *gin.Context) {
	id := c.Param("id")

	events, err := s.repo.GetInstanceEvents(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"instance_id": id,
		"events":      events,
		"count":       len(events),
	})
}

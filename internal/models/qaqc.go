package models

import "time"

type QAQCEntry struct {
	ID                int       `json:"id"`
	ProcessType       string    `json:"processType"`
	ProcessRef        string    `json:"processRef"`
	ContainersSampled int       `json:"containersSampled"`
	SampledQuantity   int       `json:"sampledQuantity"`
	SampledBy         string    `json:"sampledBy"`
	SampledOn         time.Time `json:"sampledOn"`
	ARNumber          string    `json:"arNumber"`
	ReleaseDate       time.Time `json:"releaseDate"`
	Potency           string    `json:"potency"`
	MoistureContent   string    `json:"moistureContent"`
	YieldPercent      string    `json:"yieldPercent"`
	Status            string    `json:"status"`
	AnalystRemark     string    `json:"analystRemark"`
	AnalysedBy        string    `json:"analysedBy"`
	ApprovedBy        string    `json:"approvedBy"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

// ==== API Structs ====

type CreateQAQCRequest struct {
	ProcessType       string `json:"processType" binding:"required"`
	ProcessRef        string `json:"processRef" binding:"required"`
	ContainersSampled int    `json:"containersSampled"`
	SampledQuantity   int    `json:"sampledQuantity"`
	SampledBy         string `json:"sampledBy"`
	SampledOn         string `json:"sampledOn"`
	ARNumber          string `json:"arNumber"`
	ReleaseDate       string `json:"releaseDate"`
	Potency           string `json:"potency"`
	MoistureContent   string `json:"moistureContent"`
	YieldPercent      string `json:"yieldPercent"`
	Status            string `json:"status"`
	AnalystRemark     string `json:"analystRemark"`
	AnalysedBy        string `json:"analysedBy"`
	ApprovedBy        string `json:"approvedBy"`
}

type CreateQAQCResponse struct {
	Message string `json:"message"`
}

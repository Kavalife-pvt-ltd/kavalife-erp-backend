package models

import "time"

type GRN struct {
	ID              int64  `json:"id"`
	GRNNumber       string `json:"grn_number"`
	VIRNumber       string `json:"vir_number"`
	ContainerQty    int    `json:"container_quantity"`
	Quantity        int    `json:"quantity"`
	Invoice         string `json:"invoice"`
	InvoiceDate     string `json:"invoice_date"`
	InvoiceImg      string `json:"invoice_img"`
	PackagingStatus string `json:"packaging_status"`
	CreatedBy       int64  `json:"created_by"`
	CreatedAt       string `json:"created_at"`
}

type CreateGRNRequest struct {
	VIRNumber       string `json:"virNumber" binding:"required"`
	ContainerQty    int    `json:"containerQuantity" binding:"required"`
	Quantity        int    `json:"quantity" binding:"required"`
	Invoice         string `json:"invoice"`
	InvoiceDate     string `json:"invoiceDate"`
	InvoiceImg      string `json:"invoiceImg"`
	PackagingStatus string `json:"packagingStatus"`
	CreatedBy       int64  `json:"createdBy" binding:"required"`
}

type CreateGRNResponse struct {
	GRNNumber string `json:"grn_number"`
	CreatedAt string `json:"created_at"`
}

type GRNResponse struct {
	ID              int64     `json:"id"`
	GRNNumber       string    `json:"grn_number"`
	CreatedAt       time.Time `json:"created_at"`
	ContainerQty    int       `json:"container_qty"`
	Quantity        float64   `json:"quantity"`
	Invoice         string    `json:"invoice"`
	InvoiceDate     time.Time `json:"invoice_date"`
	InvoiceImg      string    `json:"invoice_img"`
	PackagingStatus string    `json:"packaging_status"`
	CreatedBy       string    `json:"created_by"`
	VIRNumber       string    `json:"vir_number"`
	ProductName     string    `json:"product_name"`
	VendorName      string    `json:"vendor_name"`
}

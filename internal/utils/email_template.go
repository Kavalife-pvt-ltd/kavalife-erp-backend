package utils

import (
	"bytes"
	"fmt"
	"html/template"
)

type NewPONotificationData struct {
	PONumber          string
	CompanyName       string
	CompanyAddress    string
	RequestType       string
	Quantity          float64
	QuantityUnit      string
	Purity            string
	Grade             string
	ExpectedDelivery  string
	CreatedAt         string
	AdditionalComment string
}

var newPOTemplate = template.Must(template.New("new_po_email").Parse(`
<!DOCTYPE html>
<html>
  <body style="font-family: Arial, sans-serif; font-size: 14px; color: #222;">
    <h2 style="color:#0b8457; margin-bottom: 8px;">New Sales PO Created</h2>
    <p>A new sales PO has been created in Kavalife ERP.</p>

    <table cellpadding="6" cellspacing="0" style="border-collapse: collapse; margin-top: 8px;">
      <tr>
        <td style="font-weight:bold;">PO Number</td>
        <td>{{.PONumber}}</td>
      </tr>
      <tr>
        <td style="font-weight:bold;">Company</td>
        <td>{{.CompanyName}}</td>
      </tr>
      <tr>
        <td style="font-weight:bold;">Address</td>
        <td>{{.CompanyAddress}}</td>
      </tr>
      <tr>
        <td style="font-weight:bold;">Request Type</td>
        <td>{{.RequestType}}</td>
      </tr>
      <tr>
        <td style="font-weight:bold;">Quantity</td>
        <td>{{.Quantity}} {{.QuantityUnit}}</td>
      </tr>
      <tr>
        <td style="font-weight:bold;">Purity</td>
        <td>{{.Purity}}</td>
      </tr>
      <tr>
        <td style="font-weight:bold;">Grade</td>
        <td>{{.Grade}}</td>
      </tr>
      <tr>
        <td style="font-weight:bold;">Expected Delivery</td>
        <td>{{.ExpectedDelivery}}</td>
      </tr>
      <tr>
        <td style="font-weight:bold;">Created At</td>
        <td>{{.CreatedAt}}</td>
      </tr>
    </table>

    {{if .AdditionalComment}}
      <p style="margin-top:12px;">
        <strong>Sales Rep Comments:</strong><br/>
        {{.AdditionalComment}}
      </p>
    {{end}}

    <p style="margin-top:16px;">
      You can review and approve this PO in the Sales Dashboard.
    </p>

    <p style="font-size:12px; color:#777; margin-top:20px;">
      — Kavalife ERP
    </p>
  </body>
</html>
`))

func BuildNewPONotificationEmail(data NewPONotificationData) (string, string, error) {
	subject := fmt.Sprintf("[Kavalife ERP] New Sales PO %s - %s", data.PONumber, data.CompanyName)

	var buf bytes.Buffer
	if err := newPOTemplate.Execute(&buf, data); err != nil {
		return "", "", fmt.Errorf("failed to render new PO email template: %w", err)
	}

	return subject, buf.String(), nil
}

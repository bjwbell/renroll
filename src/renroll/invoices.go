package renroll
import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"github.com/jung-kurt/gofpdf"
)

func PrintInvoicesHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("printInvoicesHandler - begin")
	pdf := gofpdf.New("P", "in", "Letter", "")
	pdf.SetFont("Arial", "", 12)
	header := []string{"Base Rent", "Electricity", "Gas", "Water", "Sewage/Trash/Rec.", "Total"}
	widths := []float64{1.0, 1.0, 1.0, 1.0, 1.5, 1.0}
	dbName := r.FormValue("DbName")
	if dbName == "" {
		logError("printinvoices - no dbname set")
		return
	}
	tenants := dbReadTenants(dbName)
	formatCurrency(tenants)
	for _, tenant := range tenants {
		pdf.AddPage()
		h := 0.4
		pdf.CellFormat(5, h, "Name: " + tenant.Name + " (#" + strconv.Itoa(tenant.Id) + ")", "", 0, "", false, 0, "")
		pdf.Ln(-1)
		pdf.Ln(-1)
		for i, str := range header {
			pdf.CellFormat(widths[i], h, str, "1", 0, "", false, 0, "")
		}
		pdf.Ln(-1)
		pdf.CellFormat(widths[0], h, tenant.BaseRent, "1", 0, "", false, 0, "")
		pdf.CellFormat(widths[1], h, tenant.Electricity, "1", 0, "", false, 0, "")
		pdf.CellFormat(widths[2], h, tenant.Gas, "1", 0, "", false, 0, "")
		pdf.CellFormat(widths[3], h, tenant.Water, "1", 0, "", false, 0, "")
		pdf.CellFormat(widths[4], h, tenant.SewageTrashRecycle, "1", 0, "", false, 0, "")
		pdf.CellFormat(widths[5], h, tenant.Total, "1", 0, "", false, 0, "")
		pdf.Ln(-1)
	}
	err := pdf.OutputFileAndClose("invoices.pdf")
	if err != nil {
		logError("Unable to print invoices - write file")
		return
	}
	invoices, err := ioutil.ReadFile("invoices.pdf")
	if err != nil {
		logError("Unable to print invoices - read file")
		return
	}
	log.Print("printInvoicesTenantHandler - end")
	w.Write(invoices)
}

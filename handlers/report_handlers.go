package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/creator"
)

func init() {
	err := license.SetMeteredKey(os.Getenv("UNIDOC_LICENSE_KEY"))
	if err != nil {
		panic(err)
	}
}

func GenerateReport(c *gin.Context) {
	cr := creator.New()
	cr.NewPage()

	p := creator.NewParagraph("Hello, World!")
	p.SetPos(50, 700)
	cr.Draw(p)

	c.Writer.Header().Set("Content-Type", "application/pdf")
	err := cr.Write(c.Writer)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to generate PDF: %v", err)
		return
	}
}
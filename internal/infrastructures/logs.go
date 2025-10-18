package infrastructures

import (
	"log"
	"time"

	"github.com/julia-marcal/llm-with-go/internal/dto"
)

func LogAudit(audit *dto.LLMAudit) {
	log.Printf("📌 LLM Audit - %s", audit.Timestamp.Format(time.RFC3339))
	log.Printf("Prompt: %s", audit.Prompt)
	log.Printf("Temperature: %.2f", audit.Temperature)
	log.Printf("Duration: %s", audit.Duration)
	if audit.Err != nil {
		log.Printf("Error: %v", audit.Err)
	} else {
		log.Printf("Response: %s", audit.Response)
	}
	log.Println("-------------------------------------------------")
}

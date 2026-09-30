package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/PainQlers/backend/handlers"      // Import handlers เข้ามา
	"github.com/PainQlers/backend/pkg/telemetry" // 👈 ✨ นำเข้า Shared SDK ที่เราเพิ่งสร้าง
	"github.com/PainQlers/backend/pkg/telemetry/operation"
	"github.com/gofiber/contrib/otelfiber" // 👈 นำเข้าตัวเชื่อม Fiber เข้ากับ OTel
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"               // เพิ่มเข้ามาเพื่อโหลดไฟล์ .env
	"github.com/liushuangls/go-anthropic/v2" // เพิ่มเข้ามาเพื่อใช้ Claude SDK

	// instrument "github.com/PainQlers/backend/pkg/telemetry/instrument"
	appmetrics "github.com/PainQlers/backend/pkg/telemetry/metrics"
	middleware "github.com/PainQlers/backend/pkg/telemetry/middleware"
	// "github.com/prometheus/client_golang/prometheus"
)

// โครงสร้างสำหรับรับคำถามจาก Frontend ผ่าน body JSON
type ChatRequest struct {
	Message string `json:"message"`
}

// โครงสร้างสำหรับรับค่าจาก Claude ตอนเรียกใช้ฟังก์ชัน
type StockArgs struct {
	Ticker string `json:"ticker"`
}

// ฟังก์ชันจำลองการทำงานของ Backend เมื่อ AI สั่งเรียกใช้ Skill
func getStockPrice(ticker string) string {
	if strings.ToUpper(ticker) == "AAPL" {
		return "$180.50"
	} else if strings.ToUpper(ticker) == "GOOG" {
		return "$175.20"
	}
	return "ไม่พบข้อมูลหุ้นตัวนี้"
}

func main() {
	ctx := context.Background()

	// 1. โหลด Environment Variables จากไฟล์ .env ก่อนเริ่มระบบ
	if err := godotenv.Load(); err == nil {
		slog.Info("Loaded .env file")
	}

	// ดึง IP พิกัดเซิร์ฟเวอร์กลางมาจาก .env
	telemetryEndpoint := os.Getenv("TELEMETRY_ENDPOINT")
	if telemetryEndpoint == "" {
		telemetryEndpoint = "localhost:4317" // fallback
	}

	// เรียกเปิดใช้งานจาก Shared SDK ของเรา ชี้เป้าไปที่ Alloy พอร์ต 4317
	// ตั้งชื่อระบบนี้บนหน้า Grafana ว่า "ai-agent-service"
	shutdownTelemetry, err := telemetry.InitSharedTelemetry(ctx, telemetry.Config{
		AlloyEndpoint: telemetryEndpoint,
		ServiceName:   "ai-agent-service",
		MetricsPort:   "2112",
		Insecure:      true,
		EnableTrace:   true,
		EnableMetrics: true,
		EnableLogs:    true,
	}, "ai-agent-service")

	if err != nil {
		slog.Error("ล้มเหลวในการตั้งค่า Shared Telemetry", slog.Any("error", err))
		os.Exit(1)
	}
	defer shutdownTelemetry() // คืนคำสั่งปิดระบบเมื่อแอปสิ้นสุดการทำงาน

	appmetrics.Register()

	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	client := anthropic.NewClient(apiKey)

	// 2. สร้าง Instance ของ Fiber App
	app := fiber.New()

	// สวมหมวกดักจับข้อมูล (Middleware) โค้ดชุดนี้จะแอบเก็บสถิติ HTTP Metrics (Prometheus) ส่งอัตโนมัติ
	app.Use(middleware.RecoveryMiddleware())
	app.Use(otelfiber.Middleware())
	app.Use(middleware.MetricsMiddleware())
	app.Use(middleware.LoggingMiddleware())

	// 3. จัดกลุ่ม Route
	v1 := app.Group("/v1")

	// Route เดิมของคุณ
	v1.Get("/users/:id", handlers.GetUserByID)

	// ✨ 4. เพิ่ม Endpoint ใหม่สำหรับให้ระบบ AI Agent ทำงาน
	v1.Post("/agent/chat", func(c *fiber.Ctx) error {
		// fmt.Println("POST /agent/chat")
		// start := time.Now()

		// defer instrument.RecordAIRequestDuration(start)

		// สร้าง Custom Span (Tracing) เจาะลึก เพื่อไปส่องดูจังหวะที่ยิงไปหา Claude API ข้ามโลก
		op, agentCtx := operation.Start(
			c.UserContext(),
			"Claude AI Agent Request",
		)
		defer op.End() // สั่งจบการส่องบันทึกเวลาเมื่อพ้นฟังก์ชันนี้

		// รับข้อความคำถามจากผู้ใช้
		reqBody := new(ChatRequest)
		if err := c.BodyParser(reqBody); err != nil {
			op.Failed(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		// อ่านข้อมูลจากไฟล์ skills.md ออกมาใช้งาน
		skillsContent, err := os.ReadFile("skills.md")
		if err != nil {
			op.Failed(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "ไม่สามารถอ่านคลัง Skill ได้"})
		}

		// ดึงเฉพาะเนื้อหาที่เป็นโครงสร้าง JSON ภายในเครื่องหมาย [ ] ออกมาจาก Markdown
		fileStr := string(skillsContent)
		startIdx := strings.Index(fileStr, "[")
		endIdx := strings.LastIndex(fileStr, "]")
		if startIdx == -1 || endIdx == -1 {
			op.Failed(fmt.Errorf("invalid JSON schema in skills.md"))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "โครงสร้างข้อมูลใน skills.md ไม่ถูกต้อง"})
		}
		jsonSchema := fileStr[startIdx : endIdx+1]

		var tools []anthropic.ToolDefinition
		if err := json.Unmarshal([]byte(jsonSchema), &tools); err != nil {
			op.Failed(fmt.Errorf("failed to unmarshal JSON schema: %w", err))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "การแปลงข้อมูล JSON ล้มเหลว"})
		}

		// ส่งข้อมูลคุยกับ Claude API
		resp, err := client.CreateMessages(agentCtx, anthropic.MessagesRequest{
			Model:     "claude-3-5-sonnet-20240620",
			MaxTokens: 1000,
			Messages: []anthropic.Message{
				anthropic.NewUserTextMessage(reqBody.Message),
			},
			Tools: tools, // แนบคลังความสามารถที่ได้มาจาก skills.md ให้ Claude รู้จัก
		})

		if err != nil {
			op.Failed(err)

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "การเชื่อมต่อไปยัง Claude ล้มเหลว",
				"details": err.Error(), // ✨ ส่งรายละเอียดกลับไปบอกที่ PowerShell ด้วย
			})

		}

		slog.InfoContext(
			agentCtx,
			"Claude response received",
		)
		op.AddEvent("ได้รับคำตอบจากเซิร์ฟเวอร์ Anthropic สำเร็จแล้ว")

		// วนลูปเช็คว่า Claude ต้องการใช้งาน Skill หรือไม่
		aiResponse := ""
		for _, content := range resp.Content {
			if content.Type == anthropic.MessagesContentTypeText {
				aiResponse = *content.Text
			}

			if content.Type == anthropic.MessagesContentTypeToolUse {
				if content.Name == "get_stock_price" {
					var args StockArgs
					if err := json.Unmarshal(content.Input, &args); err != nil {
						op.Failed(fmt.Errorf("invalid tool arguments: %w", err))
						return c.Status(fiber.StatusInternalServerError).JSON(
							fiber.Map{
								"error": "invalid tool arguments",
							},
						)
					}

					// เรียกฟังก์ชันการทำงานฝั่ง Backend ของเราจริง ๆ
					result := getStockPrice(args.Ticker)

					// ส่งคำตอบสุดท้ายที่ได้จากการรันฟังก์ชันกลับไปให้ผู้ใช้
					aiResponse = fmt.Sprintf("ราคาหุ้น %s ปัจจุบันคือ %s ครับ", args.Ticker, result)
				}
			}
		}

		op.Success()

		// ส่งผลลัพธ์ตอบกลับในรูปแบบ JSON ไปหาผู้ใช้งาน API
		return c.JSON(fiber.Map{
			"reply": aiResponse,
		})
	})

	// 5. สั่งให้ Server ทำงานที่ Port 8080
	slog.Info("Server is running on port 8080...")
	if err := app.Listen(":8080"); err != nil {
		slog.Error("Server failed to start", slog.Any("error", err))
		os.Exit(1)
	}
}

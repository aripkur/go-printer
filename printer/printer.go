package printer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/alexbrainman/printer"
	"github.com/gofiber/fiber/v2"
)

type Printer struct {
}

func (p Printer) List() ([]string, error) {
	printers, err := printer.ReadNames()
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan daftar printer: %v", err)
	}
	printerDefault, err := printer.Default()
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan default printer: %v", err)
	}
	prns := []string{}
	for _, prn := range printers {
		if prn == printerDefault {
			prns = append(prns, prn+" (DEFAULT)")
		} else {
			prns = append(prns, prn)
		}
	}

	return prns, nil
}

func (p Printer) Default() (string, error) {
	printers, err := printer.Default()
	if err != nil {
		return "", fmt.Errorf("gagal mendapatkan default printer: %v", err)
	}

	return printers, nil
}
func (p Printer) PrintPdf(filePath string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("gagal mendapatkan path saat ini: %v", err)
	}
	toolPath := filepath.Join(currentDir, "PDFtoPrinter.exe")

	cmd := exec.Command(toolPath, filePath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("gagal mencetak : %v, output: %s", err, string(output))
	}

	return nil
}

func (p Printer) RunServer(port string) {
	app := fiber.New()

	app.Post("/print", func(c *fiber.Ctx) error {
		file, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Gagal menerima file",
				"error":   err.Error(),
			})
		}

		filePath := fmt.Sprintf("./%s", file.Filename)
		if err := c.SaveFile(file, filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Gagal menyimpan file",
				"error":   err.Error(),
			})
		}

		if err := p.PrintPdf(filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Gagal mencetak file",
				"error":   err.Error(),
			})
		}
		os.Remove(filePath)

		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "File berhasil dicetak",
		})
	})

	app.Get("/", func(c *fiber.Ctx) error {
		printers, err := p.List()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Gagal mendapatkan daftar printer",
				"error":   err.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"status":   "success",
			"printers": printers,
		})
	})
	app.Listen(fmt.Sprintf(":%s", port))
}

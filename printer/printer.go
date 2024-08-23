package printer

import (
	"fmt"
	"io"
	"os"

	"github.com/alexbrainman/printer"
	"github.com/gofiber/fiber/v2"
)

type Printer struct{}

func (p Printer) Print(filePath string, fileName string) error {
	pname, err := printer.Default()
	if err != nil {
		return fmt.Errorf("gagal menemukan printer default: %v", err)
	}

	prn, err := printer.Open(pname)
	if err != nil {
		return fmt.Errorf("gagal membuka printer: %v", err)
	}
	defer prn.Close()

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("gagal membuka file: %v", err)
	}
	defer f.Close()

	// Mengirimkan nama dokumen dan mode RAW
	err = prn.StartDocument(fileName, "RAW")
	if err != nil {
		return fmt.Errorf("gagal memulai dokumen: %v", err)
	}
	defer prn.EndDocument()

	err = prn.StartPage()
	if err != nil {
		return fmt.Errorf("gagal memulai halaman: %v", err)
	}
	defer prn.EndPage()

	_, err = io.Copy(prn, f)
	if err != nil {
		return fmt.Errorf("gagal menulis ke printer: %v", err)
	}

	return nil
}
func (p Printer) GetPrinters() ([]string, error) {
	printers, err := printer.ReadNames()
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan daftar printer: %v", err)
	}
	return printers, nil
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

		if err := p.Print(filePath, file.Filename); err != nil {
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
		printers, err := p.GetPrinters()
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

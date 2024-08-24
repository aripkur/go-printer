package printer

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

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

	f, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("gagal membuka file: %v", err)
	}

	// Mengirimkan nama dokumen dan mode RAW
	err = prn.StartDocument(fileName, "RAW")
	if err != nil {
		return fmt.Errorf("gagal memulai dokumen: %v", err)
	}
	defer prn.EndDocument()

	// err = prn.StartPage()
	// if err != nil {
	// 	return fmt.Errorf("gagal memulai halaman: %v", err)
	// }
	// defer prn.EndPage()

	// Menulis data PDF ke printer
	_, err = prn.Write(f)
	if err != nil {
		log.Fatalf("gagal menulis ke printer: %v", err)
	}

	return nil
}
func (p Printer) PrintTest(filePath string, fileName string) error {
	pname, err := printer.Default()
	if err != nil {
		return fmt.Errorf("gagal menemukan printer default: %v", err)
	}

	prn, err := printer.Open(pname)
	if err != nil {
		return fmt.Errorf("gagal membuka printer: %v", err)
	}
	defer prn.Close()

	err = prn.StartDocument("my document", "RAW")
	if err != nil {
		return fmt.Errorf("StartDocument failed: %v", err)
	}
	defer prn.EndDocument()
	err = prn.StartPage()
	if err != nil {
		return fmt.Errorf("StartPage failed: %v", err)
	}
	fmt.Fprintf(prn, "Hello")
	err = prn.EndPage()
	if err != nil {
		return fmt.Errorf("EndPage failed: %v", err)
	}

	return nil
}
func (p Printer) PrintSpooler(filePath string, fileName string) error {
	// Menemukan path ke SPool.exe di direktori saat ini
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("gagal mendapatkan path saat ini: %v", err)
	}
	spoolPath := filepath.Join(currentDir, "spool.exe")

	// Membangun perintah untuk mencetak
	cmd := exec.Command(spoolPath, filePath)

	// Menjalankan perintah
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("gagal mencetak dengan spool.exe: %v, output: %s", err, string(output))
	}

	fmt.Println("Print command executed successfully:", string(output))
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

		if err := p.PrintSpooler(filePath, file.Filename); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Gagal mencetak file",
				"error":   err.Error(),
			})
		}
		fmt.Println(file.Filename)
		// os.Remove(filePath)

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

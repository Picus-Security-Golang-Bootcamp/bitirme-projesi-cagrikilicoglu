package product

import (
	"encoding/csv"
	"mime/multipart"
	"strconv"
	"sync"

	"github.com/cagrikilicoglu/shopping-basket/internal/models"
	"go.uber.org/zap"
)

// readProductsWithWorkerPool: Reading a csv file concurrently and returns a products slice
func readProductsWithWorkerPool(fileHeader *multipart.FileHeader) ([]models.Product, error) {
	zap.L().Debug("product.csvService.readProductsWithWorkerPool")
	const numJobs = 5
	products := []models.Product{}
	jobs := make(chan []string, numJobs)
	results := make(chan models.Product, numJobs)
	wg := sync.WaitGroup{}

	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go toStruct(jobs, results, &wg)
	}
	go func() {
		f, err := fileHeader.Open()
		if err != nil {
			zap.L().Error("product.csvService.readProductsWithWorkerPool cannot open file", zap.Error(err))
			return
		}
		defer f.Close()

		lines, err := csv.NewReader(f).ReadAll()
		if err != nil {
			zap.L().Error("category.csvService.readProductsWithWorkerPool cannot read csv file", zap.Error(err))
			return
		}
		for _, line := range lines[1:] {
			jobs <- line
		}

		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for p := range results {
		products = append(products, p)
	}

	return products, nil
}

// toStruct: creates a product struct as the data from the file is read and send the struct to results channel
func toStruct(jobs <-chan []string, results chan<- models.Product, wg *sync.WaitGroup) {
	defer wg.Done()

	for j := range jobs {
		priceParsed, err := strconv.ParseFloat(j[2], 32)
		if err != nil {
			return
		}
		stockNumberParsed, err := strconv.Atoi(j[4])
		if err != nil {
			return
		}

		product := models.Product{
			CategoryName: &j[0],
			Name:         &j[1],
			Price:        float32(priceParsed),
			Stock: models.Stock{SKU: j[3],
				Number: uint(stockNumberParsed),
			}}
		results <- product
	}

}

package category

import (
	"encoding/csv"
	"mime/multipart"
	"sync"

	"github.com/cagrikilicoglu/shopping-basket/internal/models"
	"go.uber.org/zap"
)

// readCategoriesWithWorkerPool: Reading a csv file concurrently and returns a categories slice
func readCategoriesWithWorkerPool(fileHeader *multipart.FileHeader) ([]models.Category, error) {
	zap.L().Debug("category.csvService.readCategoriesWithWorkerPool")
	const numJobs = 5
	categories := []models.Category{}
	jobs := make(chan []string, numJobs)
	results := make(chan models.Category, numJobs)
	wg := sync.WaitGroup{}

	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go toStruct(jobs, results, &wg)
	}
	go func() {
		f, err := fileHeader.Open()
		if err != nil {
			zap.L().Error("category.csvService.readCategoriesWithWorkerPool cannot open file", zap.Error(err))
			return
		}
		defer f.Close()

		lines, err := csv.NewReader(f).ReadAll()
		if err != nil {
			zap.L().Error("category.csvService.readCategoriesWithWorkerPool cannot read csv file", zap.Error(err))
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

	for c := range results {
		categories = append(categories, c)
	}

	return categories, nil
}

// toStruct: creates a category struct as the data from the file is read and send the struct to results channel
func toStruct(jobs <-chan []string, results chan<- models.Category, wg *sync.WaitGroup) {
	defer wg.Done()

	for j := range jobs {
		category := models.Category{
			Name:        &j[0],
			Description: j[1]}
		results <- category
	}
}

package supabase

import (
	"fmt"

	"github.com/supabase-community/supabase-go"
)

const(
	PageSize        = 4000
    ProductTable    = "Products"
    RelationalTable = "relations_product_go_test"
    ModelTable      = "Model"
    MakeTable       = "Make"
    TypeTable       = "Type"
    YearsTable      = "Years"
)

type Product struct {
    ID              int    `json:"id"`
    SKU             string `json:"sku"`
    Make            string `json:"make"`
    Model           string `json:"model"`
    YearGeneration  string `json:"year_generation"`
    Type            string `json:"type"`
}

type RelationRow struct {
    ProductID int `json:"product_id"`
    TypeID    int `json:"type_id"`
    YearID    int `json:"year_id"`
    MakeID    int `json:"make_id"`
    ModelID   int `json:"model_id"`
}

func CreateRealationalTable(client, adminClient *supabase.Client) error {
	models, err := fetchAll(adminClient, ModelTable)
	if err != nil {
		return fmt.Errorf("error fetching models: %w",)
	}

	makes, err := fetchAll(adminClient, MakeTable)
	if err != nil {
		return fmt.Errorf("error fetching makes: %w", err)
	}

	types, err := fetchAll(adminClient, TypeTable)
    if err != nil {
        return fmt.Errorf("error fetching types: %w", err)
    }

    years, err := fetchAll(adminClient, YearsTable)
    if err != nil {
        return fmt.Errorf("error fetching years: %w", err)
    }

	pageIndex := 0
	hasMore := true

	for hasMore {
		products, err := fetchProducts(adminClient, pageIndex)
		if err != nil {
            return fmt.Errorf("error fetching products: %w", err)
        }

        if len(products) == 0 {
            hasMore = false
            break
        }

		if err := prepareAndBatchInsertData(adminClient, products, makes, models, years, types); err != nil {
            return fmt.Errorf("error preparing and inserting data: %w", err)
        }

        fmt.Printf("Data inserted for batch %d\n", pageIndex+1)

        if len(products) < PageSize {
            hasMore = false
        } else {
            pageIndex++
        }
	}

	return nil

}

func fetchAll(client *supabase.Client, table string) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	err := client.DB.From(table).Select("id", "name").Execute(&result)
	return result, err
}

func fetchProducts(client *supabase.Client, pageIndex int) ([]Product, error) {
	var products [Product]
	err := client.DB.From(ProductTable).
		Select("id", "sku", "make", "model", "year_generation", "type").
		Order("id", true).
		Range(pageIndex*PageSize, (pageIndex+1)*PageSize-1).
		Execute(&products)
	return products, err
}

func prepareAndBatchInsertData(client *supabase.Client, products []Product, makes, models, years, types []map[string]interface{}) error {
    var inserts []RelationRow
    batchSize := 1000

    for _, product := range products {
        row := RelationRow{
            ProductID: product.ID,
            TypeID:    1, // Assuming only car cover for now
        }

        // Find make_id, model_id, type_id (implementation similar to JavaScript version)
        // ...

        // Process year_generation
        // ...

        if len(inserts) >= batchSize {
            if err := batchInsert(client, inserts); err != nil {
                return err
            }
            inserts = []RelationRow{}
        }
    }

    if len(inserts) > 0 {
        return batchInsert(client, inserts)
    }

    return nil
}

func batchInsert(client *supabase.Client, rows []RelationRow) error {
    var result []map[string]interface{}
    err := client.DB.From(RelationalTable).Insert(rows).Execute(&result)
    if err != nil {
        return fmt.Errorf("error inserting data: %w", err)
    }
    return nil
}
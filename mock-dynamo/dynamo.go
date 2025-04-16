package mock_dynamo

import (
	"encoding/json"
	"fmt"
	"os"
)

const tableFilePath = "./mock-dynamo/table.json"

func GetItem(id string) (map[string]string, error) {
	jsonData, err := loadData()
	if err != nil {
		return nil, err
	}

	item, exists := jsonData[id]
	if !exists {
		return nil, fmt.Errorf("GUID %s not found", id)
	}

	return item, nil
}

func DeleteItem(sourceEntityId string) error {
	jsonData, err := loadData()
	if err != nil {
		return err
	}

	delete(jsonData, sourceEntityId)

	return writeData(jsonData, tableFilePath)
}

func SetItem(sourceGuid, targetOrgId, targetEntityId string) error {
	jsonData, err := loadData()
	if err != nil {
		return err
	}

	if _, exists := jsonData[sourceGuid]; exists {
		jsonData[sourceGuid][targetOrgId] = targetEntityId
	} else {
		jsonData[sourceGuid] = map[string]string{
			targetOrgId: targetEntityId,
		}
	}

	return writeData(jsonData, tableFilePath)
}

func loadData() (map[string]map[string]string, error) {
	data, err := os.ReadFile(tableFilePath)
	if err != nil {
		return nil, err
	}

	var jsonData = make(map[string]map[string]string)
	if err = json.Unmarshal(data, &jsonData); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	return jsonData, nil
}

func writeData(data any, filepath string) error {
	// Convert map to JSON with indentation for readability
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	// Write to file
	err = os.WriteFile(filepath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}

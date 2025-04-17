package mock_dynamo

import (
	"encoding/json"
	"fmt"
	"os"
)

const tableFilePath = "./mock-dynamo/table.json"

func GetTargetIdBySourceId(sourceId, orgId string) (string, error) {
	item, err := GetItem(sourceId)
	if err != nil {
		return "", nil
	}

	for _, target := range item.TargetInfo {
		if target.OrgId == orgId {
			return target.TargetEntityId, nil
		}
	}

	return "", fmt.Errorf("no target ID found for org '%s'. Source ID: '%s'", orgId, sourceId)
}

func GetItem(sourceEntityId string) (*Item, error) {
	table, err := loadData()
	if err != nil {
		return nil, err
	}

	if table == nil {
		return nil, fmt.Errorf("failed to load table")
	}

	for _, item := range table.Items {
		if sourceEntityId != item.SourceEntityId {
			continue
		}
		return &item, nil
	}

	return nil, fmt.Errorf("failed to find item for source ID '%s'", sourceEntityId)
}

func DeleteItem(sourceEntityId string) error {
	table, err := loadData()
	if err != nil {
		return err
	}

	newItemsSlice := make([]Item, 0)
	for _, item := range table.Items {
		if item.SourceEntityId != sourceEntityId {
			newItemsSlice = append(newItemsSlice, item)
		}
	}

	table.Items = newItemsSlice
	return writeData(*table, tableFilePath)
}

// UpdateItem updates or creates an item mapping between a source entity and its target information.
//
// Parameters:
//   - sourceEntityId: The unique identifier of the source entity
//   - targetOrgId: The organization ID where the target entity exists
//   - targetEntityId: The identifier of the target entity within the target organization
//
// Returns:
//   - error: Returns nil on successful update/creation, or an error if the operation fails
//
// The function performs the following operations:
//  1. Loads existing data from storage
//  2. Searches for an existing item with the given sourceEntityId
//  3. If found:
//     - Updates the targetEntityId if a matching targetOrgId exists
//     - Adds new target information if the targetOrgId is not found
//  4. If not found:
//     - Creates a new item with the provided source and target information
//  5. Persists the updated data back to storage
func UpdateItem(sourceEntityId, targetOrgId, targetEntityId string) error {
	table, err := loadData()
	if err != nil {
		return err
	}

	// Find the item with matching sourceEntityId
	var targetItem *Item
	for i := range table.Items {
		if table.Items[i].SourceEntityId == sourceEntityId {
			targetItem = &table.Items[i]
			break
		}
	}

	// If we found an existing item
	if targetItem != nil {
		// Try to update existing target info
		for i := range targetItem.TargetInfo {
			if targetItem.TargetInfo[i].OrgId == targetOrgId {
				targetItem.TargetInfo[i].TargetEntityId = targetEntityId
				return writeData(*table, tableFilePath)
			}
		}

		// If we get here, no matching target org was found, append new target info
		targetItem.TargetInfo = append(targetItem.TargetInfo, TargetInfo{
			OrgId:          targetOrgId,
			TargetEntityId: targetEntityId,
		})
	} else {
		// Create new item
		table.Items = append(table.Items, Item{
			SourceEntityId: sourceEntityId,
			TargetInfo: []TargetInfo{
				{
					OrgId:          targetOrgId,
					TargetEntityId: targetEntityId,
				},
			},
		})
	}

	return writeData(*table, tableFilePath)
}

func loadData() (*Table, error) {
	data, err := os.ReadFile(tableFilePath)
	if err != nil {
		return nil, err
	}

	var table Table
	if err = json.Unmarshal(data, &table); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	return &table, nil
}

func writeData(data Table, filepath string) error {
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

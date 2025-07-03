package models

import (
	"testing"
)

func TestCategory_Validate(t *testing.T) {
	tests := []struct {
		name     string
		category Category
		wantErr  bool
	}{
		{
			name: "valid category",
			category: Category{
				Name: "Mountain Bike",
			},
			wantErr: false,
		},
		{
			name: "invalid name",
			category: Category{
				Name: "",
			},
			wantErr: true,
		},
		{
			name: "invalid description",
			category: Category{
				Name: "Mountain Bike",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.category.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Category.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateCategoryRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateCategoryRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CreateCategoryRequest{
				Name: "Mountain Bike",
			},
			wantErr: false,
		},
		{
			name:    "missing name",
			req:     CreateCategoryRequest{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCategoryRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

package models

import (
	"testing"
)

func TestItem_Validate(t *testing.T) {
	tests := []struct {
		name    string
		item    Item
		wantErr bool
	}{
		{
			name: "valid item",
			item: Item{
				Title:       "Mountain Bike",
				Description: "High-quality mountain bike perfect for trail riding",
				ImageURL:    "https://example.com/bike.jpg",
				PricePerDay: 25.50,
				Location:    "San Francisco, CA",
			},
			wantErr: false,
		},
		{
			name: "invalid image URL",
			item: Item{
				Title:       "Mountain Bike",
				Description: "High-quality mountain bike perfect for trail riding",
				ImageURL:    "not-a-url",
				PricePerDay: 25.50,
				Location:    "San Francisco, CA",
			},
			wantErr: true,
		},
		{
			name: "title too short",
			item: Item{
				Title:       "",
				Description: "High-quality mountain bike perfect for trail riding",
				ImageURL:    "https://example.com/bike.jpg",
				PricePerDay: 25.50,
				Location:    "San Francisco, CA",
			},
			wantErr: true,
		},
		{
			name: "description too short",
			item: Item{
				Title:       "Mountain Bike",
				Description: "Too short",
				ImageURL:    "https://example.com/bike.jpg",
				PricePerDay: 25.50,
				Location:    "San Francisco, CA",
			},
			wantErr: true,
		},
		{
			name: "negative price",
			item: Item{
				Title:       "Mountain Bike",
				Description: "High-quality mountain bike perfect for trail riding",
				ImageURL:    "https://example.com/bike.jpg",
				PricePerDay: -10.0,
				Location:    "San Francisco, CA",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.item.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Item.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateItemRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateItemRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CreateItemRequest{
				Title:       "Mountain Bike",
				Description: "High-quality mountain bike perfect for trail riding",
				ImageURL:    "https://example.com/bike.jpg",
				PricePerDay: 25.50,
				Location:    "San Francisco, CA",
			},
			wantErr: false,
		},
		{
			name: "missing title",
			req: CreateItemRequest{
				Description: "High-quality mountain bike perfect for trail riding",
				ImageURL:    "https://example.com/bike.jpg",
				PricePerDay: 25.50,
				Location:    "San Francisco, CA",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateItemRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

package models

import (
	"testing"
)

func TestItemPhoto_Validate(t *testing.T) {
	tests := []struct {
		name      string
		itemPhoto ItemPhoto
		wantErr   bool
	}{
		{
			name: "valid item photo",
			itemPhoto: ItemPhoto{
				ID:     1,
				ItemID: 1,
				URL:    "https://example.com/bike.jpg",
			},
			wantErr: false,
		},
		{
			name: "invalid item photo",
			itemPhoto: ItemPhoto{
				ID:     1,
				ItemID: 0,
				URL:    "not-a-url",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.itemPhoto.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ItemPhoto.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateItemPhotoRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateItemPhotoRequest
		wantErr bool
	}{
		{
			name: "valid item photo",
			req: CreateItemPhotoRequest{
				ItemID: 1,
				Photos: []string{"https://example.com/bike.jpg"},
			},
			wantErr: false,
		},
		{
			name: "invalid item photo",
			req: CreateItemPhotoRequest{
				ItemID: 1,
				Photos: []string{"not-a-url"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateItemPhotoRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

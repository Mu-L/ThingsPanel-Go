package service

import (
	"testing"

	"project/internal/model"
)

func TestResolveMarketTemplateImageURL(t *testing.T) {
	tests := []struct {
		name     string
		data     *model.MarketTemplateFullData
		expected string
		wantNil  bool
	}{
		{
			name: "resource center cover takes precedence",
			data: &model.MarketTemplateFullData{
				CoverURL: " https://r.thingspanel.cn/api/market/templates/assets/covers/template/1.0.0.png ",
				DeviceConfig: &model.DeviceConfigPayload{
					ImageURL: "https://publisher.example/files/old.png",
				},
			},
			expected: "https://r.thingspanel.cn/api/market/templates/assets/covers/template/1.0.0.png",
		},
		{
			name: "legacy device config image is used as fallback",
			data: &model.MarketTemplateFullData{
				DeviceConfig: &model.DeviceConfigPayload{
					ImageURL: " https://publisher.example/files/legacy.png ",
				},
			},
			expected: "https://publisher.example/files/legacy.png",
		},
		{
			name:    "missing image remains empty",
			data:    &model.MarketTemplateFullData{DeviceConfig: &model.DeviceConfigPayload{}},
			wantNil: true,
		},
		{
			name:    "nil payload remains empty",
			data:    nil,
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := resolveMarketTemplateImageURL(tt.data)
			if tt.wantNil {
				if actual != nil {
					t.Fatalf("expected nil, got %q", *actual)
				}
				return
			}
			if actual == nil || *actual != tt.expected {
				t.Fatalf("expected %q, got %v", tt.expected, actual)
			}
		})
	}
}

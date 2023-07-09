package api

import (
	"scraper/entity"
	"testing"
)

func TestDownloadAllPhotos(t *testing.T) {
	type args struct {
		userID    string
		apiKey    string
		photos    []entity.Photo
		albumPath string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "Test go routine 10",
			args: args{
				userID:    "129465275@N03",
				apiKey:    "19ec237d5032c9c651a693e6efc49522",
				photos:    []entity.Photo{},
				albumPath: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DownloadAllPhotos(tt.args.userID, tt.args.apiKey, tt.args.photos, tt.args.albumPath)
		})
	}
}

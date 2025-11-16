package gh

import (
	"testing"
)

func TestGetOwnerAndRepo(t *testing.T) {
	tests := map[string]struct {
		name      string
		repo      string
		wantOwner string
		wantRepo  string
		wantErr   bool
	}{
		"valid repo": {
			repo:      "owner/repo",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantErr:   false,
		},
		"invalid - missing slash": {
			repo:    "ownerrepo",
			wantErr: true,
		},
		"invalid - too many parts": {
			repo:    "owner/repo/extra",
			wantErr: true,
		},
		"invalid - empty": {
			repo:    "",
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			owner, repo, err := getOwnerAndRepo(tt.repo)

			if tt.wantErr && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tt.wantErr {
				if owner != tt.wantOwner {
					t.Errorf("Expected owner '%s', got '%s'", tt.wantOwner, owner)
				}

				if repo != tt.wantRepo {
					t.Errorf("Expected repo '%s', got '%s'", tt.wantRepo, repo)
				}
			}
		})
	}
}

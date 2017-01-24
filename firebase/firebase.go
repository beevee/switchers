package firebase

import "github.com/zabawaba99/firego"

// Repository is a generic holder for auth and url parameters
type Repository struct {
	FirebaseURL   string
	FirebaseToken string
	firebase      *firego.Firebase
}

// Start initializes firebase connection
func (r *Repository) Start() error {
	r.firebase = firego.New(r.FirebaseURL, nil)
	r.firebase.Auth(r.FirebaseToken)

	return nil
}

// Stop does nothing
func (r *Repository) Stop() error {
	return nil
}

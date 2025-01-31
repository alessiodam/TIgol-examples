package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

// Configuration
const (
	CLIENT_ID     = "change_me"             // OAuth client ID
	CLIENT_SECRET = "change_me"             // OAuth client secret
	API_BASE_URL  = "https://api.tigol.net" // Base URL for the API
)

// User represents the expected user data fields.
type User struct {
	ID        uint   `json:"id"`
	UUID      string `json:"uuid"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// TokenResponse represents the JSON structure of the token response.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

// getUserData exchanges the authorization code for an access token and then retrieves user data.
func getUserData(code string) (*User, error) {
	// Prepare data for token exchange
	authData := map[string]string{
		"client_id":     CLIENT_ID,
		"client_secret": CLIENT_SECRET,
		"code":          code,
	}
	jsonData, err := json.Marshal(authData)
	if err != nil {
		return nil, err
	}

	// Request access token
	resp, err := http.Post(API_BASE_URL+"/auth/oidc/token", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get token: status %d", resp.StatusCode)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	// Request user data with the access token
	req, err := http.NewRequest("GET", API_BASE_URL+"/auth/v1/user/me", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
	client := &http.Client{}
	userResp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer userResp.Body.Close()
	if userResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user data: status %d", userResp.StatusCode)
	}

	var user User
	if err := json.NewDecoder(userResp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// authorizedHandler handles the /authorized route.
func authorizedHandler(w http.ResponseWriter, r *http.Request) {
	// Get authorization code from query parameters.
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing authorization code", http.StatusBadRequest)
		return
	}

	user, err := getUserData(code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving user data: %v", err), http.StatusInternalServerError)
		return
	}

	// Generate HTML to display user information.
	tmpl := `
		<html>
		<head><title>User Information</title></head>
		<body>
			<h1>User Information</h1>
			<p>TIGol OAuth2.0/OIDC Demo</p>
			<table border="1">
				<tr><th>id</th><td>{{.ID}}</td></tr>
				<tr><th>uuid</th><td>{{.UUID}}</td></tr>
				<tr><th>first_name</th><td>{{.FirstName}}</td></tr>
				<tr><th>last_name</th><td>{{.LastName}}</td></tr>
				<tr><th>username</th><td>{{.Username}}</td></tr>
				<tr><th>email</th><td>{{.Email}}</td></tr>
				<tr><th>bio</th><td>{{.Bio}}</td></tr>
				<tr><th>created_at</th><td>{{.CreatedAt}}</td></tr>
				<tr><th>updated_at</th><td>{{.UpdatedAt}}</td></tr>
			</table>
		</body>
		</html>`

	t, err := template.New("user").Parse(tmpl)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, user)
}

func main() {
	http.HandleFunc("/authorized", authorizedHandler)
	fmt.Println("Server running on http://localhost:5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}

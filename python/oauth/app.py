import requests
from flask import Flask, request

# Configuration
CLIENT_ID = "change_me"  # OAuth client ID
CLIENT_SECRET = "change_me"  # OAuth client secret
API_BASE_URL = "https://api.tigol.net"  # Base URL for the API

app = Flask(__name__)  # Initialize Flask application

def get_user_data(code):
  """
  Exchange authorization code for access token and retrieve user data.
  
  Args:
    code (str): Authorization code received from the OAuth provider.
  
  Returns:
    dict: User data retrieved from the API.
  """
  auth_data = {
    "client_id": CLIENT_ID,
    "client_secret": CLIENT_SECRET,
    "code": code
  }
  
  # Request access token using the authorization code
  response = requests.post(f"{API_BASE_URL}/auth/oidc/token", json=auth_data)
  response.raise_for_status()  # Raise an exception for HTTP errors
  token = response.json()["access_token"]  # Extract access token from response
  
  # Request user data using the access token
  user_data = requests.get(
    f"{API_BASE_URL}/auth/v1/user/me", 
    headers={"Authorization": f"Bearer {token}"}
  ).json()
  
  return user_data

@app.route("/authorized")
def authorized():
  """
  Handle the OAuth2.0/OIDC authorization callback.
  
  Returns:
    str: HTML page displaying user information.
  """
  code = request.args.get("code")  # Get authorization code from query parameters
  if not code:
    return "Missing authorization code", 400  # Return error if code is missing
  
  user_data = get_user_data(code)  # Retrieve user data using the authorization code
  
  # Generate HTML to display user information
  html = "<h1>User Information</h1>"
  html += "<p>TIgol OAuth2.0/OIDC Demo</p>"
  html += "<table border='1'>"
  
  # Add user data to the HTML table
  for key in ['id', 'uuid', 'first_name', 'last_name', 'username', 'email', 'bio', 'created_at', 'updated_at']:
    html += f"<tr><th>{key}</th><td>{user_data.get(key, '')}</td></tr>"
  
  html += "</table>"
  return html

if __name__ == "__main__":
  app.run(debug=True)  # Run the Flask application in debug mode

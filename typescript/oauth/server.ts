import express, { Request, Response } from "express";
import axios from "axios";

// Configuration
const CLIENT_ID = "change_me"; // OAuth client ID
const CLIENT_SECRET = "change_me"; // OAuth client secret
const API_BASE_URL = "https://api.tigol.net"; // Base URL for the API

const app = express();
const PORT = 5000;

// Interface for user data
interface User {
  id?: any;
  uuid: string;
  first_name: string;
  last_name: string;
  username: string;
  email: string;
  bio: string;
  created_at: string;
  updated_at: string;
}

// getUserData exchanges the authorization code for an access token and retrieves user data.
async function getUserData(code: string): Promise<User> {
  // Exchange authorization code for access token.
  const authData = {
    client_id: CLIENT_ID,
    client_secret: CLIENT_SECRET,
    code: code,
  };

  const tokenResponse = await axios.post(`${API_BASE_URL}/auth/oidc/token`, authData, {
    headers: { "Content-Type": "application/json" },
  });

  const token = tokenResponse.data.access_token;
  
  // Retrieve user data using the access token.
  const userResponse = await axios.get(`${API_BASE_URL}/auth/v1/user/me`, {
    headers: { Authorization: `Bearer ${token}` },
  });

  return userResponse.data as User;
}

// Route handler for /authorized.
app.get("/authorized", async (req: Request, res: Response) => {
  const code = req.query.code as string;
  if (!code) {
    res.status(400).send("Missing authorization code");
  }

  try {
    const userData = await getUserData(code);

    // Build an HTML response to display user information.
    const html = `
      <html>
        <head><title>User Information</title></head>
        <body>
          <h1>User Information</h1>
          <p>TIGol OAuth2.0/OIDC Demo</p>
          <table border="1">
            <tr><th>id</th><td>${userData.id ?? ""}</td></tr>
            <tr><th>uuid</th><td>${userData.uuid}</td></tr>
            <tr><th>first_name</th><td>${userData.first_name}</td></tr>
            <tr><th>last_name</th><td>${userData.last_name}</td></tr>
            <tr><th>username</th><td>${userData.username}</td></tr>
            <tr><th>email</th><td>${userData.email}</td></tr>
            <tr><th>bio</th><td>${userData.bio}</td></tr>
            <tr><th>created_at</th><td>${userData.created_at}</td></tr>
            <tr><th>updated_at</th><td>${userData.updated_at}</td></tr>
          </table>
        </body>
      </html>
    `;
    res.send(html);
  } catch (error: any) {
    res.status(500).send(`Error retrieving user data: ${error.message}`);
  }
});

app.listen(PORT, () => {
  console.log(`Server running on http://localhost:${PORT}`);
});

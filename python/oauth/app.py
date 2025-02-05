import os
import dotenv
from flask import Flask, request, redirect, url_for, make_response, render_template
from tigol import TIgolApiClient

dotenv.load_dotenv(".env")

app = Flask(__name__)
app.secret_key = os.environ.get("FLASK_SECRET_KEY", "super-secret-key")

client = TIgolApiClient(
    os.environ.get("TIGOL_CLIENT_ID"),
    os.environ.get("TIGOL_CLIENT_SECRET"),
)

USER_SESSIONS = {}

def retrieve_session():
    session_id = request.cookies.get("session_id")
    if not session_id or session_id not in USER_SESSIONS:
        return None, {}
    return session_id, USER_SESSIONS[session_id]

@app.route("/")
def index():
    return redirect(client.get_authorization_url(redirect_uri=os.environ.get("TIGOL_REDIRECT_URI")))

@app.route("/authorized")
def authorized():
    code = request.args.get("code")
    if not code:
        return "Missing authorization code", 400

    session_id = request.cookies.get("session_id")
    if not session_id:
        session_id = os.urandom(16).hex()

    try:
        token_obj = client.exchange_code_for_token(code=code)
        user_obj = client.get_user(token_obj)

        USER_SESSIONS[session_id] = {"token": token_obj, "user_data": user_obj.__dict__}

        response = make_response(redirect(url_for("display")))
        response.set_cookie("session_id", session_id)
        return response
    except Exception as e:
        return f"Error during authentication: {str(e)}", 500

@app.route("/display")
def display():
    _, session_data = retrieve_session()
    if not session_data.get("user_data"):
        return render_template("error.html", error_message="Session expired. Log in again.")

    user_data = session_data["user_data"]
    return render_template("authorized.html", user_data=user_data)

if __name__ == "__main__":
    app.run(debug=True)

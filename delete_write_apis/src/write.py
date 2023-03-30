"""
APIs for POST methods:
/package
/package/(id)
/authenticate

function names will have method specified if already
in use (other apis) 
"""
from src import app

# FOR TESTING PURPOSES
@app.route("/index", methods=["GET"])
def index():
    """Func description"""
    return "Hello again"

@app.route("/package", methods=["POST"])
def package():
    """Func description"""
    return "Package"

@app.route("/package/<string:id>", methods=["POST"])
def package_id_post(id):
    """Func description"""
    return "Package by ID"

@app.route("/authenticate", methods=["POST"])
def authenticate():
    """Func description"""
    return "Authenticate"

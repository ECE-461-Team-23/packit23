"""
APIs for DELETE methods:
/reset
/package/(id)
/package/byName/(name)
"""
from src import app

# FOR TESTING PURPOSES
@app.route("/", methods=["GET"])
def test():
    """Func description"""
    return "Hello"

@app.route("/reset", methods=["DELETE"])
def reset():
    """Func description"""
    return "Reset"

@app.route("/package/<string:id>", methods=["DELETE"])
def package_id_delete(id):
    """Func description"""
    return "Package by ID"

@app.route("/package/byName/<string:id>", methods=["DELETE"])
def package_byname_delete(id):
    """Func description"""
    return "Package by name"

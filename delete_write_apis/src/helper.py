import base64
import io
import json
import tempfile
import zipfile
from urllib.parse import urlparse

def checkGithubUrl(url: str) -> bool:
    parsed_uri = urlparse(url)
    if parsed_uri.netloc == "github.com":
        return True
    return False

def grabUrl(fileContents: str) -> str:
    # Returns the URL from package.json inside of a base64 encoded zip file
    zip_buffer = io.BytesIO(base64.b64decode(fileContents))
    zf = zipfile.ZipFile(zip_buffer)

    with tempfile.TemporaryDirectory() as dirPath:
        zf.extractall(dirPath)
        with open(dirPath + "/package.json") as file:
            package_data = json.load(file)
            return package_data["repository"]["url"]
            # print(package_data["homepage"])
            # print(package_data["repository"]["url"])

def convertZipToBase64(filePath: str):
    # Utility function only, used to generate test data
    # convertZipToBase64("/path/to/file.zip")
    # python3 helper.py > /path/to/output.txt
    with open(filePath, 'rb') as f:
        b = f.read()
        encoded = base64.b64encode(b)
        print(encoded)

# with open("/Users/ben/code/packit23/delete_write_apis/tests/example_b64.txt", "r") as file:
#     x = grabUrl(file.read())
#     print(x)

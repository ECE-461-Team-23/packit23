import base64
import json
import io
import os
import requests
import tempfile
import zipfile
from git import Repo
from urllib.parse import urlparse

def zipitem(path, ziph):
    # ziph is zipfile handle
    for root, dirs, files in os.walk(path):
        for file in files:
            ziph.write(os.path.join(root, file),
                       os.path.relpath(os.path.join(root, file),
                                       os.path.join(path, '.')))

def zipit(source_dir, zip_name):
    zipf = zipfile.ZipFile(zip_name, 'w', zipfile.ZIP_DEFLATED)
    zipitem(source_dir, zipf)
    zipf.close()

def downloadGithubRepo(url: str):
    with tempfile.TemporaryDirectory() as dir:
        Repo.clone_from(url + ".git", dir)

        # Create zip file
        zipit(dir, 'package.zip')

        # Convert to b64 
        with open('package.zip', 'rb') as f:
            b = f.read()
            return base64.b64encode(b)

def checkGithubUrl(url: str) -> bool:
    parsed_uri = urlparse(url)
    if parsed_uri.netloc == "github.com":
        return True
    return False

def grabPackageDataFromZip(fileContents: str) -> tuple[str, str, str]:
    # Returns the URL from package.json inside of a base64 encoded zip file
    zip_buffer = io.BytesIO(base64.b64decode(fileContents))
    zf = zipfile.ZipFile(zip_buffer)

    with tempfile.TemporaryDirectory() as dirPath:
        zf.extractall(dirPath)
        with open(dirPath + "/package.json") as file:
            package_data = json.load(file)
            return package_data["name"], package_data["version"], package_data["homepage"]
            # print(package_data["homepage"])
            # print(package_data["repository"]["url"])

# def convertZipToBase64(filePath: str):
#     # Utility function only, used to generate test data
#     # convertZipToBase64("/path/to/file.zip")
#     # python3 helper.py > /path/to/output.txt
#     with open(filePath, 'rb') as f:
#         b = f.read()
#         encoded = base64.b64encode(b)
#         print(encoded)

def getOwnerAndRepoFromURL(url: str) -> tuple[str, str]:
    # Returns owner, repo from URL
    parseResult = urlparse(url)
    path = parseResult.path.split("/")
    return path[1], path[2]

def grabPackageDataFromURL(url: str) -> tuple[str, str, str]:
    token = os.environ["GITHUB_TOKEN"]
    headers = { "Accept": "application/vnd.github.raw",
                "Authorization": f"Bearer {token}",
                "X-GitHub-Api-Version": "2022-11-28"}
    timeout = 60
    owner, repo = getOwnerAndRepoFromURL(url)
    response = requests.get(url=f"https://api.github.com/repos/{owner}/{repo}/contents/package.json", headers=headers, timeout=timeout)
    response_text = response.text
    assert response_text != None and response_text != ""
    package_data = json.loads(response_text)
    return package_data["name"], package_data["version"], package_data["homepage"]


def grabPackageDataFromRequest(parsed_body):
    if "Content" in parsed_body:       
        # Package contents. This is the zip file uploaded by the user. (Encoded as text using a Base64 encoding).
        # This will be a zipped version of an npm package's GitHub repository, minus the ".git/" directory." It will, for example, include the "package.json" file that can be used to retrieve the project homepage.
        # See https://docs.npmjs.com/cli/v7/configuring-npm/package-json#homepage.
        return grabPackageDataFromZip(parsed_body["Content"])
    elif "URL" in parsed_body:
        # Ingest package from public URL
        return grabPackageDataFromURL(parsed_body["URL"])

# with open("/Users/ben/code/packit23/delete_write_apis/tests/example_b64.txt", "r") as file:
#     x = grabUrl(file.read())
#     print(x)

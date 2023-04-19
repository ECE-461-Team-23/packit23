import os
from google.cloud import storage

def upload_b64_blob(contents, destination_blob_name):
    """Uploads a file to the bucket."""
    storage_client = storage.Client()
    bucket = storage_client.bucket(os.environ["BUCKET_NAME"])
    blob = bucket.blob(destination_blob_name)

    blob.upload_from_string(contents)
    print(f"File uploaded to {destination_blob_name}.")

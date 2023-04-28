import os
from google.cloud import storage

from . import helper

def upload_b64_blob(contents, destination_blob_name):
    """Uploads a file to the bucket."""
    storage_client = storage.Client()
    bucket = storage_client.bucket(os.environ["BUCKET_NAME"])
    blob = bucket.blob(destination_blob_name)

    blob.upload_from_string(contents)
    helper.log(f"File uploaded to {destination_blob_name}.")


def delete_blob(blob_name):
    """Deletes a blob from the bucket."""
    storage_client = storage.Client()
    bucket = storage_client.bucket(os.environ["BUCKET_NAME"])
    blob = bucket.blob(blob_name)
    blob.delete()
    helper.log(f"Blob {blob_name} deleted.")

def empty_bucket():
    """Empty all objects from a bucket"""
    storage_client = storage.Client()
    bucket = storage_client.bucket(os.environ["BUCKET_NAME"])
    blobs = bucket.list_blobs()
    for blob in blobs: 
        blob.delete()
    helper.log(f"Bucket emptied")

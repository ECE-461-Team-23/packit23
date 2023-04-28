import json
import os
import unittest
from unittest.mock import call

from mock import patch, Mock
from fastapi.testclient import TestClient

from src.main import app
from src.authentication import generate_jwt, get_hashed_password, validate_jwt


client = TestClient(app)
mock_jwt_secret = "mockjwtsecret"
mock_userid = "mockuserid"
mock_packagename = "mockpackagename"
mock_packageversion = "1.1.1"
mock_packageid = 10
mock_username = "mockusername"
mock_password = "mockpassword"
mock_stored_password = get_hashed_password(mock_password)
mock_packageurl = "https://www.github.com/test/url"
mock_packagerater_url = "mockurl.com"

class TestBasicConnections(unittest.TestCase):

    # @patch('src.database.connect_with_connector')
    def test_main_hello(self):
        # connect_with_connector_mock.return_value = ''
        response = client.get("/")
        assert response.status_code == 200
        assert response.json() == {"Hello": "World"}

    def test_write_hello(self):
        response = client.get("/write")
        assert response.status_code == 200
        assert response.json() == {"Hello": "Write"}

    def test_delete_hello(self):
        response = client.get("/delete")
        assert response.status_code == 200
        assert response.json() == {"Hello": "Delete"}

class TestAuthorization(unittest.TestCase):
    def test_reset_noauth(self):
        response = client.delete("/reset")
        assert response.status_code == 400
        assert response.text == "There is missing field(s) in the AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."

    def test_delete_byname_noauth(self):
        response = client.delete("/package/byName/testName")
        assert response.status_code == 400
        assert response.text == "There is missing field(s) in the PackageName/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."

    def test_delete_byid_noauth(self):
        response = client.delete("/package/1")
        assert response.status_code == 400
        assert response.text == "There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."

    def test_create_noauth(self):
        response = client.post("/package/")
        assert response.status_code == 400
        assert response.text == "There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly (e.g. Content and URL are both set), or the AuthenticationToken is invalid."

    def test_update_noauth(self):
        response = client.put("/package/1")
        assert response.status_code == 400
        assert response.text == "There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."


class TestDeleteEndpoints(unittest.TestCase):
    def setUp(self) -> None:
        os.environ["JWT_SECRET"] = mock_jwt_secret
        self.my_jwt = generate_jwt(mock_userid)
        return super().setUp()

    @patch('src.database.reset_database') #Override functions
    @patch('src.bucket.empty_bucket')
    def test_reset_successs(self, reset_database, empty_bucket):
        headers = {"X-Authorization": self.my_jwt}
        response = client.delete("/reset", headers=headers)
        assert response.status_code == 200

    @patch('src.delete.bucket')    
    @patch('src.delete.database')
    def test_delete_byname(self, database, bucket):
        database.get_all_versions_of_package.return_value = [1, 2]
        database.delete_package.side_effect = iter([3, 4])

        headers = {"X-Authorization": self.my_jwt}
        response = client.delete(f"/package/byName/{mock_packagename}", headers=headers)
        assert response.status_code == 200
        # database.delete_package.assert_called_with(1)
        database.delete_package.assert_has_calls([call(1), call(2)])
        bucket.delete_blob.assert_has_calls([call("3"), call("4")])

    @patch('src.delete.bucket')    
    @patch('src.delete.database')
    def test_delete_byid(self, database, bucket):
        database.check_if_package_exists.return_value = True
        database.delete_package.return_value = 1

        headers = {"X-Authorization": self.my_jwt}
        response = client.delete(f"/package/{mock_packageid}", headers=headers)
        assert response.status_code == 200
        database.check_if_package_exists.assert_called_with(mock_packageid)
        database.delete_package.assert_called_with(mock_packageid)
        bucket.delete_blob.assert_called_with("1")


class TestWriteEndpoints(unittest.TestCase):
    def setUp(self) -> None:
        os.environ["JWT_SECRET"] = mock_jwt_secret
        os.environ["PACKAGE_RATER_URL"] = mock_packagerater_url
        self.my_jwt = generate_jwt(mock_userid)
        return super().setUp()
  
    @patch('src.write.database')
    def test_authenticate_success(self, database):
        database.get_data_for_user.return_value = (mock_userid, mock_username, mock_stored_password)

        body = {"User": {"name": mock_username, "isAdmin": True}, "Secret": {"password": mock_password}}
        response = client.put(f"/authenticate", content=json.dumps(body))
        assert response.status_code == 200
        assert validate_jwt(response.text) == mock_userid

    @patch('src.write.database')
    def test_authenticate_badlogin(self, database):
        database.get_data_for_user.return_value = None, None, None

        body = {"User": {"name": mock_username, "isAdmin": True}, "Secret": {"password": mock_password}}
        response = client.put(f"/authenticate", content=json.dumps(body))
        assert response.status_code == 401

    @patch('requests.post')
    @patch('src.write.helper')
    @patch('src.write.database')
    def test_package_create_success(self, database, helper, requests):
        helper.grabPackageDataFromRequest.return_value = mock_packagename, mock_packageversion, mock_packageurl
        database.get_package_id.return_value = None
        requests.return_value = Mock(status_code=200, text='{"NET_SCORE":1}')
        helper.downloadGithubRepo.return_value = "filler"
        database.upload_package.return_value = 1

        headers = {"X-Authorization": self.my_jwt}
        body = {"URL": mock_packageurl}
        response = client.post("/package/", headers=headers, content=json.dumps(body))
        assert response.status_code == 201

    @patch('requests.post')
    @patch('src.write.helper')
    @patch('src.write.database')
    def test_package_update_success(self, database, helper, requests):
        helper.grabPackageDataFromRequest.return_value = mock_packagename, mock_packageversion, mock_packageurl
        database.get_package_id.return_value = mock_packageid
        requests.return_value = Mock(status_code=200, text='{"NET_SCORE":1}')
        helper.downloadGithubRepo.return_value = "filler"
        database.upload_package.return_value = 1

        headers = {"X-Authorization": self.my_jwt}
        body = {
                "metadata": {
                    "Name": mock_packagename,
                    "Version": mock_packageversion,
                    "ID": str(mock_packageid)
                },
                "data": {
                    "URL": mock_packageurl
                }
            }
        response = client.put(f"/package/{mock_packageid}", headers=headers, content=json.dumps(body))
        assert response.status_code == 200
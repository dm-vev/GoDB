import requests

class GemClient:
    def __init__(self, base_url, api_key):
        self.base_url = base_url
        self.headers = {"X-API-Key": api_key}

    def create_table(self, name):
        url = f"{self.base_url}/createTable"
        data = {"name": name}
        response = requests.post(url, json=data, headers=self.headers)
        response.raise_for_status()

    def delete_table(self, name):
        url = f"{self.base_url}/deleteTable"
        data = {"name": name}
        response = requests.post(url, json=data, headers=self.headers)
        response.raise_for_status()

    def set(self, table_name, key, value):
        url = f"{self.base_url}/set"
        data = {"table_name": table_name, "key": key, "value": value}
        response = requests.post(url, json=data, headers=self.headers)
        response.raise_for_status()

    def get(self, table_name, key):
        url = f"{self.base_url}/get"
        params = {"table_name": table_name, "key": key}
        response = requests.get(url, params=params, headers=self.headers)
        response.raise_for_status()
        return response.json()["value"]

    def delete(self, table_name, key):
        url = f"{self.base_url}/delete"
        data = {"table_name": table_name, "key": key}
        response = requests.post(url, json=data, headers=self.headers)
        response.raise_for_status()
    def exportDB(self, filename):
        url = f"{self.base_url}/exportToFile"
        data = {"filename": filename}
        response = requests.post(url, json=data, headers=self.headers)
        response.raise_for_status()

api_key = "YOUR_API_KEY"
base_url = "http://localhost:8080"

client = GemClient(base_url, api_key)

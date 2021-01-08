import uuid

from locust import HttpUser, SequentialTaskSet, task, between

class ChattoUser(SequentialTaskSet):
    def on_start(self):
        self.chatto_name = str(uuid.uuid4())

    @task
    def hi(self):
        self.client.post("/endpoints/rest", json={
            "sender": self.chatto_name,
            "text": "hi",
        })

    @task
    def bad(self):
        self.client.post("/endpoints/rest", json={
            "sender": self.chatto_name,
            "text": "bad",
        })

    @task
    def yes(self):
        self.client.post("/endpoints/rest", json={
            "sender": self.chatto_name,
            "text": "yes",
        })


class User(HttpUser):
    tasks = [ChattoUser]
    wait_time = between(1, 2)

import pytest
from requests import Session
from urllib.parse import urljoin
from tests.fixtures import ACCT_NAME


class LiveServerSession(Session):
    def __init__(self, base_url=None):
        super().__init__()
        self.base_url = base_url

    def request(self, method, url, *args, **kwargs):
        joined_url = urljoin(self.base_url, url)
        return super().request(method, joined_url, *args, **kwargs)


@pytest.fixture
def client():
    return LiveServerSession("http://localhost:3000")


@pytest.fixture
def actor_url(client):
    resp = client.get(f"/.well-known/webfinger?resource={ACCT_NAME}")
    data = resp.json()
    self_link = [x for x in data["links"] if x["rel"] == "self"][0]
    return self_link["href"]


@pytest.fixture
def actor(client):
    resp = client.get(actor_url, headers={"Content-Type": "application/activity+json"})
    assert resp.status_code == 200
    return resp.json()


@pytest.fixture
def inbox(actor):
    return actor["inbox"]


@pytest.fixture
def outbox(actor):
    return actor["outbox"]


@pytest.fixture
def followers(actor):
    return actor["followers"]

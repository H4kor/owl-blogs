import pytest
from .fixtures import ACCT_NAME


@pytest.mark.parametrize(
    ["query", "status"],
    [
        ["", 404],
        ["?foo=bar", 404],
        ["?resource=lol@bar.com", 404],
        [f"?resource={ACCT_NAME}", 200],
    ],
)
def test_webfinger_status(client, query, status):
    resp = client.get("/.well-known/webfinger" + query)
    assert resp.status_code == status


def test_webfinger(client):
    resp = client.get(f"/.well-known/webfinger?resource={ACCT_NAME}")
    assert resp.status_code == 200
    data = resp.json()
    assert data["subject"] == ACCT_NAME
    assert len(data["links"]) > 0
    self_link = [x for x in data["links"] if x["rel"] == "self"][0]
    assert self_link["type"] == "application/activity+json"
    assert "href" in self_link

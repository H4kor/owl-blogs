import pytest


def test_actor(client, actor_url):
    resp = client.get(actor_url, headers={"Content-Type": "application/activity+json"})
    assert resp.status_code == 200
    data = resp.json()
    assert "id" in data
    assert "type" in data
    assert "inbox" in data
    assert "outbox" in data
    assert "followers" in data
    assert "preferredUsername" in data
    assert "publicKey" in data
    assert len(data["publicKey"])

    pubKey = data["publicKey"]
    assert "id" in pubKey
    assert "owner" in pubKey
    assert "publicKeyPem" in pubKey

    assert pubKey["owner"] == data["id"]
    assert pubKey["id"] != data["id"]
    assert "-----BEGIN RSA PUBLIC KEY-----" in pubKey["publicKeyPem"]

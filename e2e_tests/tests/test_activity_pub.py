from pprint import pprint
from time import sleep
import uuid
import requests
from .fixtures import ensure_follow, msg_inc, sign
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


def test_following(client, inbox_url, followers_url, actor_url):
    with msg_inc(1):
        req = sign(
            "POST",
            inbox_url,
            {
                "@context": "https://www.w3.org/ns/activitystreams",
                "id": "http://mock_masto:8000/d0b5768b-a15b-4ed6-bc84-84c7e2b57588",
                "type": "Follow",
                "actor": "http://mock_masto:8000/users/h4kor",
                "object": actor_url,
            },
        )
        resp = requests.Session().send(req)

        assert resp.status_code == 200

        resp = client.get(
            followers_url, headers={"Content-Type": "application/activity+json"}
        )
        assert resp.status_code == 200
        data = resp.json()
        pprint(data)
        assert "items" in data
        assert len(data["items"]) == 1


def test_unfollow(client, inbox_url, followers_url, actor_url):
    ensure_follow(client, inbox_url, actor_url)
    sleep(0.5)
    with msg_inc(1):
        req = sign(
            "POST",
            inbox_url,
            {
                "@context": "https://www.w3.org/ns/activitystreams",
                "id": "http://mock_masto:8000/users/h4kor#follows/3632040/undo",
                "type": "Undo",
                "actor": "http://mock_masto:8000/users/h4kor",
                "object": {
                    "id": "http://mock_masto:8000/d0b5768b-a15b-4ed6-bc84-84c7e2b57588",
                    "type": "Follow",
                    "actor": "http://mock_masto:8000/users/h4kor",
                    "object": actor_url,
                },
            },
        )
        resp = requests.Session().send(req)
        assert resp.status_code == 200

        resp = client.get(
            followers_url, headers={"Content-Type": "application/activity+json"}
        )
        assert resp.status_code == 200
        data = resp.json()
        pprint(data)
        assert "totalItems" in data
        assert data["totalItems"] == 0


def test_status_code_unknown_post(client, inbox_url, followers_url, actor_url):
    req = sign(
        "POST",
        inbox_url,
        {
            "@context": "https://www.w3.org/ns/activitystreams",
            "id": f"http://mock_masto:8000/users/h4kor#like-{uuid.uuid4()}",
            "type": "Like",
            "actor": "http://mock_masto:8000/users/h4kor",
            "object": "http://localhost:3000/post/foobar/",
        },
    )
    resp = requests.Session().send(req)
    assert resp.status_code == 404
    assert resp.json()["error"] == "entry not found"


def test_status_code_unsigned(client, inbox_url, followers_url, actor_url):
    resp = requests.post(
        inbox_url,
        json={
            "@context": "https://www.w3.org/ns/activitystreams",
            "id": f"http://mock_masto:8000/users/h4kor#like-{uuid.uuid4()}",
            "type": "Like",
            "actor": "http://mock_masto:8000/users/h4kor",
            "object": "http://localhost:3000/post/foobar/",
        },
    )
    assert resp.status_code == 403
    assert resp.json()["error"] == "cannot verify signature"

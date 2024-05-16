from pprint import pprint
import requests
from .fixtures import ensure_follow, sign
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
    req = sign(
        "POST",
        inbox_url,
        {
            "@context": "https://www.w3.org/ns/activitystreams",
            "id": "https://mock_masto/d0b5768b-a15b-4ed6-bc84-84c7e2b57588",
            "type": "Follow",
            "actor": "http://mock_masto:8000/users/h4kor",
            "object": actor_url,
        },
    )
    pprint(req.headers)
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


# def test_unfollow(client, inbox_url, followers_url, actor_url):
#     ensure_follow(client, inbox_url, actor_url)

#     resp = client.post(
#         inbox_url,
#         json={
#             "@context": "https://www.w3.org/ns/activitystreams",
#             "id": "http://mock_masto:8000/users/h4kor#follows/3632040/undo",
#             "type": "Undo",
#             "actor": "http://mock_masto:8000/users/h4kor",
#             "object": {
#                 "id": "https://mock_masto/d0b5768b-a15b-4ed6-bc84-84c7e2b57588",
#                 "type": "Follow",
#                 "actor": "http://mock_masto:8000/users/h4kor",
#                 "object": actor_url,
#             },
#         },
#         headers={"Content-Type": "application/activity+json"},
#     )
#     assert resp.status_code == 200

#     resp = client.get(
#         followers_url, headers={"Content-Type": "application/activity+json"}
#     )
#     assert resp.status_code == 200
#     data = resp.json()
#     assert "items" in data
#     assert len(data["items"]) == 0

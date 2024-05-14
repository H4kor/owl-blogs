import json
from flask import Flask, request

app = Flask(__name__)


PRIV_KEY_PEM = """-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQCni8P4bvkC/3Sx
NTrDw1qw0vWtJKZMsyJ3Mcs4+1apoVqOQhujUqGqFSiRT7Vmc7OEhB0vikdiTkCk
1PcoTd/rOa/0WbG8385JcLzaJfTIG+rrRtHwZ1TwxwBju43jpGEZjpbA0dqoXMkr
J1MyD7aPLoAiVe0ikw2czSZumv4ncemOtk0VG3b2TnIxo3CMKtUOWu8xT08MMIuo
3cZRnpI6Xr/ULzvjv8e3EjIpwRJqMPECtGsfdcjFmR0yFIrjrlmkZTiW31z/Dk7i
xRGD0ADy3WEQ3lA4l3mNZeyG4S0Wi4iYe9/wegESMZcakLoME7ks+KNS388Mdbcd
DKy9NmWvAgMBAAECggEABLQAA0hHhdWv6+Lc9xkpFuTvxTV4fuyvCf4u1eGlnstg
ZF/nW1/6w8XQ8WCgbJ4mKuZz1J14FYKxfoRaj8S9MA2Ff+wd+M77gRpAuDWajRzO
LQk8OW2yd7POXKkAzvln9F9eofkCFKR4zSpPGTenCJaQkuYrQEOKfUf7oofdRzQi
w9kmp3wAxM/EseHZpknYDCgDQV7MDQAaMD7kbynL2WfXPxebktwpRlKUwgtGrevj
gagQL8J/GX6wO3ymw9sln4BhlI2+3LuiMXQdQc1tamkXFCguCuOZCu/2VRdCHmiS
nnpu+FMspBHbvxO+RXo3Cu/S6jjJgoQxD2WZTE0gqQKBgQDM6AQdqBYjISdkI9Gl
6ZLLjwZRJSYpopujtX7pun61l9kUwQevaR2Z39rMWxX62DD6arazi/ygIUBw6Kgp
s/qBEb29ec+0cESdC8aJYb3dGvDzh/8C05p7ozxj8JZQcxq5W5jql/BELlSsUONO
jfqQv8RGZNSkD9uy6TxOr4eWIwKBgQDRUuO/XRDLt8Mp10mTshxTznSQ3gAJYKeG
0WfEC3kPEukHBQb8huqFcQDiQ71oBWuEdOQWgT3aBS6L+nIMyZMT5u+BejQm7/E5
pMM+z0VRpfFSsIrCvU8yKam0aemQGlKQAfhTct1gCg+wKnYsSQMlNHKWEfDbw9I/
cns/IN+dBQKBgQC6/Of0oFVDTZgC3GUPAO3C8QwUtM/0or1hUdk1Nck3shCZzeVT
f5tRtmSWpHCUbwGTJBsCEjdBcda6srXzCJkLe8Moy6ZtxR34KqzM5fM7eMB1nJ9s
Vunc9gPAN+cUF1ZF3H7ZZjoOHjGK5m3oW8xSl41np9Acv5P/2rP8Ilaa/QKBgQDJ
YwISfitGk8mEW8hB/L4cMykapztJyl/i6Vz31EHoKr1fL4sFMZg4QfwjtCBqD6zd
hshajoU/WHTr30wS2WxTXX9YBoZeX8KpPsdJioiagRioAYm+yfuDu2m2VZ+MMIb2
Xa7YOk6Zs5RcXL3M5YHNLaSAlUoxZTjGKhJBLhN1MQKBgQCbo3ngBl7Qjjx4WJ93
2WEEKvSDCv69eecNQDuKWKEiFqBN23LheNrN8DXMWFTtE4miY106dzQ0dUMh418x
K98rXSX3VvY4w48AznvPMKVLqesFjcvwnBdvk/NqXod20CMSpOEVj6W/nGoTBQt2
0PuW3IUym9KvO0WX9E+1Qw8mbw==
-----END PRIVATE KEY-----"""

PUB_KEY_PEM = """-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAp4vD+G75Av90sTU6w8Na
sNL1rSSmTLMidzHLOPtWqaFajkIbo1KhqhUokU+1ZnOzhIQdL4pHYk5ApNT3KE3f
6zmv9FmxvN/OSXC82iX0yBvq60bR8GdU8McAY7uN46RhGY6WwNHaqFzJKydTMg+2
jy6AIlXtIpMNnM0mbpr+J3HpjrZNFRt29k5yMaNwjCrVDlrvMU9PDDCLqN3GUZ6S
Ol6/1C8747/HtxIyKcESajDxArRrH3XIxZkdMhSK465ZpGU4lt9c/w5O4sURg9AA
8t1hEN5QOJd5jWXshuEtFouImHvf8HoBEjGXGpC6DBO5LPijUt/PDHW3HQysvTZl
rwIDAQAB
-----END PUBLIC KEY-----"""

INBOX = []


@app.route("/.well-known/webfinger")
def webfinger():
    return json.dumps(
        {
            "subject": "acct:h4kor@mock_masto",
            "aliases": [
                "http://mock_masto/@h4kor",
                "http://mock_masto/users/h4kor",
            ],
            "links": [
                {
                    "rel": "http://webfinger.net/rel/profile-page",
                    "type": "text/html",
                    "href": "http://mock_masto/@h4kor",
                },
                {
                    "rel": "self",
                    "type": "application/activity+json",
                    "href": "http://mock_masto/users/h4kor",
                },
                {
                    "rel": "http://ostatus.org/schema/1.0/subscribe",
                    "template": "http://mock_masto/authorize_interaction?uri={uri}",
                },
                {
                    "rel": "http://webfinger.net/rel/avatar",
                    "type": "image/png",
                    "href": "http://assets.mock_masto/accounts/avatars/000/082/056/original/a4be9944e3b03229.png",
                },
            ],
        }
    )


@app.route("/users/h4kor")
def actor():
    return json.dumps(
        {
            "@context": [
                "http://www.w3.org/ns/activitystreams",
                "http://w3id.org/security/v1",
                {
                    "manuallyApprovesFollowers": "as:manuallyApprovesFollowers",
                    "toot": "http://joinmastodon.org/ns#",
                    "featured": {"@id": "toot:featured", "@type": "@id"},
                    "featuredTags": {"@id": "toot:featuredTags", "@type": "@id"},
                    "alsoKnownAs": {"@id": "as:alsoKnownAs", "@type": "@id"},
                    "movedTo": {"@id": "as:movedTo", "@type": "@id"},
                    "schema": "http://schema.org#",
                    "PropertyValue": "schema:PropertyValue",
                    "value": "schema:value",
                    "discoverable": "toot:discoverable",
                    "Device": "toot:Device",
                    "Ed25519Signature": "toot:Ed25519Signature",
                    "Ed25519Key": "toot:Ed25519Key",
                    "Curve25519Key": "toot:Curve25519Key",
                    "EncryptedMessage": "toot:EncryptedMessage",
                    "publicKeyBase64": "toot:publicKeyBase64",
                    "deviceId": "toot:deviceId",
                    "claim": {"@type": "@id", "@id": "toot:claim"},
                    "fingerprintKey": {"@type": "@id", "@id": "toot:fingerprintKey"},
                    "identityKey": {"@type": "@id", "@id": "toot:identityKey"},
                    "devices": {"@type": "@id", "@id": "toot:devices"},
                    "messageFranking": "toot:messageFranking",
                    "messageType": "toot:messageType",
                    "cipherText": "toot:cipherText",
                    "suspended": "toot:suspended",
                    "memorial": "toot:memorial",
                    "indexable": "toot:indexable",
                    "Hashtag": "as:Hashtag",
                    "focalPoint": {"@container": "@list", "@id": "toot:focalPoint"},
                },
            ],
            "id": "http://mock_masto/users/h4kor",
            "type": "Person",
            "following": "http://mock_masto/users/h4kor/following",
            "followers": "http://mock_masto/users/h4kor/followers",
            "inbox": "http://mock_masto/users/h4kor/inbox",
            "outbox": "http://mock_masto/users/h4kor/outbox",
            "featured": "http://mock_masto/users/h4kor/collections/featured",
            "featuredTags": "http://mock_masto/users/h4kor/collections/tags",
            "preferredUsername": "h4kor",
            "name": "Niko",
            "summary": '<p>Teaching computers to do things with arguable efficiency.</p><p>he/him</p><p><a href="http://mock_masto/tags/vegan" class="mention hashtag" rel="tag">#<span>vegan</span></a> <a href="http://mock_masto/tags/cooking" class="mention hashtag" rel="tag">#<span>cooking</span></a> <a href="http://mock_masto/tags/programming" class="mention hashtag" rel="tag">#<span>programming</span></a> <a href="http://mock_masto/tags/politics" class="mention hashtag" rel="tag">#<span>politics</span></a> <a href="http://mock_masto/tags/climate" class="mention hashtag" rel="tag">#<span>climate</span></a></p>',
            "url": "http://mock_masto/@h4kor",
            "manuallyApprovesFollowers": False,
            "discoverable": True,
            "indexable": False,
            "published": "2018-08-16T00:00:00Z",
            "memorial": False,
            "devices": "http://mock_masto/users/h4kor/collections/devices",
            "publicKey": {
                "id": "http://mock_masto/users/h4kor#main-key",
                "owner": "http://mock_masto/users/h4kor",
                "publicKeyPem": PUB_KEY_PEM,
            },
            "tag": [
                {
                    "type": "Hashtag",
                    "href": "http://mock_masto/tags/politics",
                    "name": "#politics",
                },
                {
                    "type": "Hashtag",
                    "href": "http://mock_masto/tags/climate",
                    "name": "#climate",
                },
                {
                    "type": "Hashtag",
                    "href": "http://mock_masto/tags/vegan",
                    "name": "#vegan",
                },
                {
                    "type": "Hashtag",
                    "href": "http://mock_masto/tags/programming",
                    "name": "#programming",
                },
                {
                    "type": "Hashtag",
                    "href": "http://mock_masto/tags/cooking",
                    "name": "#cooking",
                },
            ],
            "attachment": [
                {
                    "type": "PropertyValue",
                    "name": "Me",
                    "value": '<a href="http://rerere.org" target="_blank" rel="nofollow noopener noreferrer me" translate="no"><span class="invisible">http://</span><span class="">rerere.org</span><span class="invisible"></span></a>',
                },
                {
                    "type": "PropertyValue",
                    "name": "Blog",
                    "value": '<a href="http://blog.libove.org" target="_blank" rel="nofollow noopener noreferrer me" translate="no"><span class="invisible">http://</span><span class="">blog.libove.org</span><span class="invisible"></span></a>',
                },
                {"type": "PropertyValue", "name": "Location", "value": "Münster"},
                {
                    "type": "PropertyValue",
                    "name": "Current Project",
                    "value": '<a href="http://git.libove.org/h4kor/owl-blogs" target="_blank" rel="nofollow noopener noreferrer me" translate="no"><span class="invisible">http://</span><span class="">git.libove.org/h4kor/owl-blogs</span><span class="invisible"></span></a>',
                },
            ],
            "endpoints": {"sharedInbox": "http://mock_masto/inbox"},
            "icon": {
                "type": "Image",
                "mediaType": "image/png",
                "url": "http://assets.mock_masto/accounts/avatars/000/082/056/original/a4be9944e3b03229.png",
            },
        }
    )


@app.route("/users/h4kor/inbox")
def inbox():
    if request.method == "POST":
        INBOX.append(request.get_json())
    return ""


if __name__ == "__main__":
    app.run(debug=True, host="0.0.0.0", port="8000")

# Security

Chatto has many security options, which you can use to protect your bot and its data.

## Bot

You can add an authorization token to the [`/bot` endpoints](/endpoints) in the **bot.yml** file like this:

```yaml
extensions:
  my_server:
    type: REST
    url: http://localhost:8770

store:
  type: REDIS
  host: localhost
  password: pass
  ttl: 3600

auth:                           # or leave empty and use the environment
  token: this-is-a-bot-token    # variable CHATTO_BOT_AUTH_TOKEN
```

If a token is provided, requests to `/bot/predict` and `/bot/senders/<sender_id>` will require the token in the `Authorization` header as Bearer Token.

## REST Channel

Similarly, you can secure the [`/channels/rest`](/endpoints) endpoint by declaring a token in **chn.yml** like this:

```yaml
telegram:
  bot_key:

twilio:
  account_sid:
  auth_token:
  number:

rest:                           # or leave empty and use the environment 
  token: this-is-a-rest-token   # variable CHATTO_CHN_REST_TOKEN
```

If a token is provided, requests to `/channels/rest` will require the token in the `Authorization` header as Bearer Token.

## REST Extensions

### Token

REST extensions can also require a token. The `extensions.ServeREST` function reads a `token` flag. For example:

```bash
go run examples/04_trivia/ext/ext.go -token my-extension-authorization-token
```

Then, in the **bot.yml** file, in `extensions`, include the token:

```yaml
extensions:
  my_server:
    type: REST
    url: http://localhost:8770
    token: my-extension-authorization-token   # or leave empty and use the environment 
                                              # variable CHATTO_BOT_EXTENSIONS_MY_SERVER_TOKEN
```

### SSL

REST extensions can be served over HTTPS. The `extensions.ServeREST` function reads the `ssl-key` and `ssl-cert` flags. For example:

```bash
go run examples/04_trivia/ext/ext.go -ssl-key localhost.key -ssl-cert localhost.crt
```

Then, in the **bot.yml** file, in `extensions`, use HTTPS:

```yaml
extensions:
  my_secure_server:
    type: REST
    url: https://localhost:8770
```

For local testing you can follow these steps to generate a certificate for *localhost*:

```bash
openssl req -x509 -out localhost.crt -keyout localhost.key \
  -newkey rsa:2048 -nodes -sha256 \
  -subj '/CN=localhost' -extensions EXT -config <( \
   printf "[dn]\nCN=localhost\n[req]\ndistinguished_name = dn\n[EXT]\nsubjectAltName=DNS:localhost\nkeyUsage=digitalSignature\nextendedKeyUsage=serverAuth")
sudo cp localhost.crt /usr/local/share/ca-certificates/
sudo dpkg-reconfigure ca-certificates
```

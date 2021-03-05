# Endpoints

Apart from the [channels](/channels), once your Chatto bot is running, you can interact with it via the HTTP endpoints.

## HTTP Endpoint

Chatto will run on port `4770` by default. You can specify a different one with the `-port` flag.

Send a `POST` request to `/channels/rest` with the following body structure:

```json
{
    "sender": "foo",
    "text": "bar"
}
```

Example with cURL:

```bash
curl --request POST 'http://localhost:4770/channels/rest' \
--header 'Content-Type: application/json' \
--data-raw '{
    "sender": "foo",
    "text": "bar"
}'
```

The bot will respond as such:

```json
[
    {
        "text": "some answer",
        "image": "some image"
    }
]
```

## CLI

You can use the Chatto CLI tool by downloading the `chatto-cli` binary, which launches a command line interface where you can send and receive messages from your bot. This is a useful mode when debugging.

<p align="center">
<img src="https://media.giphy.com/media/DFIxYClozxyMg9wnil/giphy.gif" alt="demo" width="480"/>
</p>

## Prediction

You can test the command predictions with a `POST` request to the `/bot/predict` endpoint with the following body structure:

```json
{
    "text": "good"
}
```

Example with cURL:

```bash
curl --request POST 'http://localhost:4770/bot/predict' \
--header 'Content-Type: application/json' \
--data-raw '{
    "text": "foo"
}'
```

The resulting prediction will look like this:

```json
{
    "original": "foo",
    "predicted": "good",
    "probability": 0.3274336283185841
}
```

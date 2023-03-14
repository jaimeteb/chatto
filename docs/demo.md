# Demo

In this chat window you can interact with your chatbot running on localhost, or with another Chatto example bot.

## Test it out

<div style="min-height:300px;">
    <iframe
        src="https://chatto-examples-ui.web.app/"
        style="width:100%;min-height:300px;"
        frameborder="0"
        marginwidth="0"
    ></iframe>
</div>

## Local REST Endpoint

To use this demo, run a Chatto bot on your localhost, on the default port `4770`. The UI will use the REST endpoint. Make sure to enable CORS on the REST endpoint by adding the following to the `bot.yml` file:

```yaml
# bot.yml
enable_rest_cors: true
```

## Chatto Trivia Pro

The demo called [Chatto Trivia Pro](https://github.com/jaimeteb/chatto_trivia_pro/) contains a very simple trivia game chatbot. You can check out the details about the configuration files at the repository, but it is basically a Chatto instance using an extensions server and a Redis store. The chatbot is deployed on [Google Cloud Functions](https://cloud.google.com/functions) using the [Serverless Framework](https://www.serverless.com/framework/docs). It uses [Open Trivia Database](https://opentdb.com) for the questions and answers.

This demo is also available via [**Telegram**](http://t.me/ChattoTriviaProBot).

Architecture diagram:

![Chatto Trivia Pro](/img/chatto_trivia_pro.svg)

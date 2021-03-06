# Channels

In the **chn.yml** you can insert the credentials for a Telegram Bot, Twilio phone number, and/or Slack App.

```yaml
telegram:
  bot_key: MY_BOT_KEY

twilio:
  account_sid: MY_ACCOUNT_SID
  auth_token: MY_AUTH_TOKEN
  number: MY_NUMBER

slack:
  token: MY_SLACK_TOKEN
  app_token: MY_SLACK_APP_TOKEN
```

!!! note
    You can leave these fields empty and set the corresponding environment variable names with uppercase and underscores, for example:

    ```yaml
    telegram:
      bot_key:    # TELEGRAM_BOT_KEY=this-is-my-telegram-bot-key
    ```

### Telegram

You can connect your Chatto bot to [Telegram](https://core.telegram.org/bots) by providing your Telegram Bot Key, either directly in the **chn.yml** file or by setting the `TELEGRAM_BOT_KEY` environment variable.

You must [set the bot's webhook](https://core.telegram.org/bots/api#setwebhook) to the `/channels/telegram` endpoint in order to receive messages.

<p align="center">
<img src="/img/telegram_channel.jpg" alt="Telegram" width="300"/>
</p>

### Twilio

Similarly, connect your bot to [Twilio](https://www.twilio.com/messaging-api) by adding your credentials to the file or by setting the corresponding environment variables (`TWILIO_ACCOUNT_SID`, etc.).

You must set the webhooks to the `/channels/twilio` endpoint in order to receive messages.

<p align="center">
<img src="/img/twilio_channel.jpg" alt="Twilio" width="300"/>
</p>

### Slack

You can connect your bot to your Slack workspace by adding your [Slack App](https://api.slack.com/apps) Tokens to the **chn.yml** file directly or set the `SLACK_TOKEN` and `SLACK_APP_TOKEN` environment variables.

#### Event Subscriptions

You can use Slack Event Subscriptions to interact with your bot. To receive messages make sure you:

* Enable Event Subscriptions and set the request URL to `/channels/slack`.
* Subscribe to `app_mention` and `message.im` events.

#### Socket Mode

You can also use your bot in Slack's socket mode. To do this:

* Enable [Socket Mode](https://api.slack.com/apis/connections/socket#toggling) on your app.
* Add the generated token to the **chn.yml** as `app_token`, along with the bot token. The file should look like this:

    ```yaml
    slack:
      token: xoxb-my-bot-token
      app_token: xapp-1-my-app-token
    ```

<p align="center">
<img src="/img/slack_channel.jpg" alt="Slack" width="300"/>
</p>
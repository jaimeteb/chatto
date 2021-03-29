# Miscellaneous Bot

[**Miscellaneous Bot**](https://github.com/jaimeteb/chatto/tree/master/examples/02_misc) is a bot that can do different tasks and interact with a lot of external APIs.

The bot can:

* **Tell the weather**: The bit will ask your location and provide a weather forecast for that place, using the [Weather API](https://www.weatherapi.com/).
* **Tell a joke**:  The bot will return a randomly selected joke using the [Joke API](https://jokeapi.dev/).
* **Get a random quote**: Get a random quote with the help of [Quotable](https://github.com/lukePeavey/quotable).
* **Answer random questions**: Try to answer general questions using [Scale SERP](https://www.scaleserp.com/).


This bot demonstrates the integration capabilities of the *Extensions* and a slightly more complex *Classifier*.

## Diagram

This bot's Finite State Machine can be visualized like this:

![Misc](/img/chatto_misc.svg)

## Run it

!!! warning
    If you want to run this example you have to get your own API Keys for [Weather API](https://www.weatherapi.com/) and [Scale SERP](https://www.scaleserp.com/).

To run this example:

```bash
go run examples/02_misc/ext/ext.go
```

And in other terminal:

```bash
export WEATHER_API_KEY=<your weatherapi.com api key>
export SCALE_SERP_API_KEY=<your scaleserp.com api key>
chatto --path examples/02_misc/
```

<!-- ## Test it

You can test a live version of this example on Telegram. Click [**here**](https://t.me/chatto_misc_bot) to use the bot.

<p align="center">
<img src="/img/chatto_misc_telegram.jpg" alt="Misc" width="300"/>
</p> -->

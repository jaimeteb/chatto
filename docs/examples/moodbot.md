# Mood Bot

[**Mood Bot**](https://github.com/jaimeteb/chatto/tree/master/examples/01_moodbot) is a Chatto version of [Rasa's Mood Bot](https://github.com/RasaHQ/rasa/tree/master/examples/moodbot). It aims to demonstrate a basic conversation flow.

## Diagram

This bot's Finite State Machine can be visualized like this:

![Mood Bot](/img/chatto_mood_bot.svg)

!!! note
    The **clf.yml** used in this example is not very good at classifying messages, since it uses a Naïve-Bayes Classifier, instead of [Rasa's NLU Pipeline](https://rasa.com/docs/rasa/tuning-your-model/).

## Run it

To run this example:

```bash
chatto -path examples/01_moodbot/
```

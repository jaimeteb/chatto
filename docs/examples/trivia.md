# Trivia

The [**Trivia Quiz**](https://github.com/jaimeteb/chatto/tree/master/examples/04_trivia) is a simple three-question quiz. At the end of the quiz you'll receive your score. You can exit the trivia at any time.

This example demonstrates how stored [slots](/finitestatemachine/#slots) can be used in the conversation. Also, this example provides a [Python extension server](/extension/#other-languages).

The extensions for this bot are used to:

* Reject the user's input if their answer is not a number from the list.
* Calculate the user's score based on the valid stored answers.

## Diagram

This bot's Finite State Machine can be visualized like this:

![Trivia](/img/chatto_trivia.svg)

## Run it

To run this example:

* Using Go extensions:

    ```bash
    go run examples/04_trivia/ext/ext.go
    ```

* Using Python extensions:

    ```bash
    python3 examples/04_trivia/ext/ext.py
    ```

    !!! note
        This Python example requires Flask to be installed:

        ```bash
        pip install flask
        ```

And in other terminal:

```bash
chatto -path examples/04_trivia/
```

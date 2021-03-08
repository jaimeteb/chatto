# Finite State Machine

The Finite State Machine (FSM) is based on the one shown in [this article](https://levelup.gitconnected.com/implement-a-finite-state-machine-in-golang-f0438b6bc0a8). The states, commands, default messages and transitions are described in the **fsm.yml** file:

```yaml
transitions:
  - from:
      - "off"
    into: "on"
    command: "turn_on"
    answers:
      - text: "Turning on."

  - from:
      - "on"
    into: "off"
    command: "turn_off"
    answers:
      - text: "Turning off."
      - text: "❌"

defaults:
  unknown: "Can't do that."
```

Under `transitions` you can list the transitions available for the FSM. Each object describes the states of the transition (`from` a list of states `into` another state) if `command` is executed; `answers` is the list of messages to send to the user.

In other words:

```yaml
  - from:
      - "off"                # the list of origin states
    into: "on"               # the target state
    command: "turn_on"       # the command that will do the transition
    answers:
      - text: "Turning on."  # the answers that are sent after the transition
```

The commands used in the transitions must correspond with the ones listed in the [classifier](/classifier).

## Answers

Answers are formed by a *text* field and/or an *image* URL. For example:

```yaml
    answers:
      - text: "Single message."

    answers:
      - text: "A list"
      - text: "of simple messages"

    answers:
      - text: "An image will be attached"
      - text: "to this message"
        image: https://i.imgur.com/8MU0IUT.jpeg
```

## *Any*

The special state **any** can help you go from any state into another, if the command is executed.

```yaml
  # The "end" command will move from any state into the
  # initial state, sending the message "Bye bye!"
  - from:
      - any
    into: initial
    command: end
    answers:
      - text: "Bye bye!"
```

Also, the special command **any** is used to transition between states, regardless of the command predicted.

```yaml
  # When in the "search_pokemon" state, any command will do the
  # transition into the initial state
  - from:
      - search_pokemon
    into: initial
    command: any
    extension: search_pokemon
```

## Default Messages

In the **fsm.yml** file, the `defaults` section is used to set the messages that will be returned when the following events happen:

- **`unknown`**: The current state does not transition into another one with the predicted command.
- **`unsure`**: The command prediction confidence was below the threshold.
- **`error`**: An error ocurred during the execution of an extension.

Here's an example of default messages:

```yaml
defaults:
  unknown: "Can't do that transition."
  unsure: "Sorry, I didn't understand that."
  error: "An error ocurred."
```

## Slots

You can save information from the user's input by using `slot` objects. In the **fsm.yml** file, slots are declared as such:

```yaml
  - from:
      - question_1
    into: question_2
    command: any
    extension: val_ans_1
    slot:
      name: answer_1
      mode: whole_text
```

In this example, in the transition from **question_1** to **question_2**, a slot called **answer_1** will be saved, in other words, the user's message is stored in memory.

The supported slot modes are:

* **`whole_text`**: Saves the entire input in the slot.

    ```yaml
    slot:
      name: answer_1
      mode: whole_text
    ```

* **`regex`**: Saves the first Regular Expression match found.

    ```yaml
    slot:
      name: answer_1
      mode: regex
      regex: "[0-9]+"
    ```

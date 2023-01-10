# Changelog

## v0.8.5

* Use `time.Duration` for channel delay.
* Use `time.Duration` for TTL and Purge in FSM Stores.
* Consider "from" and "into" states in StateTable.

---

## v0.8.4

* Add `command` in extension request.
* Add optional delay in channel messages.

---

## v0.8.3

* Fix `chatto init` files.
* Use default error message when an extension fails.

---

## v0.8.2

* Use Cobra for command-line interface.
* Unify `chatto`, `chatto-cli` and `chatto-init` into one single `chatto` binary.
* Add SQL option for FSM store:

	```yaml
	# bot.yml

	store:
	  type: sql
	  rdbms: mysql
	  host: localhost
	  user: root
	  password: root
	  database: chatto
	  ttl: 20
	```

---

## v0.8.1

* Add K-Nearest Neighbors classifier with [fastText](https://fasttext.cc/docs/en/pretrained-vectors.html) sentence-wise average word vectors as features.
	* Truncate word vectors to a certain percentage.
	* Optionally skip out-of-vocabulary words.
* Add Tf-Idf option for Naïve Bayes classifier.
* Load and/or save models from files.

Example `model` object in **clf.yml** file:

```yaml
model:
  classifier: knn
  parameters:
    k: 5
  directory: ./model
  save: true
  load: false
  word_vectors:
    file_name: ./model/wiki.en.vec
    truncate: 0.01
    skip_oov: true
```

---
## v0.8.0

* Multiple extension servers.
	* An alias for each extension server must be specified in **bot.yml**.
	* A server name and an extension name are required when using extensions in **fsm.yml**.

		```yaml
      	# bot.yml
      	extensions:
		  my_rest_server:                 # this server will be referenced as
		    type: REST                    # "my_rest_server" in the fsm.yml file
		    url: http://localhost:8770

		  my_rpc_server:                  # this server will be referenced as
			type: RPC                     # "my_rpc_server" in the fsm.yml file
			host: localhost
			port: 8770
		```

		```yaml
		# fsm.yml
		transitions:
		  - from:
			  - initial
			into: another_state
			command: foo
			extension:                    # this transition will execute the
			  server: my_rest_server      # extension "foo_extension" from the
			  name: foo_extension         # "my_rest_server" extension server

		  - from:
			  - initial
			into: another_state
			command: bar
			extension:                    # this transition will execute the
			  server: my_rpc_server       # extension "bar_extension" from the 
			  name: bar_extension         # "my_bar_server" extension server
		```

---
## v0.7.1

* Automatically detect commands and states from transitions.
* Make environment variables more k8s-friendly:
	* For **bot.yml** variables have the `CHATTO_BOT` prefix.
	* For **chn.yml** variables have the `CHATTO_CHN` prefix.	
* Add example Kubernetes deployment.
* Auto-reload **fsm.yml** and **clf.yml** files on change.

---
## v0.7.0

* Support token auth for the `chatto-cli` REST channel client.
* Rename and reorganize **fsm.yml** to make transitions more intuitive to write:
	* Rename FSM `functions` to `transitions`
	* Rename fsm `message` to `answers`
	* Rename extension `CommandFuncs` to `Extensions`
	* Renane `extension` pakcage to `extensions`
* Make channel handlers public.
* Add `bot.cleanAnswers` to remove empty parameters in response body.

!!! warning
	Package and function names were changed since the last version.

---
## v0.6.2

* Add `bot.Answer` un public `bot` package.
* Add channel name in extension request body.

---
## v0.6.1

* Add `chatto-init` binary.

---
## v0.6.0

* Add SSL options to `extension.ServeREST`:

	```bash
	go run examples/04_trivia/ext/ext.go -ssl-keyfile localhost.key -ssl-certificate localhost.crt
	```

* Add token authorization option for extensions and bot endpoints:

	```bash
	go run examples/04_trivia/ext/ext.go -token my-authorization-token
	```

* Add default messages per conversation.
* Add version endpoint to REST and RPC extensions.
* Add ability to transition from multiple states:

	```yaml
	- transition:
		from:
		  - "state_1"
		  - "state_2"
		into: "state_3"
	```

---
## v0.5.1

* Use goreleaser.
* Provide a public server and client `bot` package.
* Move private APIs into `internal` directory.
* Make `chatto-cli` as a separate binary.
* Add `version` flag.

---
## v0.5.0

* Upgrade to Go 1.15
* Add Socket Mode for the Slack channel.
* Add `extension` field in **fsm.yml** functions to indicate the extension to run.

	!!! note
		Extension names do not need to begin with **"ext_"** any more.

* The `message` field in **fsm.yml** functions can only contain a list of message objects:

	```yaml
		message:
		- text: "Single message."

		message:
		- text: "A list"
		- text: "of simple messages"

		message:
		- text: "An image will be attached"
		- text: "to this message"
			image: https://i.imgur.com/8MU0IUT.jpeg
	```

* Extension functions must return a list of [`query.Answer`](https://godoc.org/github.com/jaimeteb/chatto/query#Answer)s, either by creating the objects or by using the [`query.Answers`](https://godoc.org/github.com/jaimeteb/chatto/query#Answers) function.

	!!! warning
		Package and function names were changed since the last version. Please refer to the [extensions documentation section](https://chatto.jaimeteb.com/extensions) for the updated names.

* Use `runtime.Formatter` when `DEBUG=true`.

---
## v0.4.2

* Add new example [*Misc Bot*](https://github.com/jaimeteb/chatto/tree/master/examples/02_misc).

---
## v0.4.1

* Default log level to **INFO**. Change to **DEBUG** with the environment variable `DEBUG` set to `true`.
* Include *Sender* in `ext.Request` as `Sen`.
* Add *regex* slot mode, which saves the first Regular Expression match found.

    ```yaml
    slot:
      name: answer_1
      mode: regex
      regex: "[0-9]+"
    ```

---
## v0.4.0

* Add image support

	Messages can be either simple strings, or objects formed by an *image* URL and/or *text*.
	
	Example on **fsm.yml**:

	```yaml
	message:
      - text: "Oh don't be sad :("
      	image: https://i.imgur.com/8MU0IUT.jpeg
      - "Did that help?"
	```

	Example on extensions:

	```go
	func greetFunc(req *ext.Request) (res *ext.Response) {
		return &ext.Response{
			FSM: req.FSM,
			Res: cmn.Message{
				Text:  "Hello Universe",
				Image: "https://i.imgur.com/pPdjh6x.jpg",
			},
		}
	}
	```

* Moved *Message* objects into separate **common** package
* Moved *Extensions* into separate **ext** package
* Added Slack channel

	Add the Slack Token in the **chn.yml** file or set the `SLACK_TOKEN` environment variable.

	```yaml
	slack:
  	  token: MY_SLACK_TOKEN
	```

	And set the Event Subscriptions request URL to `/endpoints/slack`

* Refactored channels to allow for easier integration

---
## v0.3.1

* Improve CLI appearance
* Expire FSM in Cache and Redis
    * Use [patrickmn/go-cache](https://github.com/patrickmn/go-cache)
* Support easier integration for REST extensions (see below)

An extension REST service must implement these routes:

- `GET /ext/get_all_funcs`

	This route should return an array with the names of the registered extensions.

	Example response:
	```json
	[
		"ext_val_ans_1",
		"ext_val_ans_2",
		"ext_score"
	]
	```

- `POST /ext/get_func`

	This route should return the resulting Finite State Machine object after the extension's execution, along with a message or messages.

	Example request body:

	```json
	{
		"fsm": {
			"state": 2,
			"slots": {
				"answer_1": "3"
			}
		},
		"req": "ext_val_ans_1",
		"txt": "3",
		"dom": {
			"state_table": {
				"any": -1,
				"initial": 0,
				"question_1": 1,
				"question_2": 2,
				"question_3": 3
			},
			"command_list": [
				"start",
				"end"
			],
			"default_messages": {
				"unknown": "Not sure I understood, try again please.",
				"unsure": "Not sure I understood, try again please.",
				"error": "I'm sorry, there was an error."
			}
		}
	}
	```

	Example response:
	```json
	{
		"fsm": {
			"state": 2,
			"slots": {
				"answer_1": "3"
			}
		},
		"res": "Select one of the options."
	}
	```

---
## v0.3.0

* Add support for **RPC** as well as **REST** extensions
* Added **bot.yml** file to load extensions client-side configurations and store configurations
```yaml
bot_name: "test_bot"
extensions:
  type: REST
  url: http://localhost:8888
store:
  type: REDIS
  host: localhost
  password: pass
```
* Values for **bot.yml** can be loaded from environment variables
* Moved **pl.yml** contents into **clf.yml**
* Log with [logrus](https://github.com/sirupsen/logrus) and set log level with environment variable **LOG_LEVEL**
* Add *error* default message
* Add `-port` flag for chatto and extensions

---
## v0.2.4

* Move channel configurations to a separate file (**chn.yml**)
* Added Twilio channel

```yaml
telegram:
  bot_key:

twilio:
  account_sid:
  auth_token:
  number:
```
Values for **chn.yml** can be loaded from environment variables, for example: ```TELEGRAM_BOT_KEY``` and ```TWILIO_AUTH_TOKEN```

* Added new example (*Pokemon bot*)

---
## v0.2.3

* Add optional **pl.yml** file
* Improve docker compose
* Use *unsure* message when classification is below threshold

---
## v0.2.2

* Add extension host and port environment variables
* Add chatto service port environment variable
* Provide docker and docker-compose examples
* Implement github action for docker

---
## v0.2.1

* Allow multiple messages in fsm file

---
## v0.2.0

* Extensions over RPC
* Telegram support
* Redis support

---
## v0.1.0

* Implemented extensions (Golang plugin build mode)
* Implemented slot saving
* Updated fsm.yml structure


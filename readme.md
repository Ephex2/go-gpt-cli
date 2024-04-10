# Overview

This project is meant to provide a CLI tool to allow you to interact with servers that implement the OpenAI API spec.

It offers the ability to create different profiles for different endpoints, allowing for the ability to quickly swap between different sets of saved parameters for different endpoints. For example, there could be a profile with a system prompt that is optimized for completing requests for code, while another could be setup for a more conversational experience.

This project is still a work in progress - let me know if anything needs immediate improvement!


<br/>

## Quickstart

Build the project, and set up the api key:

``` bash
cd src
go build
./go-gpt-cli config apikey mykey12345
```

Query away!

``` bash
./go-gpt-cli chat prompt "Write a short poem about a sunrise"
```
```
In the stillness of dawn's embrace,
Soft hues of light begin to grace.
Silent whispers in the sky,
Golden rays that gently fly.

A canvas painted fresh and new,
Nature's beauty in full view.
As darkness fades and colors blend,
A sunrise story without end.
```

<br/>

## Change Base Url  

If you want to connect to another API which implements the OpenAI API spec, the URL can be changed quickly using the seturl command:

``` bash
./go-gpt-cli config seturl http://my.local.instance:PORT
```

<br/>

## Profiles and Endpoints

The term used for the routes which offer different functionality (image handling, chat completions, etc.) in this project is 'endpoints'.

Examples of endpoints are:

- audio
- chat
- image
- embeddings
- and more

Each endpoint which can be used to make calls to the API and requires parameters to be set (e.g. model is excluded) will have an associated profile.

These profiles can be used to modify the values sent in a request to the api.

To list profile commands available use:

``` bash
./go-gpt-cli profile -h
```

You can create profiles by specifying an endpoint and name for them, as seen below:

``` bash
./go-gpt-cli profile create chat codereview

# ./go-gpt-cli profile create <endpointName> <profileName>
```

Profiles are represented as json; they be read and updated as required:

``` bash
./go-gpt-cli profile read chat codereview > codeprofile.json
# --- perform editing as desired ---

# update the profile
./go-gpt-cli profile updated chat ./codeprofile.json
```

A profile can be set as the default profile to be used for an endpoint by using the 'default' command:

``` bash
./go-gpt-cli profile default chat codereview
```

A default profile (as specified by values set in the code) will be created and set as the default for an endpoint if no default profile for an endpoint exists when a command using that endpoint is run.

To list existing default profiles, use the ```config get``` command:

``` bash
./go-gpt-cli config get

{
  "ApiKey": "abc123",
  "BaseUrl": "https://api.openai.com",
  "audioDefaultProfile": "default",
  "chatDefaultProfile": "default",
  "embeddingsDefaultProfile": "default",
  "fileDefaultProfile": "default",
  "finetuningDefaultProfile": "default",
  "imageDefaultProfile": "default"
}
```

**Note:** This will also output your apikey in plain-text.

<br/>

## Cobra completions

This project uses standard cobra completions to help autocomplete shell commands. To setup documentation for completions, use ```completion -h```:

``` bash
./go-gpt-cli completion -h

To load completions:

Bash:

  $ source <(go-gpt-cli completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ go-gpt-cli completion bash > /etc/bash_completion.d/go-gpt-cli
  # macOS:
  $ go-gpt-cli completion bash > $(brew --prefix)/etc/bash_completion.d/go-gpt-cli

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ go-gpt-cli completion zsh > "${fpath[1]}/_go-gpt-cli"

  # You will need to start a new shell for this setup to take effect.

fish:

  $ go-gpt-cli completion fish | source

  # To load completions for each session, execute once:
  $ go-gpt-cli completion fish > ~/.config/fish/completions/go-gpt-cli.fish

PowerShell:

  PS> go-gpt-cli completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> go-gpt-cli completion powershell > go-gpt-cli.ps1
  # and source this file from your PowerShell profile.

Usage:
  go-gpt-cli completion [bash|zsh|fish|powershell]

Flags:
  -h, --help   help for completion
```


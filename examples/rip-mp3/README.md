# Bashbot Example - Rip-Mp3 Command

In this example, a youtube url is passed to the rip-mp3 command with the [youtube-dl tool](https://youtube-dl.org/) installed on the bashbot host. An audio file is stripped from the youtube video, and returned as a message with the `bashbot send-file` command. A full example can be found at [rip-mp3.yaml](rip-mp3.yaml)

```text
!bashbot rip-mp3 https://www.youtube.com/watch?v=dQw4w9WgXcQ
```

<img src="https://i.imgur.com/hOPmOli.gif">

## Bashbot Configuration

This command depends on [ffmpeg](https://ffmpeg.org/), which also depends on python3. By default, Bashbot container use a non-root user at runtime, and cannot install additional tools with apk, and in this example, these dependencies are installed using the [latest-root](https://hub.docker.com/r/mathewfleisch/bashbot/tags?page=1&ordering=last_updated&name=latest-root) base-image, from the `dependencies` object (though it would be better to bake these steps into your base image).

```yaml
dependencies:
  - name: "Install youtube-dl tool"
    install:
      - "apk add python3 py3-pip ffmpeg &&"
      - "ln -s /usr/bin/python3 /usr/local/bin/python &&"
      - "wget https://yt-dl.org/downloads/latest/youtube-dl -O /usr/local/bin/youtube-dl &&"
      - "chmod +x /usr/local/bin/youtube-dl"
```

In the values.yaml of the bashbot helm chart, the runtime user can be overridden to root by using the desired tag with `-root` appended (i.e. `latest-root` or `v2.0.5-root`)

```yaml
image:
  repository: mathewfleisch/bashbot
  tag: latest-root
```

The command configuration for rip-mp3 defines the dependencies that must be present on the host, and one parameter that uses a regular expression to match a valid youtube url. That url is passed to the youtube-dl tool, and an mp3 file is then passed back to the channel it was triggered from using the bashbot cli tool.

```yaml
name: Rip mp3 from Youtube
description: Use the youtube-dl tool to download the audio file from a youtube videop
envvars: []
dependencies:
  - python3
  - ffmpeg
  - youtube-dl
help: "!bashbot rip-mp3 [youtube-url]"
trigger: rip-mp3
location: ./
command:
  - "youtube-dl -q -o 'bashbot-youtube-rip-${TRIGGERED_AT}.%(ext)s' -x --audio-format mp3 ${url} &&"
  - "bashbot send-file --channel ${TRIGGERED_CHANNEL_ID} --file ${PWD}/bashbot-youtube-rip-${TRIGGERED_AT}.mp3"
parameters:
  - name: url
    allowed: []
    description: A valid url (Expecting youtube)
    match: ^(https)://([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:/~+#-]*[\w@?^=%&/~+#-])?$
log: true
ephemeral: false
response: text
permissions:
  - all
```

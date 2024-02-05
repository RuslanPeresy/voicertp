# voicertp

Change these lines:

```
var Token = ""     // Set bot token
var ChannelID = "" // Set channel id
var GuildID = ""   // Set server id
```

## How to use

Run bot:

```
go run client.go
```

or build and run:

```
go build
./voicertp
```

If credentials are correct, you will see the bot joined the channel in your Discord.
As soon as any speakers will talk to the same channel, you will be able to listen to RTP packets on rtp://127.0.0.1:8080

```
ffmpeg -protocol_whitelist 'file,rtp,udp' -i example.sdp audio.ogg
```

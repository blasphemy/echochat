echochat
========

[![Build Status](https://travis-ci.org/blasphemy/echochat.svg?branch=master)](https://travis-ci.org/blasphemy/echochat)[![Project Stats](https://www.openhub.net/p/echochat-ircd/widgets/project_thin_badge.gif)](https://www.openhub.net/p/echochat-ircd)


A simple ircd written in Go.

Requirements
============

* Go

No really, that's all. It only depends on the Go standard library. It can run anywhere that go runs. I even got it to handle a few thousand users an an android tablet.

Building
========

This part probably really goes against the go workflow, but oh well. Here's the short version.
```bash
git clone https://github.com/blasphemy/echochat.git
cd echochat/echat
go build
```

That's it. You should also be able to do it the normal way using ```go get```. To run it, just do ```./echat```, or ```echat.exe``` on windows. Configuration explained below.

Configuration
=============
compile and run echochat to make a sample config file. It will be saved as echochat.json.

Here's an explanation

```javascript
{
  "ServerName": "test.net.local",
  "ServerDescription": "A test server", //Server description or note.
  "DefaultKickReason": "Your behavior is not conductive of the desired environment.",
  "DefaultKillReason": "Your behavior is not cunductive of the desired environment.",
  "DefaultQuitReason": "Leaving",
  "DefaultPartReason": "Leaving",
  "PingTime": 45, //Time between pings and disconnects due to ping timeouts.
  "PingCheckTime": 20, //Even though there is a ping time, it is only checked at this invertal
  "ResolveHosts": true, // Hostname resolving (localhost.localdomain vs 127.0.0.1)
  "DefaultCmode": "nt", // Default modes set on a channel upon creation
  "StatTime": 30, // Some stats are dumped to the terminal at this invertal
  "Debug": false, // Debug statements. You probably don't want this unless you're hacking on it
  "Cloaking": false, // Cloak hostnames
  "Salt": "default", // Salt, used for cloaking hostnames, and possibly any other cryptographic operations in the ircd.
  "Privacy": true, //Don't log things like messages, could be considered a violation of privacy.
  "OpersKickable", : false, //Are you able to kick an OPER from a channel?
  "SystemJoinChannels": true, //Joins the "system" user to all channels.
  "Logfile": "echochat.log", //File to write logs to
  "ListenIPs": [ // List of IPs to listen on
    "0.0.0.0"
  ],
  "ListenPorts": [ //List of ports to listen on
    6667,
    6668,
    6669
  ],
  "LogChannels": [ //You can log to channels if you like. By default, these channels will have mode +A, so only opers can join
    "#opers",
    "#log"
  ],
  "Opers": { //List of opers. Takes a plaintext username/password combo.
    "default": "password"
  },
  "AutoJoin": [ //List of channels to join a user to on connect. 
    "#default",
    "#home"
  ]
}
```

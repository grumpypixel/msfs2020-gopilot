# msfs2020-gopilot

## tl;dr
The GoPilot is just another browser-based VFR map for [Microsoft Flight Simulator 2020](https://www.flightsimulator.com) (MSFS2020).

## Features

* VFR Map
    * with maps from OpenStreetMap, Stamen Design, Mapbox and Bing Maps
    * a HUD displaying some vital/basic information
    * a simple map marker
* Teleport Service
    * Pick a point on the map, select the desired altitude, heading and airspeed and off you go

## Releases

Download the latest version [here](https://github.com/grumpypixel/msfs2020-gopilot/releases)

## Screenshot(s)
<img src="https://user-images.githubusercontent.com/28186486/106658243-54170e00-659d-11eb-84e6-24e1bf66447e.png" width="20%"></img>

## How do I build GoPilot?

Assuming you have installed Go on your machine and cloned/downloaded the repo, you can build & run GoPilot as follows:

Bash:
```console
$ build.sh
$ run.sh
```

Or manually (also Bash):
```console
$ CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o gopilot.exe gopilot/main.go gopilot/request_manager.go
$ ./gopilot.exe
```

Windows Command Prompt:
```console
$ build.bat
$ gopilot.exe
```

## GoPilot command-line options

You can run the GoPilot with the following options:

* Connection Name: --name <YOUR_CONNECTION_NAME> (default: "GoPilot")
* DLL Search Path: --searchpath <PATH_TO_SIMCONNECT_DLL> (default: ".")
* Server Address: --address <ADDR:PORT> (default: "0.0.0.0:8888")
* Request Interval: --requestinterval <INTERVAL_IN_MILLISECONDS> (default: 250)
* Timeout: --timeout <TIMEOUT_IN_SECONDS> (default: 600)

Example:
```console
$ gopilot.exe --name POTATOSQUAD --searchpath ../.. --address 0.0.0.13370 --timeout 1000
```

## The GoPilot is running. Now what?

The gopilot.exe starts a local web server which you can connect to with a browser.

Now open the browser of your choice and type the following into the address bar:

```console
http://localhost:8888/
```
Note: This will only work, of course, when the GoPilot and the browser are running on the same machine.

If you want to connect to the web server from another machine, find the IP address of the computer the gopilot.exe is running on (probably something like: 192.168.1.73) and go like this:
```console
http://192.168.1.73:8888/
```

## Web Server API

The following web server routes are available:
* `/` or `/vfrmap` opens the VFR map
* `/teleport` opens the Teleport Service. Be advised not to teleport yourself into the ground mistakenly.
* `/mehmap` opens a plain & simple map without distractions. (So no HUD. No nothing. Meh.)
* `/setdata` opens an *experimental* page where you can manually set data on the *sim object*. DO NOT USE THIS if you don't know what you're doing. This might (and probably will) CRASH your simulator. Seriously.
* `/simvars` display all registered simulation variables (no auto-update)
* `/debug` display debug information (also no auto-update)

Examples:
* `http://localhost:8888/vfrmap` or simply: `http://localhost:8888`
* `http://localhost:8888/teleport`
* `http://localhost:8888/mehmap`
* `http://localhost:8888/setdata`

### VFR Map Options

The VFR map comes with a bunch of options (*URL Parameters*) which can be specified in the address bar.

The format is as follows:\
`http://localhost:8888?<PARAM_NAME>=<PARAM_VALUE>`

Or a combination of multiple options:\
`http://localhost:8888?<PARAM_NAME_1>=<VALUE_1>&<PARAM_NAME_2>=<VALUE_2>`

List of available options:
* `dms_coords=true|false` display latitude/longitude coordinates in DMS format (*Degrees Minutes Seconds*) (default: true)
* `hud=<true|false>` show/hide the HUD (default: true)
* `layer_control=true|false` show/hide the layer control (default: true)
* `scale_control=true|false` show/hide the scale control (default: true)
* `strip_title=true|false` remove "Asobo" from the plane's type name (default: false)
* `units=<true|false>` show/hide the units in the HUD (default: true)
* `watermark=true|false` display a watermark in the right bottom corner (default: false)
* `watermark_size=<value>` specify the watermark's size (default: 128px). *128px*/*128em*/*128* will result in different sizes.
* `watermark_url=<url>` specify the url of the watermark (default: load a placeholder)
* `watermark_position=bottomright|bottomleft|topleft|topright` set the watermark's position (default: bottomright)
* `zoom_control=true|false` show/hide the zoom control (default: true)
* `plane_size=<number>` specifiy the size of the plane (default: 64)
* `plane_opacity=<decimal>` specify the plane's opacity as a decimal value (default: 1.0)
* `plane_style=black|gray|green|white` set the plane's color (default: white)
* `open_in=<bing|google>` open a marked spot in either Google Maps or Bing Maps (default: bing)
* `marker_event=click|dblclick|contextmenu` specify the mouse event with which the map marker is placed (default: click)
* `mapbox_token=<token>` use your own Mapbox token since the one provided is limited
* `bing_key=<key>` usw your own Bing Maps key since the one provided is also limited

Note: Boolean parameters can be entered as true or false, 1 or 0, yes or no.

Example:
```console
http://localhost:8888?plane_style=green&plane_size=128&plane_opacity=0.73&open_in=bing&dms_coords=0&watermark=1&watermark_url=https://media.giphy.com/media/SgwPtMD47PV04/giphy.gif&marker_event=contextmenu&zoom_control=false
```

## How do I find my IP address?

Look here for help: [Microsoft Support](https://support.microsoft.com/en-us/windows/find-your-ip-address-f21a9bbc-c582-55cd-35e0-73431160a1b9)

Or open the command prompt and enter:
```console
$ ipconfig
```

There you are looking for the something like this:
```console
Wireless LAN adapter Wi-Fi:

   Connection-specific DNS Suffix  . : fritz.box
   IPv4 Address. . . . . . . . . . . : 192.168.1.73
   Subnet Mask . . . . . . . . . . . : 255.255.255.0
   Default Gateway . . . . . . . . . : 192.168.1.1
```
In this case, *192.168.1.73* would be your local IP address which you can use to connect to from another computer.

## Where's this SimConnect.DLL?

The SimConnect.dll comes bundled with the GoPilot executable and will be extracted automatically if the DLL cannot be found in the given search paths. Please be aware that this bundled versoin may no be the latest version of SimConnect.

If you want/need the latest version of the SimConnect.dll, you can find it in the MSFS2020 SDK within the following directory:\
/MSFS SDK/Samples/SimvarWatcher/bin/x64/Release/

Please have a look at the [Flight Simulator Forums](https://forums.flightsimulator.com/t/how-to-getting-started-with-the-sdk-dev-mode/123241) for further instructions if you want to know how to get and install the SDK.

## My virus scanner thinks my Go distribution or compiled binary is infected. What the heck?!

From the official Golang website https://golang.org/doc/faq#virus:

"This is a common occurrence, especially on Windows machines, and is almost always a false positive. Commercial virus scanning programs are often confused by the structure of Go binaries, which they don't see as often as those compiled from other languages.

If you've just installed the Go distribution and the system reports it is infected, that's certainly a mistake. To be really thorough, you can verify the download by comparing the checksum with those on the downloads page.

In any case, if you believe the report is in error, please report a bug to the supplier of your virus scanner. Maybe in time virus scanners can learn to understand Go programs."

## Are there any bugs?

Well. Every program has bugs and I'm pretty sure this one does too. Bugs are annoying and I will fix them as soon as they come up.

## I like this but I would like it more if...

Just send me an e-mail and tell me about it.

## What an outrage! This is all bullcrap! I absolutely dislike it!

Same as above. Just send me an e-mail and let me know why. Or, alternatively, clone this repo, improve on it and build something fresh and astonishing. This way we can all learn from you.

## Motivation

There are some really fine VFR maps already available on GitHub. All of them are fabulous and assist well when exploring the virtual flight simulator world. So why build another one?

- I was always missing something when using other VFR maps. Of course, I could've just cloned another project and build "my stuff" on top of that but then I probably wouldn't have learned much, which brings me to:
- I wanted to delve into Golang more intensively.
- I had the perfect excuse to start up the simulator and shout "CLEAR PROP!" out of the window because I needed to test and evaluate newly implemented features.
- And last but not least and most importantly: Fun-coding. "Go will make you love programming again", they said. "We [friggin'] promise", they said. And yes, they were right and kept their promise. Visit [@golang](https://twitter.com/golang)

## Shoutouts

A big thank you goes out to Microsoft and especially to all the peepz at [Asobo](https://www.asobostudio.com) for creating such an awesome simulator/program/game/dillydaller.

Big, hearty thank yous go out to [Pilot Emilie](https://www.youtube.com/c/Pilotemilie/about) and [Squirrel](https://www.youtube.com/squirrel/about) for their fun and educational videos.
And thanks to the whole, supporting Potato Gang.

Cheers,\
coffeecookie#884@EDDL

CLEAR PROP!!!

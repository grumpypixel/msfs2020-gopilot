# msfs2020-gopilot

## tl;dr
GoPilot is just another browser-based [VFR Map](https://en.wikipedia.org/wiki/Visual_flight_rules) for [Microsoft Flight Simulator 2020](https://www.flightsimulator.com) (MSFS2020).

## Main Features

* VFR Map
    * With maps and overlays from [OpenStreetMap](https://www.openstreetmap.org), [Carto](https://carto.com), [Stamen Design](https://stamen.com), [OpenTopoMap](https://opentopomap.org), [Mapbox](https://www.mapbox.com/maps), [Bing Maps](https://www.bing.com/maps), [openAIP](https://www.openaip.net) and [ESRI](https://de.wikipedia.org/wiki/ESRI)
    * a HUD displaying vital information as heading, airspeed, altitude, vspeed, elevator/rudder trim, outside air temperature, wind velocity...
    * Your latitude/longitude position
    * a humble map marker
* Teleport Service
    * Pick a point on the map, select the desired altitude, heading and airspeed and off you go! *wheee!*
* Airport Finder
    * Locate airports anywhere on the map within a specific radius using data from [OurAirports](https://ourairports.com)

## Releases

Download the latest version [here](https://github.com/grumpypixel/msfs2020-gopilot/releases).

Unzip the archive, run gopilot.exe and browse to `http://localhost:8888` (or: `http://whatever-your-ip-address-may-be:8888`).

## Screenshot(s)
<img src="https://user-images.githubusercontent.com/28186486/107130271-09163700-68cc-11eb-8273-60000598fd9d.png" width="18%"></img>
<img src="https://user-images.githubusercontent.com/28186486/107130274-0b789100-68cc-11eb-80ed-815dccc20fff.png" width="18%"></img>
<img src="https://user-images.githubusercontent.com/28186486/107130276-0fa4ae80-68cc-11eb-8a6e-4b9d2638fe5e.png" width="18%"></img>
<img src="https://user-images.githubusercontent.com/28186486/107130277-129f9f00-68cc-11eb-9120-14a9ecce852d.png" width="18%"></img>
<img src="https://user-images.githubusercontent.com/28186486/108609536-d2154a80-73ce-11eb-935d-498923f8b538.png" width="18%"></img>

## Run GoPilot

Run GoPilot:

```console
$ gopilot.exe
```

Or explicitly specify the config file:

```console
$ gopilot.exe --cfg configs/config.yml
```

Or run with your own configuration file:

```console
$ gopilot.exe --cfg configs/my-config-file.yml
```

## GoPilot is running. Now what?

The gopilot executable starts a local web server which you can connect to with a browser.

VFR Map:\
`http://localhost:8888`\
or\
`http://whatever-your-ip-address-may-be:8888`

Teleport Service:\
`http://localhost:8888/teleport`

Airport Finder:\
`http://localhost:8888/airports`

## Web Server API

The following routes are available:
* `/vfrmap` or `/` opens the VFR map
* `/airports` opens the Airport Finder
* `/teleport` opens the Teleport Service. Be advised not to teleport yourself into the ground mistakenly.
* `/mehmap` opens a plain & simple map without distractions. (No HUD. No nothing. Meh.)
* `/setdata` opens an *experimental* and *hideous* page where you can manually set data on the *sim object*. DO NOT USE THIS if you don't know what you're doing. This might (and probably will) CRASH your simulator. Seriously.
* `/simvars` displays all registered simulation variables (no auto-update)
* `/debug` displays debug information (also no auto-update)

Examples:
* `http://localhost:8888/vfrmap` or simply: `http://localhost:8888`
* `http://localhost:8888/airports`
* `http://localhost:8888/teleport`
* `http://localhost:8888/mehmap`
* `http://localhost:8888/setdata`
* `http://localhost:8888/simvars`
* `http://localhost:8888/debug`

## VFR Map Options

The VFR map comes with a bunch of options (*URL Parameters*) which can be specified in the address bar.

The format is as follows:\
`http://localhost:8888?<PARAM_NAME>=<PARAM_VALUE>`

Or a combination of multiple options:\
`http://localhost:8888?<PARAM_NAME_1>=<VALUE_1>&<PARAM_NAME_2>=<VALUE_2>`

List of available options:

* `default_map=<map-name>` - set the default map (default: OpenStreetMap)
* `layer_control=true|false` - show/hide the layer control (default: true)
* `scale_control=true|false` - show/hide the scale control (default: true)
* `zoom_control=true|false` - show/hide the zoom control (default: true)
* `dms_coords=true|false` - display latitude/longitude coordinates in DMS format (*Degrees Minutes Seconds*) (default: true)
* `hud=<true|false>` - show/hide the HUD (default: true)
* `strip_title=true|false` - remove "Asobo" from the plane's type name (default: false)
* `units=<true|false>` - show/hide the units in the HUD (default: true)
* `watermark=true|false` - display a watermark in the right bottom corner (default: false)
* `watermark_size=<value>` - specify the watermark's size (default: 128px). *128px*/*128em*/*128* will result in different sizes.
* `watermark_url=<url>` - specify the url of the watermark (default: load a placeholder)
* `watermark_position=bottomright|bottomleft|topleft|topright` - set the watermark's position (default: bottomright)
* `attributions=true|false` - show/hide the attributions (default: true)
* `position_overlay=true|false` - show/hide the latitude-longitude position overlay (default: true)
* `plane_overlay=true|false` - show/hide the follow-plane/center-on-plane overlay (default: true)
* `plane_size=<number>` - specifiy the size of the plane (default: 64)
* `plane_opacity=<decimal>` - specify the plane's opacity as a decimal value (default: 1.0)
* `plane_style=black|gray|green|white` - set the plane's color (default: white)
* `open_in=<bing|google>` - open a marked spot in either Google Maps or Bing Maps (default: bing)
* `marker_event=click|dblclick|contextmenu` - specify the mouse event with which the map marker is placed (default: click)
* `mapbox_token=<token>` - use your own token for [Mapbox](https://docs.mapbox.com/help/tutorials/get-started-tokens-api/) since the one provided is limited
* `bing_key=<key>` - use your own key for [Bing Maps](https://docs.microsoft.com/en-us/bingmaps/getting-started/bing-maps-dev-center-help/getting-a-bing-maps-key) since the one provided is also limited
* `bing_culture=<culture-code>` - set the culture code for Bing Maps (default: en-US). See [Bing Maps Culture Codes](https://docs.microsoft.com/en-us/bingmaps/rest-services/common-parameters-and-types/supported-culture-codes)

Note: Boolean parameters can be entered as *true* or *false*, *1* or *0*, *yes* or *no*, *yay* or *nay*, *yeah* or *nope*.

### Examples for VFR Map Options

Random example:
```console
http://localhost:8888/?plane_style=green&plane_size=128&plane_opacity=0.73&open_in=bing&dms_coords=0&watermark=1&watermark_url=https://media.giphy.com/media/SgwPtMD47PV04/giphy.gif&watermark_size=256px&marker_event=contextmenu&layer_control=0&scale_control=0&zoom_control=false&attributions=0
```

Or go zen and switch off all overlays:
```console
http://localhost:8888/?hud=0&plane_overlay=0&position_overlay=0&scale_control=0&layer_control=0&zoom_control=0&attributions=0
```

Set a default map:
```console
# Show CartoDB Positron per default:
http://localhost:8888/?default_map=CartoDB.Positron

# Show Bing Maps Aerial per default:
http://localhost:8888/?default_map=BingMaps.Aerial
```

Default map options:
* BingMaps.Light
* BingMaps.Dark
* BingMaps.Gray
* BingMaps.Aerial
* CartoDB.DarkMatter
* CartoDB.Positron
* CartoDB.Voyager
* Esri.WorldGrayCanvas
* Esri.WorldImagery
* Esri.WorldStreetMap
* Esri.WorldTopoMap
* Mapbox
* OpenStreetMap
* OpenStreetMap.DE
* OpenStreetMap.FR
* OpenStreeMap.Gray
* OpenStreetMap.HOT
* OpenTopoMap
* Stamen.Terrain
* Stamen.Toner
* Stamen.TonerLight

Bing Maps examples:
```console
# Display simplified Chinese Bing Maps
http://localhost:8888/?bing_culture=zh-Hans

# Display Spanish Bing Maps
http://localhost:8888/?bing_culture=es-ES

# Specify your own Bing Maps Key and display German maps
http://localhost:8888/?bing_key=<YOUR_BING_KEY>&bing_culture=de-de

# Bing-Maps-Light/French Combo
http://localhost:8888/?default_map=BingMaps.Light&bing_culture=fr-FR
```

You should probably get your own key for Bing Maps [here](https://www.microsoft.com/en-us/maps/create-a-bing-maps-key).

## VFR Map Keyboard Shortcuts

The available keyboard shortcuts for the VFR map are:

* C - Center on Plane
* F - Follow Plane
* T - Toggle Fullscreen

## How do I build GoPilot myself?

Assuming you have installed [Go](https://golang.org/dl/) on your machine and cloned or downloaded the repository, you can build & run GoPilot as follows:

With [bash](https://gitforwindows.org):

```console
$ ./scripts/build.sh
```

Or manually (with bash):

```console
$ CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o gopilot.exe ./cmd/gopilot/main.go
```

Without bash:

```console
$ ./scripts/build.bat
```

Or manually (without bash):

```console
set CGO_ENABLED=1
set GOOS=windows
set GOARCH=amd64
go build -o gopilot.exe ./cmd/gopilot/main.go
```

## How do I find my IP address?

Look here for help: [Microsoft Support](https://support.microsoft.com/en-us/windows/find-your-ip-address-f21a9bbc-c582-55cd-35e0-73431160a1b9)

Or open the command prompt and enter:
```console
$ ipconfig
```

ipconfig will display all current TCP/IP network configuration values.\
There you are looking for something like this:

```console
Wireless LAN adapter Wi-Fi:

   Connection-specific DNS Suffix  . : wicked.wifi.box
   IPv4 Address. . . . . . . . . . . : 192.168.11.73
   Subnet Mask . . . . . . . . . . . : 255.255.255.0
   Default Gateway . . . . . . . . . : 192.168.11.1
```

In this case, *192.168.11.73* would be your local IP address which you can use to connect to from another computer.

So instead of `http://localhost:8888` you would enter `http://192.168.11.73:8888` in your browser's address bar.

## Where's this SimConnect.dll?

The SimConnect.dll comes bundled with the GoPilot executable and will be extracted automatically if the DLL cannot be found in the given search paths. Please be aware that this bundled version may not be the latest version of SimConnect.

If you want/need the latest version of the SimConnect.dll, install the MSFS2020 SDK. Please consult the [Flight Simulator Forums](https://forums.flightsimulator.com/t/how-to-getting-started-with-the-sdk-dev-mode/123241) for further instructions if you want to know how to download and install the SDK.

By default, the DLL will be located within the directory:\
/MSFS SDK/Samples/SimvarWatcher/bin/x64/Release/

## My virus scanner thinks my Go distribution or compiled binary is infected. What the heck?!

From the official Golang website https://golang.org/doc/faq#virus:

"This is a common occurrence, especially on Windows machines, and is almost always a false positive. Commercial virus scanning programs are often confused by the structure of Go binaries, which they don't see as often as those compiled from other languages.

If you've just installed the Go distribution and the system reports it is infected, that's certainly a mistake. To be really thorough, you can verify the download by comparing the checksum with those on the downloads page.

In any case, if you believe the report is in error, please report a bug to the supplier of your virus scanner. Maybe in time virus scanners can learn to understand Go programs."

## Are there any bugs I should be aware of?

Well. Every program has bugs and I'm pretty sure this one does too. Bugs are annoying and they will fixed with best effort. Report bugs [here](https://github.com/grumpypixel/msfs2020-gopilot/issues).

## I like this but I would like it even more if...

Nice. [Tell](https://github.com/grumpypixel/msfs2020-gopilot/issues) me about it.

## This is all bullcrap! I absolutely dislike it!

Same as above. Open a new [issue](https://github.com/grumpypixel/msfs2020-gopilot/issues) and let me know about your feels. Or, alternatively, clone this repo, improve on it and build something fresh and astonishing. This way we can all learn from your creation.

## Help! Please!

In case you are experiencing any sort of technical issues with the program or if you have awe-inspiring ideas for incredible improvements, please open an [issue](https://github.com/grumpypixel/msfs2020-gopilot/issues).

## Motivation

There are some really fine VFR maps already available on GitHub like [this](https://github.com/lian/msfs2020-go) and [this](https://github.com/hankhank10/MSFS2020-cockpit-companion). All of them are fabulous and assist well when exploring the virtual flight simulator world. So why build another one?

- I was always missing something when using other VFR maps. Of course, I could've just cloned another project and build "my stuff" on top of that but then I probably wouldn't have learned much, which brings me to:
- I wanted to delve into Golang more intensively.
- I had the perfect excuse to start up the simulator and shout "CLEAR PROP!" out of the window because I needed to test and evaluate newly implemented features.
- And last but not least and most importantly: Fun-coding. ["Go will make you love programming again"](https://twitter.com/golang), they said. "We [friggin'] promise", they said. And yes, they were right and kept their promise.

## OurAirports Database

The underlying airport data was downloaded from [OurAirports.com](https://ourairports.com/data/).

### Terms of use for the data
From OurAirports:

"DOWNLOAD AND USE AT YOUR OWN RISK! We hereby release all of these files into the [Public Domain](https://en.wikipedia.org/wiki/Public_domain), with no warranty of any kind — By downloading any of these files, you agree that OurAirports.com, Megginson Technologies Ltd., and anyone involved with the web site or company hold no liability for anything that happens when you use the data, including (but not limited to) computer damage, lost revenue, flying into cliffs, or a general feeling of drowsiness that persists more than two days.

Do you agree with the above conditions? If so, then download away!"

## Shoutouts

A big thank-you goes out to Microsoft and especially to all the peeps at [Asobo](https://www.asobostudio.com) for creating such an awesome simulator/program/game/dillydaller.

Big, hearty thank-yous go out to [Pilot Emilie](https://www.youtube.com/c/Pilotemilie/about) and [Squirrel](https://www.youtube.com/squirrel/about) for their fun and educational videos. And thanks to the whole, supporting Potato Gang.

Big thanks go out to [David Megginson](http://ourairports.com/about.html#credits) and all contributors of OurAirports.com. What an impressive piece of collected awesomeness!

Cheers,\
grumpypixel/coffeecookie#884@EDDL

CLEAR PROP!!!

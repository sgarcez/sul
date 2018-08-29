# Sul - Strava Uploader

Simple Strava activity uploader inspired by [pi-python-garmin-strava](https://github.com/thegingerbloke/pi-python-garmin-strava).

It creates Strava activities by uploading raw files from your device (only `.fit` files supported currently).

## Features

- Starts a local server to guide the process of obtaining an access token with write permissions.
- Uploads all files in a directory concurrently.

## Todo

- Add support for more file types(`TCX`, `GPX`, etc)
- Add tombstones for processed files
- Wrap errors properly in `pkg/uploader`
- Ensure singleton server in `pkg/auth`

## Installation

### Binaries

Download your preferred flavour from the [releases page](https://github.com/sgarcez/sul/releases) and install manually.

For example, for raspberry pi (find the correct [ARM architecture version](https://en.wikipedia.org/wiki/Raspberry_Pi#Specifications)):

```shell
$ curl -L https://github.com/sgarcez/sul/releases/download/v0.0.2/sul_0.0.2_Linux_armv7.tar.gz | tar xz
```

### Using go get

```shell
$ go get -u github.com/sgarcez/sul
```

## Usage

### Setup: obtaining an access token with write permissions

Once you create a Strava Application you can retrieve its id and secret here: [https://www.strava.com/settings/api](https://www.strava.com/settings/api)

```
$ sul auth -i <app-id> -s <app-secret>
```

This will provide you with a URL to visit which will redirect you to a local server once the application has been authorised. The local server will capture and print out the new access token and exit the command.

### Uploading files

You can now use the token obtained in the previous step to upload activities to Strava, for example all files in a mounted Garmin device:

```
$ sul upload -t <token> -d /Volumes/GARMIN/ACTIVITY/
```

## Automatic uploads on device mount

### systemd

`/etc/systemd/system/garmin-sul.service`

```
[Unit]
Description=Garmin Sul trigger
Requires=media-usb0.mount
After=media-media-usb0.mount

[Service]
ExecStart=/opt/sul/run-sul.sh

[Install]
WantedBy=media-usb0.mount
```

You can find the device unit with:

```
$ sudo systemctl list-units -t mount
```

`/opt/sul/run-sul.sh`

```
#!/bin/bash

/opt/sul/sul upload \
    -t <token> \
    -d /media/usb/GARMIN/ACTIVITY/ >> /opt/sul/log 2>&1
```

Enable the Service:

```
$ sudo systemctl start garmin-sul.service
$ sudo systemctl enable garmin-sul.service
```

### udev

See the documentation at [pi-python-garmin-strava](https://github.com/thegingerbloke/pi-python-garmin-strava)

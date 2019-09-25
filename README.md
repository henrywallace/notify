# notify

Send notifications with services like gmail.

To use this with Gmail do:

- Create a new Google API credentials on GCP with oauth "send only" credentials
  [1], and download the credentials file to somewhere like `$HOME/.secrets`.
  Name the credentials `google-credentials.json`.

- Setup with `NOTIFY_FROM=from@gmail.com notify --setup`. Use the `--qr` flag
  if you wish to print a QR code of the setup URL if you do not readily have
  browser access, nor wish to labor over copying the long URL.

- Once setup, send an email like:
```sh
NOTIFY_SECRETS=$HOME/.secrets \
NOTIFY_FROM=from@gmail.com \
NOTIFY_TO=to@yahoo.com \
    notify --type gmail --subject "Hello friend :)"
```

For an example of how to use this in a systemd service that checks for system
updates on arch, see `examples/pkgupdate`.


## References
- [1] https://console.developers.google.com/apis/credentials

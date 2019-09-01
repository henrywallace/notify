# pkgupdate

Deploy systemctl user unit that will check for new package updates, and send an
email notification listing those available updates. By default checking happens
every 30m, and will retry for a few times before giving up within that 30m
window. A cache is kept of the previous state of available updates if any, and
an email will only be sent if the cached state has changed.

## Setup

What follows is the way I like to setup it up for myself. Use this as an
prototype on which to fashion it how you like.

In particular, I like to have an email account for each of my machines, so that
notifications are single historied. If it were ever exposed, the intruder could
only get the notification history for one of the machines.

Here are the steps:
- Obtain a new machine with SSH access at HOSTNAME.
- Create a new API credentials on GCP with oauth "send only" credentials.
- Create a new gmail acount with an address corresponding to HOSTNAME. You
  can avoid providing a number if you create one through your phone's browser.
  Turn on MFA.
- Build a new version of github.com/henrywallace/homelab/go/notify for
  HOSTNAME. If compiling on ARM such as with a rasberry pi, see [1]. For a more
  general list see [2].
- Move that binary to your machine: `ssh HOSTNAME mkdir -p '~/bin'` and then
  `scp notify HOSTNAME:'~/bin/'`.
- Setup `notify` on the machine like `NOTIFY_SECRETS=$HOME/.secrets
  NOTIFY_FROM=hostmachine3000@gmail.com notify --setup`.
- Fashion a copy of github.com/henyrwallace/dotfiles/blob/master/bin/upd.
- Set the unit files for systemd [4]:
```sh
ssh HOSTNAME mkdir -p '~/.config/systemd/user'
scp up up.timer HOSTNAME:'~/.config/systemd/user'
ssh HOSTNAME mkdir -p '~/bin'
scp upd HOSTNAME:'~/bin'
```
- Setup the environment variables for the unit [5] (note too [3] which is more
  broad):
```sh
path="~/.config/systemd.d/user/up.service.d/override.conf"
ssh HOSTNAME mkdir -p '~/.config/systemd.d/user/up.service.d/'.
scp <(echo """
[Service]
ExecStart=$HOME/bin/upd
Environment="NOTIFY_SECRETS=$HOME/.secrets"
Environment="NOTIFY_TO=machine.notifications@mycutedomain.pizza"
Environment="NOTIFY_FROM=hostmachine3000@gmail.com"
Environment="NOTIFY_BIN=$HOME/bin/notify"
""") HOSTNAME:'~/.config/systemd.d/user/up.service.d/override.conf'
```
- Enable and start:
```sh
systemctl -H HOSTNAME --user enable up.timer
systemctl -H HOSTNAME --user enable up --now
```
- Watch how it's doing `journalctl --user -f -u up`.
- If you want to fiddle with the unit files or anything else do
```sh
systemctl --user daemon-reload
systemctl --user restart up up.timer
journalctl --user -f -u up
```

## References
- [1] https://github.com/golang/go/wiki/GoArm
- [2] https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63
- [3] https://www.freedesktop.org/software/systemd/man/environment.d.html
- [4] https://wiki.archlinux.org/index.php/systemd/User#How_it_works
- [5] https://serverfault.com/a/413408/430816

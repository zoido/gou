This was a short-lived hobby project. In the end I managed to achieve the same functionality
via [televison][tv] and following channel.

```toml
[metadata]
name = "zellij-xurls"
description = "Select URLs from zelij dump"
requirements = ["xurls"]

[source]
command = [
  { name = "strict", run = "xurls < /tmp/zellij-screen-dump | sort -u" },
  { name = "lax", run = "xurls -r < /tmp/zellij-screen-dump | sort -u" },
]

[keybindings]
ctrl-o = "actions:open"
ctrl-y = "copy_entry_to_clipboard"
enter = "actions:open-and-exit"

[actions.open]
description = "Open URL"
command = "xdg-open '{}'"
mode = "fork"

[actions.open-and-exit]
description = "Open URL and turn off television"
command = "xdg-open '{}'"
mode = "execute"
```

[tv]: https://alexpasmantier.github.io/television/

---

# Go\[U\]

Go to URL. Small program that allows opening URLs in terminal without touching
the mouse.
Designed to be used along terminal multiplexers.

Scans input for URLs and displays a list you use to select URL to
open or to copy to your clipboard.

An alternative to [urlscan][1] or [urlview][2].

[1]: https://github.com/firecat53/urlscan
[2]: https://github.com/sigpipe/urlview

## Shoulders of Giants

Combines the prior available work.

- [xurls](mvdan.cc/xurls/v2)
- [charm.land](https://charm.land/)
- [aymanbagabas/go-osc52](https://github.com/aymanbagabas/go-osc52/v2 )
- [github.com/rkoesters/xdg](https://github.com/rkoesters/xdg)

## Keys

- `q`: quit
- `y`: copy to clipboard via OSC 52 escape sequence.
- `o`: open URL with [XDG](https://en.wikipedia.org/wiki/Freedesktop.org)

More will be available when customisation is implemented.

# Untear worlds torn apart by PaperMC

You can use untear to revert Paper worlds to vanilla worlds.

Untear will look for 3 directories with the same prefix (`world` by default).
- `{prefix}`
- `{prefix}_nether`
- `{prefix}_the_end`

If all are present, it will create a new directory `vanilla_{prefix}`, and merge the 3 worlds into it.

It is non-destructive to the original Paper worlds.

I have tested that it preserves mobs, but..<br>
⚠️ **Remember to check all dimensions vanilla version before permanently deleting Paper worlds!**

## Installation
1. Go
```bash
$ go install github.com/emilyxfox/untear@latest
```

## Usage
```bash
$ ls
# world world_nether world_the_end

$ untear
# ...
# INFO 🎉 Success!
# INFO 🟩 Worlds have been merged

$ ls
# vanilla_world world world_nether world_the_end
```
## Help
```bash
$ untear --help
# Untear lets you rejoin the three Minecraft dimensions into a single world file so it usable with vanilla servers and in singleplayer.
#
# Usage:
#   untear [path] [flags]
#
# Examples:
# untear server-files/ --prefix world
#
# Flags:
#   -h, --help            help for untear
#   -p, --prefix world    world prefix (default is world) (default "world")
#   -v, --verbose false   enable debug logging (default is false)
```
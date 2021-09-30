# Lardoon

Lardoon is a web repository that provides a clean interface for listing, searching, and downloading ACMI files. Lardoon is intended to be used alongside [TacView](https://www.tacview.net/product/about/en/) and [jambon](https://github.com/b1naryth1ef/jambon) to automatically record and import server-side TacView recordings.

## Web UI

Lardoon comes with a web UI that displays a list of all imported replays, and allows you to view individual flight tracks within those replays. Additionally Lardoon can quickly trim TacView files to only the specific portion related to a given flight, reducing ACMI file/download size greatly.

![Example](https://i.imgur.com/GP21sED.png)

## Example Setup

Lets explore an example lightweight setup that allows you to automatically record and import ACMIs remotely via the TacView realtime protocol.

### Recording

We can use [jambon](https://github.com/b1naryth1ef/jambon) to automatically record TacViews remotely, simply by running a process somewhere like so:

```
$ while :; do jambon record --server coldwar.hoggitworld.com --output saw-$(date +"%m-%d-%y-%H-%M").acmi; sleep 30; done 
```

This doesn't account for some hiccups like network or server outages, but will give you a good starting point.

### Importing

Lardoon requires ACMI files to be processed to extract metadata. This can be done iteratively on files while they are being recorded, so the following is enough:

```
$ while :; do lardoon import -p /mnt/bigdata/lardoon/; sleep 300; done
```

### Pruning

If you decide to delete files from disk you will need to run a prune:

```
lardoon prune --no-dry-run
```
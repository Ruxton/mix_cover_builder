Mix Cover Builder
====================================

Overview
--------

A go CLI application that builds covers from Virtual DJ tracklists.

It parses a sub-section of Virtual DJ tracklist.txt (split with
https://github.com/Ruxton/virtual_dj_tracklist_split) into a 9x9 image with an
optional overlay.

Images are found in order using:
 * iTunes
 * Google images searches for bandcamp with some checks

Binaries
---------

Compiled [binaries are available on github](https://github.com/Ruxton/mixcloud_uploader/releases)

Source Requirements
------------

* GoLang > 1.4.2
* The Internet

Using From Source
--------------------------------

  1. Run bin/build
  1. Run the built packages as below from pkg/

Using Pre-compiled Packages
---------------------------

Using compiled packages:

  `buildcover --tracklist $tracklist --output $outputfile --overlay $overlayfile`
  OR
  `buildcover $tracklist $outputfile --overlay $overlayfile`

Meta
----

* Code: `git clone git://github.com/ruxton/mix_cover_builder.git`

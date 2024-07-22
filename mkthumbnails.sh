#!/bin/sh

set -e

ind=input/
force="$1"

for x in . $ind/*.jpg; do
	#  Safety (literal "input/*.jpg")
	if [ ! -f "$x" ]; then
		continue
	fi

	# Don't create thumbnails for thumbnails (-:
	if echo $x | grep -q '[0-9]\+w.jpg$'; then
		continue
	fi

	y=`basename "$x" .jpg`
	for w in 320 480 640 800 1000; do
		z=$(echo $ind/$y-${w}w.jpg)
		if [ -n "$force" ] || [ ! -f "$z" ]; then
			magick "$x" -thumbnail "${w}x${w}>" -unsharp 0x.5 -auto-orient "$z"
		fi
	done
done

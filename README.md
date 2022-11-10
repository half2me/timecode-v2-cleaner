# timecode-v2-cleaner
Linter and cleaner for the timecode v2 format

This will remove any garbage from the timecode file, as well as make sure the timecode is monotone increasing (required by [mp4fpmsod](https://github.com/nu774/mp4fpsmod)). It fixes timecode by taking bad timestamp values and increasing them by a little* compared to the previous timestamp. The precision is 4 decimal places by default, but this can be adjusted with the `--precision n` flag.

Usage:
```
timecode-v2-cleaner --out fixed.txt original_timecode_v2_file.txt
```

Example timecode files:
https://github.com/nmaier/mkvtoolnix/blob/master/examples/example-timecodes-v2.txt

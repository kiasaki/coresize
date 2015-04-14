# coresize

_core, because it's simple, resize, because it serves resized versions of your images_

**coresize** aims to be a small server that can load images from local disk or S3 and serve resized and aligned versions of those images "on the fly". Here's what it supports:

- Serve images from a local folder
- Serve images directly from S3
- Specify width and height of new image
- Specify alignment of image in new canvas (top, center, bottom, left, center, right)
- Make the server respond to hashes based on original content (checksum) instead of original filenames
- Expose known hashes as a json map at `/filenames.json`

## CLI usage

```
Usage of coresize:
  -port=8080: Port to listen on
  -aws-client-key="": Only used when pull-from=s3
  -aws-secret-key="":
  -bucket="": S3 bucket
  -v=false: Be more verbose
```

## GET `/filenames.json`

Example response:

```json
{
  "hashes": true,
  "files": {
    "agreda.png": "8e64ea6d-agreda.png",
    "aguia_branca.png": "9095612e-aguia_branca.png",
    "alsa.png": "9a8b98ad-alsa.png",
  }
}
```

## GET `/i/:filename`

Serves a file resized on-the-fly to the right format. The filename must include the hash if **coresize** was started with the `-hash` flag.

Parameters:

- `filename` (string) Filename to render

Query string parameters:

- `width` (int) Width of rendered image
- `height` (int) Height of rendered image
- `align` (enum{tl,tc,tr,cl,cc,cr,bl,bc,br}) How to align image
  - First character is `x` axis alignment {top, center, bottom}
  - Second character is `y` axis alignment {left, center, right}

```
+--------------+
|tl    tc    tr|
|              |
|cl    cc    cr|
|              |
|bl    bc    br|
+--------------+
```

Example request:

```
/i/1c2cc361-brewster.png?x=600&y=300&align=tc
```

Example response:

```
Some nice binary :D
```

## License

MIT

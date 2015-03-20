# coresize

_core, because it's simple, resize, because it serves resized versions of your images_

**coresize** aims to be a small server that can load images from local disk, S3,
or zip and serve resized and alligned versions of those images. Here's what it supports:

- Use local folder of images
- Pull down folder of images from S3
- Pull down folder of images from a `.zip` (not yet implemented)
- Serve over HTTP or HTTPS
- Specify width and height of new image
- Specify alignment of image in new canvas (t,b,l,r)
- Make the server respond to hashes based on original content instead of filenames
- Expose json hashes map at `/filenames.json`

## CLI usage

```
Usage of coresize:
  -aws-client-key="": Only used when pull-from=s3
  -aws-secret-key="":
  -folder-name="files/": Local folder where images to serve are located and will be pulled to
  -force-pull=false: Force fetching images from remote
  -hash=false: Answer to hashed filenames
  -port=8080: Port to listen on
  -pull-from="": Either 's3' or 'http'
  -pull-from-url="": S3 location or http location
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

- `x` (int) Width of rendered image
- `y` (int) Height of redered image
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

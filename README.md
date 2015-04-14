# coresize

_core, because it's simple, resize, because it serves resized versions of your images_

**coresize** aims to be a small server that can load images from local disk or S3 and serve resized and aligned versions of those images "on the fly". Here's what it supports:

- Serve images directly from S3
- Cache images locally to skip the S3 roundtrip next request
- Specify width and height of new image
- Specify alignment of image in new canvas (top, center, bottom, left, center, right)

## CLI usage

```
Usage of coresize:
  -port=8080: Port to listen on
  -aws-client-key="": Only used when pull-from=s3
  -aws-secret-key="":
  -bucket="": S3 bucket
  -v=false: Be more verbose
```

## GET `/v1/i/:filename`

Serves a file resized on-the-fly to the right format.

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
/v1/i/brewster.png?hash=1c2cc361&width=600&height=300&align=tc
```

Example response:

```
Some nice binary :D
```

## License

MIT

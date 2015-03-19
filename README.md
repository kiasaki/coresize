# coresize

_core, because it's simple, resize, because it serves resized versions of your images_

**coresize** aims to be a small server that can load images from local disk, S3,
or zip and serve resized and alligned versions of those images. Here's what it supports:

- Use local folder of images
- Pull down folder of images from S3
- Pull down folder of images from a `.zip`
- Serve over HTTP or HTTPS
- Specify width and height of new image
- Specify alignment of image in new canvas (t,b,l,r)
- Make the server respond to hashes based on original content instead of filenames
- Expose json hashes map at `/filepaths.json`

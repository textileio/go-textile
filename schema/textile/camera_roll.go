package textile

var CameraRoll = `
{
  "name": "camera_roll",
  "pin": true,
  "links": {
    "raw": {
      "use": ":file",
      "mill": "/blob"
    },
    "exif": {
      "use": "raw",
      "mill": "/image/exif"
    },
    "thumb": {
      "use": "raw",
      "pin": true,
      "mill": "/image/resize",
      "opts": {
        "width": "320",
        "quality": "80"
      }
    }
  }
}
`

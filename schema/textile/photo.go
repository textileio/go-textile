package textile

var Photos = `
{
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
    "large": {
      "use": "raw",
      "mill": "/image/resize",
      "opts": {
        "width": "1600",
        "quality": "80"
      }
    },
    "medium": {
      "use": "raw",
      "mill": "/image/resize",
      "opts": {
        "width": "800",
        "quality": "80"
      }
    },
    "small": {
      "use": "raw",
      "mill": "/image/resize",
      "opts": {
        "width": "320",
        "quality": "80"
      }
    },
    "thumb": {
      "use": "raw",
      "pin": true,
      "mill": "/image/resize",
      "opts": {
        "width": "100",
        "quality": "80"
      }
    }
  }
}
`

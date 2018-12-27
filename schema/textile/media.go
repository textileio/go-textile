package textile

var Media = `
{
  "name": "media",
  "pin": true,
  "links": {
    "large": {
      "use": ":file",
      "mill": "/image/resize",
      "opts": {
        "width": "800",
        "quality": "80"
      }
    },
    "small": {
      "use": ":file",
      "mill": "/image/resize",
      "opts": {
        "width": "320",
        "quality": "80"
      }
    },
    "thumb": {
      "use": "large",
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

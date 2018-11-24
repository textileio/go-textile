package textile

var Avatars = `
{
  "name": "avatar",
  "pin": true,
  "links": {
    "large": {
      "use": ":file",
      "pin": true,
      "plaintext": true,
      "mill": "/image/resize",
      "opts": {
        "width": "320",
        "quality": "75"
      }
    },
    "small": {
      "use": ":file",
      "pin": true,
      "plaintext": true,
      "mill": "/image/resize",
      "opts": {
        "width": "100",
        "quality": "75"
      }
    }
  }
}
`

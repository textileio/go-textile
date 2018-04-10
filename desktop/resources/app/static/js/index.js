let index = {
  init: function() {
    // Init
    asticode.loader.init();
    asticode.modaler.init();
    asticode.notifier.init();

    // Wait for astilectron to be ready
    document.addEventListener('astilectron-ready', function() {
      // Listen
      index.listen();

      let pairNewDevice = function () {
        // Get a thumbnail (i know it exists)
        astilectron.sendMessage({name: "peer.qr", payload: ""}, function (message) {
          // Check error
          console.log(message)
          if (message.name === "error") {
            asticode.notifier.error("Error");
            return
          }
          let qrCode = document.querySelector('.peerQRCode');
          qrCode.setAttribute('src', "data:image/png;base64," + message.payload.png + "")
          let pairCode = document.querySelector('.pairCode');
          pairCode.innerText = message.payload.code
          let modal = document.querySelector('.modal');
          modal.classList.toggle('modal-open');
        })
      }
      setupQRModal(pairNewDevice)
    })
  },
  listen: function() {
    astilectron.onMessage(function(message) {
      switch (message.name) {
        case "new.log":
          let lines = document.getElementById("log").innerText;
          document.getElementById("log").innerText = message.payload + "<br/>" + lines
          break;
        case "new.image":
          addNewGalleryImage("https://gateway.ipfs.io/ipfs/QmR8mGCutBWDPBc9zdpPZPoRYRAJS7BMZhJtqHeFtJp2ma/thumb.jpg")
          break;
      }
    });
  },
};

function setupQRModal(pairMethod) {
  let modal = document.querySelector('.modal');
  let closeButtons = document.querySelectorAll('.close-modal');
  // set open modal behaviour
  document.getElementById('pairNewDevice').addEventListener('click', pairMethod);
  // set close modal behaviour
  for (i = 0; i < closeButtons.length; ++i) {
    closeButtons[i].addEventListener('click', function() {
      modal.classList.toggle('modal-open');
    });
  }
  // close modal if clicked outside content area
  document.querySelector('.modal-inner').addEventListener('click', function() {
    modal.classList.toggle('modal-open');
  });
  // prevent modal inner from closing parent when clicked
  document.querySelector('.modal-content').addEventListener('click', function(e) {
    e.stopPropagation();
  });
}

function addNewGalleryImage(src) {
  let entry = document.createElement('div');
  entry.className = 'two columns entry';
  let gallery = document.getElementById('gallery');
  gallery.insertBefore(entry, gallery.firstChild);

  let thumb = document.createElement('div');
  thumb.className = 'row thumb';
  entry.appendChild(thumb);

  console.log(src)
  let img = document.createElement("img");
  img.className = 'u-max-full-width';
  thumb.appendChild(img);
  img.setAttribute('src', src);

  let label = document.createElement('div');
  label.className = 'row label u-max-full-width';
  label.innerText = "Yesterday"
  entry.appendChild(label);
}
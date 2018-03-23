let index = {
    // about: function(html) {
    //     let c = document.createElement("div");
    //     c.innerHTML = html;
    //     asticode.modaler.setContent(c);
    //     asticode.modaler.show();
    // },
    // addFolder(name, path) {
    //     let div = document.createElement("div");
    //     div.className = "dir";
    //     div.onclick = function() { index.explore(path) };
    //     div.innerHTML = `<i class="fa fa-folder"></i><span>` + name + `</span>`;
    //     document.getElementById("dirs").appendChild(div)
    // },
    init: function() {
        // Init
        asticode.loader.init();
        asticode.modaler.init();
        asticode.notifier.init();

        // Wait for astilectron to be ready
        document.addEventListener('astilectron-ready', function() {
            // Listen
            index.listen();

            document.getElementById("path").innerHTML = "derp";

            // just to see the peerid in the interface
            asticode.notifier.info("derp")
            astilectron.sendMessage({name: "ipfs.peerId", payload: "get"}, function(message) {
                asticode.notifier.info(message)

                // Get a thumbnail (i know it exists)
                astilectron.sendMessage({name: "ipfs.getPath", payload: "QmSuyWQoXNyXSmkyFXVb1xsPHuXqVZgHvjpjB3R3Y9YZvj/thumb.JPG"}, function(message) {
                    asticode.notifier.info(message)
                    // Check error
                    if (message.name === "error") {
                        asticode.notifier.error(message.payload);
                        return
                    }
                    document.getElementById("test").src = "data:image/png;base64, " + message.payload;
                })


                // Check error
                if (message.name === "error") {
                    asticode.notifier.error(message.payload);
                    return
                }
                document.getElementById("message").innerHTML = message.payload;
            })

        })
    },
    listen: function() {
        astilectron.onMessage(function(message) {
            console.log(message.name)
            switch (message.name) {
                case "about":
                    // index.about(message.payload);
                    return {payload: "payload"};
                    break;
                case "check.out.menu":
                    asticode.notifier.info(message.payload);
                    break;
            }
        });
    }
};
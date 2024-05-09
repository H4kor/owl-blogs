/**
 * Add "drop to upload" to a text area
 * @param {string} id id of a textarea
 */
function addFileDrop(id) {
    // deactivate file drop on body 
    // this prevents accidentally opening the file instead uploading
    document.body.ondrop = (ev) => {
        ev.preventDefault();
        return false
    }
    document.body.ondragover = (ev) => {
        ev.preventDefault();
        return false
    }

    // get field
    const textArea = document.getElementById(id)

    /**
     * Uploads a single file and add markdown to textarea
     * @param {File} file file object to upload
     */
    function processFile(file) {
        console.log(`name = ${file.name}`);
        console.log(`size = ${file.size}`);
        console.log(`type = ${file.type}`);
        console.log(file)

        const formData = new FormData()
        formData.append("file", file)

        textArea.classList.add("drop-file-process")
        fetch(
            "/admin/api/binaries",
            {
                method: "POST",
                body: formData
            }
        ).then((resp) => {
            return resp.json()
        }).then(data => {
            if (file.type.split("/")[0] == "image") {
                textArea.value += `\n![](${data.location})`
            } else {
                textArea.value += `\n[${file.name}](${data.location})`
            }
            textArea.classList.remove("drop-file-process")
        }).catch(err => {
            console.error(err);
            textArea.classList.add("drop-file-error")
            setTimeout(() => {
                textArea.classList.remove("drop-file-error")
            }, 2000)
        }).finally(() => {
            textArea.classList.remove("drop-file-process")
        })

    }

    function dropHandler(ev) {
        textArea.classList.remove("drop-file")
        ev.preventDefault();

        if (ev.dataTransfer.items) {
            // Use DataTransferItemList interface to access the file(s)
            [...ev.dataTransfer.items].forEach((item, i) => {
                // If dropped items aren't files, reject them
                if (item.kind === "file") {
                    const file = item.getAsFile();
                    processFile(file)
                }
            });
        } else {
            // Use DataTransfer interface to access the file(s)
            [...ev.dataTransfer.files].forEach((file, i) => {
                processFile(file)
            });
        }
    }

    function dragOverHandler(ev) {
        textArea.classList.add("drop-file")
        ev.preventDefault();
    }

    function dragLeaveHandler(ev) {
        textArea.classList.remove("drop-file")
        ev.preventDefault();
    }    

    textArea.ondrop = dropHandler;
    textArea.ondragover = dragOverHandler;
    textArea.ondragleave = dragLeaveHandler;

}
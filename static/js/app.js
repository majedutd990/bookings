// Prompt is our java script module for all alerts notification and custom pop ups

function Prompt() {
    let toast = function (c) {
        const {
            msg = "",
            icon = 'success',
            position = 'top-end'
        } = c;
        const Toast = Swal.mixin({
            toast: true,
            title: msg,
            position: position,
            icon: icon,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.addEventListener('mouseenter', Swal.stopTimer)
                toast.addEventListener('mouseleave', Swal.resumeTimer)
            }
        })

        Toast.fire({})
    }
    let success = function (c) {
        const {
            msg = '',
            title = '',
            footer = '',
        } = c;
        Swal.fire({
            icon: 'success',
            title: title,
            text: msg,
            footer: footer
        })
    }
    let error = function (c) {
        const {
            msg = '',
            title = '',
            footer = '',
        } = c;
        Swal.fire({
            icon: 'error',
            title: title,
            text: msg,
            footer: footer
        })
    }
    let custom = async function (c) {
        const {
            icon = "",
            msg = "",
            title = "",
            showConfirmButton = true,
        } = c;
        const {value: result} = await Swal.fire({
            icon: icon,
            title: title,
            html: msg,
            width: 600,
            showConfirmButton: showConfirmButton,
            willOpen() {
                if (c.willOpen !== undefined) {
                    c.willOpen();
                }
            },
            backdrop: false,
            showCancelButton: true,
            focusConfirm: false,
            preConfirm: () => {
                if (c.preConfirm() !== undefined) {
                    c.preConfirm()
                }
            },
            didOpen: () => {
                if (c.didOpen() !== undefined) {
                    c.didOpen();
                }
            }
        })
        if (result) {
            if (result.dismiss !== Swal.DismissReason.cancel) {
                if (result.value !== "") {
                    if (c.callback !== undefined) {
                        c.callback(result);
                    } else {
                        c.callback(false);
                    }
                } else {
                    c.callback(false);
                }
            }
        }
    }

    return {
        toast: toast,
        success: success,
        error: error,
        custom: custom,
    }
}


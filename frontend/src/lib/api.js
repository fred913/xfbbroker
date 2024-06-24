async function get(url, params) {
    let sessionId = localStorage.getItem("sessionId")
    let p = new URLSearchParams(params)
    if (sessionId.length > 0) {
        p.set("sessionId", sessionId)
    }
    return await fetch(url + '?' + p.toString())
}

async function getConfig() {
    return await get("http://localhost:8000/_/config").then(r => r.json())
}

async function checkSignpay() {
    return await get("http://localhost:8000/_/xfb/signpay").then(r => {
        if (r.status == 201) {
            return r.headers.get('Location')
        } else if (r.status == 200) {
            return ""
        } else {
            throw Error(`unknown statusCode ${r.status}: ${r.text()}`)
        }
    })
}

export {
    getConfig, checkSignpay
}
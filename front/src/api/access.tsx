export function sendAccess(encodedFilePath: string) {
    const url = new URL('/api/access', window.location.origin);
    url.searchParams.set('path', encodedFilePath);

    if (encodedFilePath === "") {
        return
    }

    fetch(url.toString(), {
        method: 'GET',
    }).then(res => {
        if (!res.ok) {
            console.error("Failed to send progress");
        }
    }).catch(err => {
        console.error("Network error while sending progress", err);
    });
}

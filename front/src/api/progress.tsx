export function sendProgress(encodedFilePath: string, currentPosition: string, progress: number) {
    const url = new URL('/api/progress', window.location.origin);
    url.searchParams.set('path', encodedFilePath);
    url.searchParams.set('position', currentPosition);
    url.searchParams.set('progress', progress.toString());

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

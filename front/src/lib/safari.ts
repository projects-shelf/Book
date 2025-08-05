export function isSafari(): boolean {
    const ua = navigator.userAgent;
    return (
        /Safari/.test(ua) &&
        !/Chrome/.test(ua) &&
        !/Chromium/.test(ua) &&
        !/Android/.test(ua)
    );
}
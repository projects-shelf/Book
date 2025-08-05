import { useEffect, useMemo, useRef, useState } from 'react';
import { ViewerLayout } from './layout';
import { ViewerOptions } from './sheet';
import { Slider } from "@/components/ui/slider"
import { sendProgress } from '@/api/progress';
import { sendAccess } from '@/api/access';
import { useWindowSize } from '@/hooks/windowSize';
import { isSafari } from '@/lib/safari';

interface CBRViewerProps {
    fileUrl: string;
    initialPage?: number;
}

const big_number = 99999

export function CBRViewer({ fileUrl, initialPage = 1 }: CBRViewerProps) {
    const [numPages, setNumPages] = useState<number | null>(null);
    const [pageNumber, setPageNumber] = useState(initialPage);

    const [direction, setDirection] = useState<"ltr" | "rtl">("ltr")
    const [spread, setSpread] = useState<"none" | "odd" | "even">("none")
    // const [fontSize, setFontSize] = useState<number | "">(16)

    const [scale, setScale] = useState(0.5);
    const [sliderValue, setSliderValue] = useState([initialPage]);

    const imgRef = useRef<HTMLImageElement>(null);

    useEffect(() => {
        const saved = localStorage.getItem("viewerOptions")
        if (saved) {
            try {
                const parsed = JSON.parse(saved) as ViewerOptions
                setDirection(parsed.direction)
                setSpread(parsed.spread)
                // setFontSize(parsed.fontSize)
            } catch (e) {
                console.warn("Failed to parse viewerOptions from localStorage")
            }
        }


    }, [])

    const toPrev = () => {
        var delta = 1
        switch (spread) {
            case "odd":
                if (pageNumber % 2 === 1) {
                    delta = 2
                } else {
                    delta = 3
                }
                break;
            case "even":
                if (pageNumber % 2 === 1) {
                    delta = 3
                } else {
                    delta = 2
                }
                break;
        }

        const newPage = Math.max(pageNumber - delta, 1)
        setPageNumber(newPage)
        setSliderValue([newPage])
        sendProgress(encodedFilePath ?? "", newPage.toString(), newPage / (numPages ?? big_number))
    }

    const toNext = () => {
        if (numPages !== null) {
            var delta = 1
            switch (spread) {
                case "odd":
                    if (pageNumber % 2 === 1) {
                        delta = 2
                    }
                    break;
                case "even":
                    if (pageNumber % 2 === 0) {
                        delta = 2
                    }
                    break;
            }

            const newPage = Math.min(pageNumber + delta, numPages)
            setPageNumber(newPage)
            setSliderValue([newPage])
            sendProgress(encodedFilePath ?? "", newPage.toString(), newPage / (numPages ?? big_number))
        }
    }

    const standerisedPageNumber = useMemo(() => {
        switch (spread) {
            case "odd":
                return pageNumber % 2 === 0 ? pageNumber - 1 : pageNumber;
            case "even":
                return (pageNumber !== 1 && pageNumber % 2 === 1) ? pageNumber - 1 : pageNumber;
            default:
                return pageNumber;
        }
    }, [pageNumber, spread]);

    const isSpreads = useMemo(() => {
        if (!numPages) return false;
        switch (spread) {
            case "odd":
                return standerisedPageNumber !== numPages;
            case "even":
                return standerisedPageNumber !== 1 && standerisedPageNumber !== numPages;
            default:
                return false;
        }
    }, [standerisedPageNumber, spread, numPages]);

    const { width: windowWidth, height: windowHeight } = useWindowSize()
    const windowAspect = windowWidth / windowHeight

    const encodedFilePath = useMemo(() => {
        try {
            const url = new URL(fileUrl, window.location.origin);
            return url.searchParams.get("path");
        } catch (e) {
            console.warn("Invalid URL:", fileUrl);
            return "";
        }
    }, [fileUrl]);

    useMemo(() => {
        sendAccess(encodedFilePath ?? "");
    }, [fileUrl]);

    useEffect(() => {
        const fetchPageCount = async () => {
            try {
                const res = await fetch(`/book/cbr/pages?path=${encodedFilePath}`);
                if (!res.ok) throw new Error("Failed to fetch page count");
                const data = await res.json();
                const total = data.pages ?? 1;

                setNumPages(total);

                if (initialPage < 1 || initialPage > total) {
                    setPageNumber(1);
                    setSliderValue([1]);
                }
            } catch (err) {
                console.error("Error fetching CBR page count:", err);
            }
        };

        fetchPageCount();
    }, []);

    const [mLoaded, setMLoaded] = useState(false);
    const [nLoaded, setNLoaded] = useState(false);

    if (isSafari() && !(mLoaded && nLoaded)) {
        return (
            <div>
                <p>Loading...</p>
                <div style={{ position: 'absolute', width: 0, height: 0, overflow: 'hidden' }}>
                    <img style={{ width: "auto", height: "auto", maxWidth: "none", maxHeight: "none" }}
                        src={`/book/pdf?path=${encodedFilePath}&page=${standerisedPageNumber}`}
                        onLoad={() => setMLoaded(true)}
                        onError={() => setMLoaded(true)} />
                    <img style={{ width: "auto", height: "auto", maxWidth: "none", maxHeight: "none" }}
                        src={`/book/pdf?path=${encodedFilePath}&page=${standerisedPageNumber + 1}`}
                        onLoad={() => setNLoaded(true)}
                        onError={() => setNLoaded(true)} />
                </div>
            </div>
        );
    }

    return (
        <ViewerLayout
            onOptionChanged={(opt) => {
                setDirection(opt.direction)
                setSpread(opt.spread)
                // setFontSize(opt.fontSize)
            }}
            onLeft={() => {
                if (direction === "ltr") {
                    toPrev()
                } else {
                    toNext()
                }
            }}
            onRight={() => {
                if (direction === "ltr") {
                    toNext()
                } else {
                    toPrev()
                }
            }}>
            <div className="relative w-full max-w-full h-screen overflow-hidden">
                <div style={{ display: "flex", justifyContent: 'center', transform: `scale(${scale})`, transformOrigin: 'top center' }}>
                    {direction === 'rtl' && isSpreads && (
                        <img style={{ width: "auto", height: "auto", maxWidth: "none", maxHeight: "none" }}
                            key={`l-/book/cbr?path=${encodedFilePath}&page=${Math.min(standerisedPageNumber + 1, numPages ?? 1)}`}
                            src={`/book/cbr?path=${encodedFilePath}&page=${Math.min(standerisedPageNumber + 1, numPages ?? 1)}`} />
                    )}
                    <img style={{ width: "auto", height: "auto", maxWidth: "none", maxHeight: "none" }}
                        key={`m-/book/cbr?path=${encodedFilePath}&page=${standerisedPageNumber}`}
                        src={`/book/cbr?path=${encodedFilePath}&page=${standerisedPageNumber}`}
                        ref={imgRef} onLoad={() => {
                            const pageWidth = imgRef.current!.naturalWidth;
                            const pageHeight = imgRef.current!.naturalHeight;
                            const pagesWidth = isSpreads ? pageWidth * 2 : pageWidth;
                            const pageAspect = pagesWidth / pageHeight;

                            if (pageAspect > windowAspect) {
                                setScale(windowWidth / pagesWidth);
                            } else {
                                setScale(windowHeight / pageHeight);
                            }
                        }} />
                    {direction === 'ltr' && isSpreads && (
                        <img style={{ width: "auto", height: "auto", maxWidth: "none", maxHeight: "none" }}
                            key={`r-/book/cbr?path=${encodedFilePath}&page=${Math.min(standerisedPageNumber + 1, numPages ?? 1)}`}
                            src={`/book/cbr?path=${encodedFilePath}&page=${Math.min(standerisedPageNumber + 1, numPages ?? 1)}`} />
                    )}

                    {/* Shadow Page */}
                    <div style={{ position: 'absolute', width: 0, height: 0, overflow: 'hidden' }}>
                        <img style={{ width: "auto", height: "auto", maxWidth: "none", maxHeight: "none" }}
                            src={`/book/cbr?path=${encodedFilePath}&page=${Math.max(standerisedPageNumber - 1, 1)}`} />
                        <img style={{ width: "auto", height: "auto", maxWidth: "none", maxHeight: "none" }}
                            src={`/book/cbr?path=${encodedFilePath}&page=${Math.max(standerisedPageNumber - 2, 1)}`} />
                        <img style={{ width: "auto", height: "auto", maxWidth: "none", maxHeight: "none" }}
                            src={`/book/cbr?path=${encodedFilePath}&page=${Math.min(standerisedPageNumber + 1, numPages ?? 1)}`} />
                        <img style={{ width: "auto", height: "auto", maxWidth: "none", maxHeight: "none" }}
                            src={`/book/cbr?path=${encodedFilePath}&page=${Math.min(standerisedPageNumber + 2, numPages ?? 1)}`} />
                        <img style={{ width: "auto", height: "auto", maxWidth: "none", maxHeight: "none" }}
                            src={`/book/cbr?path=${encodedFilePath}&page=${Math.min(standerisedPageNumber + 3, numPages ?? 1)}`} />
                    </div>
                    {/* Shadow Page */}

                </div>
            </div>

            <div className="absolute left-0 bottom-3 w-full" style={{ opacity: 0.25 }}>
                <Slider
                    value={sliderValue}
                    min={1}
                    max={numPages ?? big_number}
                    step={1}
                    onValueChange={(values) => {
                        setSliderValue(values)
                    }}
                    onValueCommit={(values) => {
                        if (values.length > 0) {
                            setPageNumber(values[0])
                            sendProgress(encodedFilePath ?? "", values[0].toString(), values[0] / (numPages ?? big_number))
                        }
                    }}
                    dir={direction === 'rtl' ? "rtl" : "ltr"}
                    tabIndex={-1}
                    onKeyDown={(e) => {
                        if (["ArrowLeft", "ArrowRight", "ArrowUp", "ArrowDown"].includes(e.key)) {
                            e.preventDefault();
                        }
                    }}
                />
            </div>
        </ViewerLayout >
    );
}
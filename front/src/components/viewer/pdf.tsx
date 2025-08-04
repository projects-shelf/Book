import { useEffect, useMemo, useState } from 'react';
import { Document, Page, pdfjs } from 'react-pdf';
import 'react-pdf/dist/Page/AnnotationLayer.css';
import 'react-pdf/dist/Page/TextLayer.css';
import { ViewerLayout } from './layout';
import { ViewerOptions } from './sheet';
import { useWindowSize } from '@/hooks/windowSize';
import { Slider } from "@/components/ui/slider"
import { sendProgress } from '@/api/progress';
import { sendAccess } from '@/api/access';

pdfjs.GlobalWorkerOptions.workerPort =
    new Worker("/pdf.worker.min.js");

interface PDFViewerProps {
    fileUrl: string;
    initialPage?: number;
}

const scaling_factor = 3
const big_number = 99999

export function PDFViewer({ fileUrl, initialPage = 1 }: PDFViewerProps) {
    const [numPages, setNumPages] = useState<number | null>(null);
    const [pageNumber, setPageNumber] = useState(initialPage);

    const [direction, setDirection] = useState<"ltr" | "rtl">("ltr")
    const [spread, setSpread] = useState<"none" | "odd" | "even">("none")
    // const [fontSize, setFontSize] = useState<number | "">(16)

    const [scale, setScale] = useState(0.5);
    const [sliderValue, setSliderValue] = useState([initialPage]);

    const onDocumentLoadSuccess = ({ numPages }: { numPages: number }) => {
        setNumPages(numPages);
        if (initialPage < 1 || initialPage > numPages) {
            setPageNumber(1);
            setSliderValue([1]);
        }
    };

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

    const { width: windowWidth, height: windowHeight } = useWindowSize()
    const windowAspect = windowWidth / windowHeight

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

    const encodedFilePath = useMemo(() => {
        try {
            const url = new URL(fileUrl, window.location.origin);
            return url.searchParams.get("path");
        } catch (e) {
            console.warn("Invalid URL:", fileUrl);
            return "";
        }
    }, [fileUrl]);

    sendAccess(encodedFilePath ?? "")

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
                <Document file={fileUrl} onLoadSuccess={onDocumentLoadSuccess}>
                    <div style={{
                        display: "flex", justifyContent: 'center', transform: `scale(${scale})`, transformOrigin: 'top center'
                    }}>
                        {direction === 'rtl' && isSpreads && (
                            <Page pageNumber={Math.min(standerisedPageNumber + 1, numPages ?? 1)} scale={scaling_factor} />
                        )}
                        <Page
                            pageNumber={standerisedPageNumber}
                            scale={scaling_factor}
                            onLoadSuccess={(page) => {
                                const { width: pageWidth, height: pageHeight } = page.getViewport({ scale: scaling_factor });
                                const pagesWidth = (isSpreads ? pageWidth * 2 : pageWidth)
                                const pageAspect = pagesWidth / pageHeight
                                if (pageAspect > windowAspect) {
                                    setScale(windowWidth / pagesWidth)
                                } else {
                                    setScale(windowHeight / pageHeight)
                                }
                            }}
                        />
                        {direction === 'ltr' && isSpreads && (
                            <Page pageNumber={Math.min(standerisedPageNumber + 1, numPages ?? 1)} scale={scaling_factor} />
                        )}

                        {/* Shadow Page */}
                        <div style={{ position: 'absolute', width: 0, height: 0, overflow: 'hidden' }}>
                            <Page pageNumber={Math.max(standerisedPageNumber - 1, 1)} scale={scaling_factor} />
                            <Page pageNumber={Math.max(standerisedPageNumber - 2, 1)} scale={scaling_factor} />
                            <Page pageNumber={Math.min(standerisedPageNumber + 1, numPages ?? 1)} scale={scaling_factor} />
                            <Page pageNumber={Math.min(standerisedPageNumber + 2, numPages ?? 1)} scale={scaling_factor} />
                            <Page pageNumber={Math.min(standerisedPageNumber + 3, numPages ?? 1)} scale={scaling_factor} />
                        </div>
                        {/* Shadow Page */}

                    </div>
                </Document>
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
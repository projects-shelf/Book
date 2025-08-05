import { useEffect, useMemo, useRef, useState } from 'react';
import { ReactReader, ReactReaderStyle, type IReactReaderStyle } from 'react-reader';
import type { Rendition } from 'epubjs'
import { ViewerLayout } from './layout';
import { Slider } from "@/components/ui/slider"
import { ViewerOptions } from './sheet';
import { sendAccess } from '@/api/access';
import { sendProgress } from '@/api/progress';

type EPUBViewerProps = {
    fileUrl: string;
    initialPage?: string;
};

const customReaderStyle: IReactReaderStyle = {
    ...ReactReaderStyle,
    arrow: {
        ...ReactReaderStyle.arrow,
        display: 'none',
    },
    arrowHover: {
        ...ReactReaderStyle.arrowHover,
        display: 'none',
    },
};

export function EPUBViewer({ fileUrl }: EPUBViewerProps) {
    const [direction, setDirection] = useState<"ltr" | "rtl">("ltr")
    const [fontSize, setFontSize] = useState<number | "">(16)

    const [sliderValue, setSliderValue] = useState([0.0]);

    const [location, setLocation] = useState<string | number>(0)
    const [locationsGenerated, setLocationsGenerated] = useState(false)
    const rendition = useRef<Rendition | undefined>(undefined)

    useEffect(() => {
        const saved = localStorage.getItem("viewerOptions")
        if (saved) {
            try {
                const parsed = JSON.parse(saved) as ViewerOptions
                setDirection(parsed.direction)
                // setSpread(parsed.spread)
                setFontSize(parsed.fontSize)
            } catch (e) {
                console.warn("Failed to parse viewerOptions from localStorage")
            }
        }
    }, [])

    useEffect(() => {
        if (rendition.current) {
            const size = fontSize === "" ? "16px" : `${fontSize}px`;
            rendition.current.themes.fontSize(size);
        }
    }, [fontSize]);

    const toPrev = () => {
        rendition.current?.prev()
    }

    const toNext = () => {
        rendition.current?.next()
    }

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
        const stopArrowKeys = (e: KeyboardEvent) => {
            if (e.key === "ArrowLeft" || e.key === "ArrowRight") {
                e.stopPropagation();
                e.preventDefault();
            }
        };

        window.addEventListener("keydown", stopArrowKeys, true);

        return () => window.removeEventListener("keydown", stopArrowKeys, true);
    }, []);

    return (
        <ViewerLayout
            onOptionChanged={(opt) => {
                setDirection(opt.direction)
                // setSpread(opt.spread)
                setFontSize(opt.fontSize)
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
            <div tabIndex={-1} style={{ width: '100vw', height: '100vh' }}>
                <ReactReader
                    url={fileUrl}
                    location={location}
                    locationChanged={(loc) => {
                        setLocation(loc)
                        if (locationsGenerated && rendition.current?.book.locations) {
                            const percentage = rendition.current?.book.locations.percentageFromCfi(loc)
                            setSliderValue([percentage])
                            sendProgress(encodedFilePath ?? "", loc, percentage)
                        }
                    }}
                    getRendition={(_rendition: Rendition) => {
                        rendition.current = _rendition
                        const size = fontSize === "" ? "16px" : `${fontSize}px`;
                        _rendition.themes.fontSize(size);

                        _rendition.book.ready.then(() => {
                            console.log("Book ready, generating locations...");
                            return _rendition.book.locations.generate(1024);
                        }).then(() => {
                            setLocationsGenerated(true);
                        }).catch((error) => {
                            console.error("Error generating locations:", error);
                        });
                    }}
                    showToc={false}
                    readerStyles={customReaderStyle}
                    isRTL={direction === 'rtl' ? true : false}
                />
            </div>
            <div className="absolute left-0 bottom-3 w-full z-50" style={{ opacity: 0.25 }}>
                <Slider
                    value={sliderValue}
                    min={0}
                    max={1.0}
                    step={0.001}
                    onValueChange={(values) => {
                        setSliderValue(values)
                    }}
                    onValueCommit={(values) => {
                        if (values.length > 0) {
                            if (locationsGenerated && rendition.current?.book.locations) {
                                rendition.current.display(rendition.current?.book.locations.cfiFromPercentage(values[0]))
                            }
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

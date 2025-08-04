import { useState, useEffect } from "react";
import { PDFViewer } from "@/components/viewer/pdf";

function useQuery() {
    const [query, setQuery] = useState<URLSearchParams | null>(null);

    useEffect(() => {
        setQuery(new URLSearchParams(window.location.search));
    }, []);

    return query;
}

export function PDFViewerPage() {
    const query = useQuery();

    if (!query) return <p>Loading...</p>;

    const path = query.get("path") ?? "";
    const parsed = parseInt(query.get("position") ?? "", 10);
    const position = isNaN(parsed) ? 1 : parsed;

    return <PDFViewer fileUrl={`/book/pdf?path=${encodeURIComponent(path)}`} initialPage={position} />;
}

import { useState, useEffect } from "react";
import { CBZViewer } from "@/components/viewer/cbz";

function useQuery() {
    const [query, setQuery] = useState<URLSearchParams | null>(null);

    useEffect(() => {
        setQuery(new URLSearchParams(window.location.search));
    }, []);

    return query;
}

export function CBZViewerPage() {
    const query = useQuery();

    if (!query) return <p>Loading...</p>;

    const path = query.get("path") ?? "";
    const parsed = parseInt(query.get("position") ?? "", 10);
    const position = isNaN(parsed) ? 1 : parsed;

    return <CBZViewer fileUrl={`/book/cbz?path=${encodeURIComponent(path)}`} initialPage={position} />;
}
